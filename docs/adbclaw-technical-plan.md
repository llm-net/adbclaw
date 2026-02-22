# adbclaw 技术方案

> 基于 DroidRun / mobile-use 深度分析后形成的技术方案。
> adbclaw 是一个**纯 CLI 工具**，不包含 LLM/Agent 逻辑，供 OpenClaw 或其他 agent/bot 通过 skill 调用。

---

## 1. 定位与边界

### adbclaw 是什么

一个**设备控制工具层**，类似 `kubectl` 之于 Kubernetes、`gh` 之于 GitHub：
- 提供完整的 Android 设备操作 CLI
- 所有命令输出结构化 JSON，供 agent 解析
- 附带 skill 描述文件，让 agent 知道自己能调什么
- 内置 sendevent 反侦测输入引擎
- 单 binary 分发，零依赖运行

### adbclaw 不是什么

- 不包含 LLM 调用、prompt、agent loop
- 不做任务规划、屏幕理解、决策推理
- 不依赖 Python 运行时
- 不需要在手机上安装 APK（Portal App 类的方案排除）

### 与竞品的关键差异

| | DroidRun | mobile-use | adbclaw |
|---|---|---|---|
| 定位 | LLM Agent + 设备控制一体 | LLM Agent + 设备控制一体 | **纯设备控制工具** |
| 语言 | Python | Python | Go + C（设备端） |
| LLM 依赖 | llama-index (必须) | LangGraph + LangChain (必须) | **无** |
| 输入注入 | `adb shell input`（可检测） | `adb shell input`（可检测） | **sendevent /dev/input**（高隐蔽） |
| 反侦测 | 无 | 无 | **stealth / deep-stealth** |
| 手机端依赖 | Portal APK + AccessibilityService | UIAutomator2 server | **无 APK**，仅推送轻量 binary |
| UI 树获取 | Portal Accessibility Service | UIAutomator2 dump_hierarchy | `uiautomator dump`（原生命令） |
| Agent 集成 | 内嵌 agent | 内嵌 agent | **skill 文件 + JSON stdout** |
| 分发 | pip install | pip install | **单 binary** |

---

## 2. 从竞品分析中提取的设计决策

### 2.1 从 DroidRun 学到的

**采纳：**
- `DeviceDriver` 抽象接口设计思路好，但 adbclaw 用 CLI 子命令代替 — 每个 driver 方法 = 一个子命令
- `ToolRegistry` 的 capability gating 概念 — adbclaw 在 `device probe` 时探测设备能力，skill 文件中标注哪些命令可用
- UI 树的 `IndexedFormatter`（给元素编号）— adbclaw 的 `ui tree` 输出应包含编号，方便 agent 用 index 引用元素
- 元素坐标解析 `get_element_coords(index)` — adbclaw 提供 `ui tap --index 3` 直接按编号点击

**不采纳：**
- Portal APK 依赖 — 要求用户安装额外 App + 开启无障碍服务，增加侦测面（AccessibilityService 可被枚举）
- llama-index Workflow 编排 — adbclaw 不做 agent
- 文本输入走 Portal IME — 改用 ADB `input text` + 对特殊字符的转义处理，或用 ADB broadcast 方式

### 2.2 从 mobile-use 学到的

**采纳：**
- `Target` 三级回退（bounds → resource_id → text）— adbclaw 的 `tap` 命令支持 `--index`、`--id`、`--text`、`--bounds` 四种定位方式
- 截图 + UI 层级并行获取 — `adbclaw observe` 命令同时返回截图和 UI 树
- `scratchpad` 持久化 kv 存储 — adbclaw 提供 `note` 子命令，agent 可跨步骤存取数据
- 截图压缩（JPEG quality 控制）— agent 不需要全分辨率 PNG

**不采纳：**
- UIAutomator2 Python client — 需要在手机端运行 UIAutomator2 server，是一个可被检测的自动化框架
- LangGraph 状态图 — adbclaw 不做 agent
- 8 个 LLM agent 分工 — 那是 agent 层的事

### 2.3 两个项目共同的缺陷 → adbclaw 的机会

| 缺陷 | adbclaw 的解决方案 |
|------|-------------------|
| 输入用 `adb shell input`，deviceId=-1 可检测 | sendevent 引擎，deviceId 真实 |
| 无时间抖动/轨迹模拟 | 人类化引擎（高斯抖动 + 贝塞尔曲线） |
| 依赖 Python 运行时 | Go 单 binary |
| 需要安装 APK / UIAutomator2 | 仅推送 `/data/local/tmp/adbclawd`，不装 APK |
| UI 树获取依赖 Portal/UIAutomator2 server | 用 `uiautomator dump` 原生命令 + 自行解析 XML |
| 输出格式面向 LLM prompt | 输出标准 JSON，agent 自行格式化 |

---

## 3. 整体架构

```
┌─────────────────────────────────────────────────────────┐
│  OpenClaw / 其他 Agent / Bot                             │
│  (读 skill.json 知道 adbclaw 能做什么，                    │
│   调用 CLI 命令，解析 JSON stdout)                         │
└────────────────────┬────────────────────────────────────┘
                     │ 子进程调用 / stdin-stdout pipe
┌────────────────────▼────────────────────────────────────┐
│  adbclaw CLI (Go binary, Mac/Linux)                      │
│                                                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐ │
│  │ device   │ │ input    │ │ observe  │ │ app        │ │
│  │ 连接管理  │ │ 输入注入  │ │ 屏幕/UI  │ │ App 管理   │ │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └─────┬──────┘ │
│       │            │            │              │        │
│  ┌────▼────────────▼────────────▼──────────────▼──────┐ │
│  │              ADB Transport Layer                    │ │
│  │  (自实现 ADB 协议，不依赖 adb binary)                │ │
│  └────────────────────┬───────────────────────────────┘ │
└───────────────────────┼─────────────────────────────────┘
                        │ USB / TCP
┌───────────────────────▼─────────────────────────────────┐
│  Android 设备                                            │
│                                                          │
│  /data/local/tmp/adbclawd (设备端辅助 binary)             │
│  ├── sendevent 批量写入引擎                                │
│  ├── 快速截屏 (framebuffer / screencap)                   │
│  ├── UI 树 dump + 解析                                    │
│  └── 设备能力探测                                          │
└─────────────────────────────────────────────────────────┘
```

---

## 4. 核心模块设计

### 4.1 设备端辅助程序 `adbclawd`

**为什么需要设备端程序：**

DroidRun 每次 tap 都要启动一个 `adb shell input tap` 进程（~50ms 开销）。mobile-use 同理。sendevent 更糟——一次 tap 需要 6-8 条 sendevent 命令，每条都是一个进程。

adbclawd 是一个**常驻在手机端的轻量 C 程序**（类似 scrcpy-server 或 minitouch 的设计），通过 ADB forward 的 socket 与 Mac 端通信：

```
Mac (adbclaw) ──ADB forward──> Android (adbclawd:unix socket)
```

**adbclawd 职责：**

```
adbclawd
├── /input     — 批量 sendevent 写入（一次 socket message = 一个完整手势）
├── /probe     — 探测 /dev/input 设备节点、触屏参数、协议类型
├── /screen    — screencap 并通过 socket 传回（避免 adb pull 开销）
├── /uitree    — 执行 uiautomator dump + 读取 XML 并返回
└── /info      — 设备信息（分辨率、DPI、Android 版本等）
```

**编译：**
- 用 Android NDK 交叉编译为 ARM64/ARM32 static binary
- 体积控制在 < 500KB
- 首次连接时自动 `adb push` 到 `/data/local/tmp/adbclawd`
- 运行方式: `adb shell /data/local/tmp/adbclawd --socket adbclawd.sock`

**为什么不装 APK：**
- APK 安装需要用户确认，自动化场景不友好
- 已安装 App 列表可被其他 App 枚举（`PackageManager.getInstalledPackages`）
- AccessibilityService 可被检测
- `/data/local/tmp` 下的 binary 对普通 App 不可见

### 4.2 输入注入引擎

```
┌─────────────────────────────────────────────┐
│            Input Engine (Mac 端)              │
│                                              │
│  tap(x, y) / swipe(x1,y1,x2,y2) / type()   │
│              │                               │
│              ▼                               │
│  ┌─────────────────────┐                    │
│  │   Humanizer 层       │ (可选, --humanize) │
│  │   ├─ 坐标高斯偏移    │                    │
│  │   ├─ 压力曲线生成    │                    │
│  │   ├─ 贝塞尔轨迹插值  │                    │
│  │   └─ 时间高斯抖动    │                    │
│  └──────────┬──────────┘                    │
│             ▼                                │
│  ┌─────────────────────┐                    │
│  │ Event Serializer     │                    │
│  │ 生成 sendevent 序列  │                    │
│  │ (适配 TypeA/TypeB)   │                    │
│  └──────────┬──────────┘                    │
│             ▼                                │
│  Socket message → adbclawd                   │
└─────────────────────────────────────────────┘

        │ ADB forward socket
        ▼

┌─────────────────────────────────────────────┐
│            adbclawd (设备端)                  │
│                                              │
│  收到 event 序列                              │
│              │                               │
│              ▼                               │
│  open("/dev/input/eventX", O_WRONLY)         │
│  write(events, sizeof(input_event) * N)      │
│  (一次 write 调用，批量写入)                    │
└─────────────────────────────────────────────┘
```

**sendevent 协议细节：**

连接时 `adbclawd` 自动探测并上报：

```json
{
  "touchscreen": {
    "device": "/dev/input/event2",
    "protocol": "B",
    "abs_mt_position_x": { "min": 0, "max": 1079 },
    "abs_mt_position_y": { "min": 0, "max": 2399 },
    "abs_mt_pressure":   { "min": 0, "max": 255 },
    "abs_mt_touch_major": { "min": 0, "max": 255 },
    "max_slots": 10
  }
}
```

Mac 端据此生成正确的 event 序列。

### 4.3 屏幕与 UI 观测

**截屏：**

| 方案 | 延迟 | 说明 |
|------|------|------|
| adbclawd socket 截屏 | ~80ms | adbclawd 内部执行 screencap，通过 socket 直传回 Mac |
| `adb exec-out screencap -p` | ~200ms | 走 ADB stdout pipe |
| scrcpy server (可选) | ~30ms | H.264 编码视频流，适合高频截屏场景 |

默认用方案 1（adbclawd socket）。如果 agent 需要高频截屏（如看视频），可启用 scrcpy 流模式。

**UI 树：**

```bash
# adbclawd 内部执行：
uiautomator dump /data/local/tmp/adbclaw_ui.xml
# 读取 XML，解析成结构化 JSON 返回
```

adbclaw 在 Mac 端将 XML 解析为带编号的元素列表：

```json
{
  "elements": [
    {
      "index": 0,
      "class": "android.widget.TextView",
      "resource_id": "com.example:id/title",
      "text": "Settings",
      "content_desc": "",
      "bounds": { "left": 0, "top": 120, "right": 1080, "bottom": 180 },
      "center": { "x": 540, "y": 150 },
      "clickable": true,
      "scrollable": false
    }
  ],
  "focused_app": "com.android.settings/.Settings",
  "screen": { "width": 1080, "height": 2400 }
}
```

**合并观测命令 `observe`：**

借鉴 mobile-use 并行获取截屏 + UI 树的做法：

```bash
adbclaw observe
# 同时返回截屏（base64）和 UI 树，减少 agent 调用次数
```

```json
{
  "screenshot": "base64://iVBORw0KGgo...",
  "screenshot_format": "jpeg",
  "ui": { "elements": [...], "focused_app": "...", "screen": {...} },
  "device": { "battery": 85, "wifi": true, "time": "14:30" },
  "timestamp": "2026-02-21T14:30:00Z"
}
```

### 4.4 文字输入

两个竞品都遇到的问题：`adb shell input text` 不支持 Unicode、不支持特殊字符。

DroidRun 的方案：Portal APK 的自定义 IME（需装 APK）。
mobile-use 的方案：UIAutomator2 的 FastInputIME（需 UIAutomator2 server）。

**adbclaw 的方案：ADB broadcast**

```bash
# 利用 Android 原生 am broadcast 传递文本
adb shell am broadcast -a ADB_INPUT_TEXT --es text "你好世界🌍"
```

这需要 adbclawd 注册一个 BroadcastReceiver——但那又需要 APK。

**务实方案：分层处理**

```
文字输入
├── ASCII 文本: adb shell input text "hello" (直接，无需额外依赖)
├── Unicode 文本: adb shell "echo '你好' | /data/local/tmp/adbclawd ime"
│   (adbclawd 通过 /dev/uinput 模拟 HID 键盘输入)
└── 剪贴板方案: adb shell "echo '你好' > /data/local/tmp/cb.txt"
    + adbclawd 模拟 Ctrl+V 粘贴
```

首选剪贴板方案——最简单且覆盖所有字符：
1. 写文本到临时文件
2. `adb shell am broadcast` 设置剪贴板内容（Android 10+ 限制需在前台）
3. sendevent 模拟长按触发粘贴，或模拟 Ctrl+V

### 4.5 App 管理

从两个竞品中提取的最小必要操作集：

| 命令 | 实现 | 说明 |
|------|------|------|
| `app list` | `pm list packages -3` | 列出第三方 App |
| `app list --all` | `pm list packages` | 包含系统 App |
| `app current` | `dumpsys window \| grep mCurrentFocus` | 当前前台 Activity |
| `app launch <pkg>` | `monkey -p <pkg> -c android.intent.category.LAUNCHER 1` | 启动 App |
| `app launch <pkg>/<act>` | `am start -n <pkg>/<act>` | 启动指定 Activity |
| `app stop <pkg>` | `am force-stop <pkg>` | 停止 App |
| `app install <apk>` | `pm install -r <apk>` | 安装 APK |
| `app uninstall <pkg>` | `pm uninstall <pkg>` | 卸载 App |
| `app info <pkg>` | `dumpsys package <pkg>` | App 详情 |

DroidRun 的 `open_app` 需要一个 LLM 调用来匹配 app name → package name。adbclaw 不做这件事——agent 自行决定 package name，adbclaw 只接受精确的 package name。

---

## 5. CLI 命令设计（定稿）

### 5.1 全局选项

```bash
adbclaw [global-options] <command> [command-options]

Global Options:
  -s, --serial <id>      目标设备 serial（默认使用唯一连接设备）
  -o, --output <format>  输出格式: json (默认) | text | quiet
  --stealth <level>      反侦测级别: off (默认) | on | deep
  --humanize             启用人类化输入模拟
  --timeout <ms>         命令超时（默认 30000ms）
  --verbose              调试输出到 stderr
```

### 5.2 命令总览

```bash
# ── 设备 ──
adbclaw device list                        # 列出已连接设备
adbclaw device info                        # 设备详情 (型号/分辨率/Android版本/电量)
adbclaw device connect <ip:port>           # 无线 ADB 连接
adbclaw device disconnect [ip:port]        # 断开无线连接
adbclaw device probe                       # 探测设备能力 (触屏协议/节点/参数范围)
adbclaw device shell <command>             # 执行 shell 命令

# ── 观测 ──
adbclaw observe                            # 截屏 + UI树 + 设备状态 (合并)
adbclaw screenshot [-o file] [--quality 80] [--format jpeg|png]
adbclaw ui tree                            # UI 元素树 (带编号)
adbclaw ui find --text "Login"             # 查找元素
adbclaw ui find --id "btn_submit"          # 按 resource-id 查找
adbclaw ui find --index 3                  # 获取编号为 3 的元素详情

# ── 输入 ──
adbclaw tap <x> <y>                        # 点击坐标
adbclaw tap --index <n>                    # 点击 UI 树中编号为 n 的元素
adbclaw tap --id <resource-id>             # 点击指定 resource-id 的元素
adbclaw tap --text "Login"                 # 点击包含文字的元素
adbclaw long-press <x> <y> [--duration 1000]
adbclaw swipe <x1> <y1> <x2> <y2> [--duration 300]
adbclaw type <text>                        # 输入文字
adbclaw key <HOME|BACK|ENTER|TAB|DEL|POWER|VOLUME_UP|VOLUME_DOWN>
adbclaw gesture <file.json>                # 执行复杂手势序列

# ── App ──
adbclaw app list [--all]                   # 列出已安装 App
adbclaw app current                        # 当前前台 App/Activity
adbclaw app launch <package>               # 启动 App
adbclaw app stop <package>                 # 停止 App
adbclaw app install <path.apk>             # 安装 APK
adbclaw app uninstall <package>            # 卸载 App

# ── 文件 ──
adbclaw file push <local> <remote>         # 推送文件到设备
adbclaw file pull <remote> <local>         # 从设备拉取文件

# ── 工具 ──
adbclaw note set <key> <value>             # 存储 kv 数据 (跨命令持久)
adbclaw note get <key>                     # 读取 kv 数据
adbclaw note list                          # 列出所有 key
adbclaw note clear                         # 清空

# ── 管理 ──
adbclaw version                            # 版本信息
adbclaw doctor                             # 环境检查 (ADB可用/设备连接/adbclawd状态)
adbclaw daemon start                       # 启动 adbclawd (通常自动)
adbclaw daemon stop                        # 停止 adbclawd
adbclaw skill                              # 输出 skill.json (供 agent 读取)
```

### 5.3 输出格式示例

所有命令统一 JSON envelope：

```json
{
  "ok": true,
  "command": "tap",
  "data": {
    "x": 540,
    "y": 150,
    "method": "sendevent",
    "element": {
      "index": 3,
      "text": "Login",
      "resource_id": "com.example:id/btn_login"
    }
  },
  "duration_ms": 45,
  "timestamp": "2026-02-21T14:30:00.123Z"
}
```

错误输出：

```json
{
  "ok": false,
  "command": "tap",
  "error": {
    "code": "ELEMENT_NOT_FOUND",
    "message": "No element found with text 'Login'",
    "suggestion": "Try 'adbclaw ui tree' to see available elements"
  }
}
```

---

## 6. Skill 文件设计

adbclaw 不嵌入 agent，但需要告诉 agent 自己能做什么。通过 `adbclaw skill` 输出标准化的 skill 描述文件：

```json
{
  "name": "adbclaw",
  "version": "0.1.0",
  "description": "Android device control CLI with anti-detection input injection",
  "capabilities": {
    "input_methods": ["sendevent", "adb_input"],
    "stealth_levels": ["off", "on", "deep"],
    "humanize": true,
    "unicode_input": true
  },
  "tools": [
    {
      "name": "observe",
      "description": "Capture screenshot + UI element tree + device state in one call. Use this as the primary way to understand what's on screen.",
      "usage": "adbclaw observe [--quality 80]",
      "returns": "screenshot (base64 JPEG), UI elements with index/id/text/bounds, focused app, screen dimensions",
      "when_to_use": "Before deciding any action. After every action to verify the result."
    },
    {
      "name": "tap",
      "description": "Tap on a UI element or coordinate. Prefer --index (from observe) for reliability.",
      "usage": "adbclaw tap {--index N | --id RESOURCE_ID | --text TEXT | X Y}",
      "parameters": {
        "index": "Element index from 'observe' output (most reliable)",
        "id": "Android resource-id string",
        "text": "Visible text on the element",
        "x,y": "Raw pixel coordinates (least reliable, use as fallback)"
      },
      "when_to_use": "To click buttons, links, input fields, or any interactive element."
    },
    {
      "name": "long_press",
      "description": "Long press on a coordinate or element.",
      "usage": "adbclaw long-press {--index N | X Y} [--duration 1000]",
      "when_to_use": "For context menus, drag initiation, or special long-press actions."
    },
    {
      "name": "swipe",
      "description": "Swipe from one point to another.",
      "usage": "adbclaw swipe X1 Y1 X2 Y2 [--duration 300]",
      "when_to_use": "To scroll content, dismiss notifications, or navigate between pages."
    },
    {
      "name": "type",
      "description": "Input text into the currently focused field.",
      "usage": "adbclaw type TEXT",
      "when_to_use": "After tapping an input field to focus it."
    },
    {
      "name": "key",
      "description": "Press a system key.",
      "usage": "adbclaw key {HOME|BACK|ENTER|TAB|DEL|POWER|VOLUME_UP|VOLUME_DOWN}",
      "when_to_use": "To navigate back, go home, confirm input, or control device."
    },
    {
      "name": "app_launch",
      "description": "Launch an app by package name.",
      "usage": "adbclaw app launch PACKAGE_NAME",
      "when_to_use": "To open a specific app. Use 'app list' first to find the package name."
    },
    {
      "name": "app_stop",
      "description": "Force stop an app.",
      "usage": "adbclaw app stop PACKAGE_NAME",
      "when_to_use": "To close/reset an app."
    },
    {
      "name": "app_list",
      "description": "List installed apps with package names.",
      "usage": "adbclaw app list [--all]",
      "when_to_use": "To find the package name of an app before launching it."
    },
    {
      "name": "app_current",
      "description": "Get the currently focused app and activity.",
      "usage": "adbclaw app current",
      "when_to_use": "To check which app is in the foreground."
    },
    {
      "name": "ui_tree",
      "description": "Get the UI element tree with indices, resource-ids, text, and bounds.",
      "usage": "adbclaw ui tree",
      "when_to_use": "When you need detailed element information beyond what 'observe' provides."
    },
    {
      "name": "ui_find",
      "description": "Find specific UI elements by text or resource-id.",
      "usage": "adbclaw ui find {--text TEXT | --id RESOURCE_ID}",
      "when_to_use": "To locate a specific element before interacting with it."
    },
    {
      "name": "screenshot",
      "description": "Capture a screenshot.",
      "usage": "adbclaw screenshot [--quality 80] [--format jpeg|png] [-o file]",
      "when_to_use": "When you need only the visual state without UI tree."
    },
    {
      "name": "note",
      "description": "Persistent key-value scratchpad for cross-step data.",
      "usage": "adbclaw note {set KEY VALUE | get KEY | list | clear}",
      "when_to_use": "To remember information across multiple steps (e.g., a price found in one app to enter in another)."
    },
    {
      "name": "device_info",
      "description": "Get device details: model, screen size, Android version, battery, network.",
      "usage": "adbclaw device info",
      "when_to_use": "At the start of a session to understand the device."
    }
  ],
  "workflow_hints": {
    "observe_first": "Always run 'observe' before deciding an action and after performing an action to verify results.",
    "prefer_index": "When tapping, prefer --index over coordinates. Indices come from the 'observe' output.",
    "type_after_focus": "Always tap an input field first, then use 'type' to enter text.",
    "scroll_pattern": "To scroll down: swipe from bottom-center to top-center (e.g., swipe 540 1800 540 600).",
    "error_handling": "If an action fails, re-observe to see the current state before retrying."
  }
}
```

### Agent 集成方式

**方式 1：CLI 子进程调用（最简单）**

Agent 直接 spawn adbclaw 子进程，读 JSON stdout：

```
Agent loop:
  1. state = exec("adbclaw observe")
  2. LLM(state.screenshot, state.ui) → decision
  3. exec("adbclaw tap --index 3")
  4. goto 1
```

**方式 2：MCP Server（标准化）**

```bash
adbclaw mcp serve --port 8400
# 暴露所有命令为 MCP tools
# Claude Desktop / OpenClaw 通过 MCP 协议调用
```

**方式 3：长连接模式（高性能）**

```bash
adbclaw interactive
# stdin 接收 JSON 命令，stdout 返回 JSON 结果
# 省去每次子进程启动开销
{"cmd": "observe"}
{"ok": true, "data": {...}}
{"cmd": "tap", "args": {"index": 3}}
{"ok": true, "data": {...}}
```

---

## 7. 技术选型

| 组件 | 选型 | 理由 |
|------|------|------|
| CLI 主体 | **Go** | 项目技术栈一致，单 binary，交叉编译到 macOS/Linux |
| CLI 框架 | **cobra** | Go 生态标准 CLI 框架 |
| ADB 通信 | **go-adb** 或自实现 ADB 协议 | 避免依赖 adb binary，更可控 |
| 设备端程序 | **C** (NDK 交叉编译) | 直接操作 `/dev/input`，性能关键，体积小 |
| JSON 输出 | **encoding/json** | 标准库足够 |
| UI XML 解析 | **encoding/xml** | 解析 uiautomator dump 输出 |
| MCP Server | **Go** (JSON-RPC over stdio) | 与 CLI 同进程，MCP 协议简单 |

### 目录结构

```
adbclaw/
├── cmd/                    # CLI 命令定义 (cobra)
│   ├── root.go
│   ├── device.go           # device list/info/connect/probe
│   ├── observe.go          # observe (截屏+UI树)
│   ├── input.go            # tap/swipe/long-press/type/key
│   ├── app.go              # app list/current/launch/stop
│   ├── ui.go               # ui tree/find
│   ├── file.go             # file push/pull
│   ├── note.go             # note set/get/list
│   ├── skill.go            # skill (输出 skill.json)
│   ├── mcp.go              # mcp serve
│   └── doctor.go           # doctor
│
├── pkg/
│   ├── adb/                # ADB 协议实现
│   │   ├── transport.go    # USB/TCP 连接
│   │   ├── device.go       # 设备操作
│   │   ├── shell.go        # shell 命令执行
│   │   └── forward.go      # 端口转发
│   │
│   ├── input/              # 输入注入引擎
│   │   ├── engine.go       # 输入引擎接口
│   │   ├── sendevent.go    # sendevent 实现
│   │   ├── adbinput.go     # adb shell input 实现 (fallback)
│   │   ├── humanize.go     # 人类化模拟 (抖动/曲线/压力)
│   │   └── probe.go        # 触屏参数探测
│   │
│   ├── observe/            # 屏幕与 UI 观测
│   │   ├── screenshot.go   # 截屏
│   │   ├── uitree.go       # UI 树解析 + 编号
│   │   └── combined.go     # observe 合并命令
│   │
│   ├── daemon/             # adbclawd 管理
│   │   ├── push.go         # 推送 binary 到设备
│   │   ├── lifecycle.go    # 启动/停止/心跳
│   │   └── protocol.go     # socket 通信协议
│   │
│   ├── stealth/            # 反侦测
│   │   ├── checker.go      # 检查当前侦测暴露面
│   │   └── mitigate.go     # 执行缓解措施
│   │
│   └── output/             # 输出格式化
│       ├── json.go
│       ├── text.go
│       └── envelope.go     # 统一 {ok, command, data, error} 封装
│
├── device/                 # adbclawd 设备端程序 (C)
│   ├── main.c
│   ├── input.c             # sendevent 批量写入
│   ├── screen.c            # 截屏
│   ├── uitree.c            # uiautomator dump 调用
│   ├── probe.c             # 设备探测
│   ├── socket.c            # Unix socket server
│   └── Makefile            # NDK 交叉编译
│
├── skill/                  # Skill 描述文件
│   └── skill.json          # Agent 读取的能力描述
│
├── go.mod
├── go.sum
├── Makefile                # 构建: make build / make device / make all
└── README.md
```

---

## 8. 实现优先级

### Phase 1 — 基础可用（MVP）

最小可用版本，agent 能完成基本交互循环：

| 命令 | 优先级 | 说明 |
|------|--------|------|
| `device list` | P0 | 发现设备 |
| `device info` | P0 | 获取屏幕尺寸 |
| `screenshot` | P0 | 截屏（先用 `adb exec-out screencap -p`，不需要 adbclawd） |
| `ui tree` | P0 | UI 树（先用 `adb shell uiautomator dump`） |
| `observe` | P0 | 合并截屏 + UI 树 |
| `tap X Y` | P0 | 坐标点击（先用 `adb shell input tap`） |
| `tap --index N` | P0 | 元素编号点击 |
| `swipe` | P0 | 滑动 |
| `type` | P0 | 文字输入 |
| `key` | P0 | 按键 |
| `app list/current/launch/stop` | P0 | App 管理 |
| `skill` | P0 | 输出 skill.json |

Phase 1 **不需要** adbclawd，全部通过 `adb shell` 命令实现。

### Phase 2 — sendevent 引擎

| 功能 | 说明 |
|------|------|
| adbclawd C 程序 | 设备端 binary，socket 通信 |
| `device probe` | 探测触屏参数 |
| sendevent tap/swipe | 通过 `/dev/input` 注入 |
| `--stealth on` | 默认使用 sendevent |
| 自动 fallback | sendevent 失败时降级到 `adb shell input` |

### Phase 3 — 人类化 + 高级功能

| 功能 | 说明 |
|------|------|
| `--humanize` | 高斯抖动 + 贝塞尔曲线 + 压力模拟 |
| `gesture` | 复杂手势序列 |
| `note` | 持久化 kv 存储 |
| `mcp serve` | MCP Server |
| `interactive` | 长连接模式 |
| Unicode 文字输入 | 剪贴板方案 |
| 无线 ADB 优先 | `device connect` 智能选择 |

### Phase 4 — deep-stealth + 高级观测

| 功能 | 说明 |
|------|------|
| `--stealth deep` | Root 设备的深度反侦测 |
| scrcpy 截屏流 | 高频截屏场景 |
| `stealth status` | 侦测暴露面检查 |

---

## 9. 与 Agent 的典型交互流程

```
Agent (OpenClaw)                          adbclaw CLI
      │                                       │
      │  exec: adbclaw device info             │
      │──────────────────────────────────────►  │
      │  ◄─ {model, screen: 1080x2400, ...}   │
      │                                        │
      │  exec: adbclaw observe --quality 60    │
      │──────────────────────────────────────►  │
      │  ◄─ {screenshot: base64, ui: [...]}    │
      │                                        │
      │  [LLM 分析截图+UI树 → 决定点击 Login]    │
      │                                        │
      │  exec: adbclaw tap --index 5           │
      │──────────────────────────────────────►  │
      │  ◄─ {ok: true, x: 540, y: 1200}       │
      │                                        │
      │  exec: adbclaw observe                 │
      │──────────────────────────────────────►  │
      │  ◄─ {screenshot: ..., ui: [...]}       │
      │                                        │
      │  [LLM 看到登录页 → 输入用户名]           │
      │                                        │
      │  exec: adbclaw tap --id "input_email"  │
      │──────────────────────────────────────►  │
      │  ◄─ {ok: true}                         │
      │                                        │
      │  exec: adbclaw type "user@example.com" │
      │──────────────────────────────────────►  │
      │  ◄─ {ok: true}                         │
      │                                        │
      │  ... (observe → decide → act loop)     │
```

---

## 10. 关键风险与对策

| 风险 | 影响 | 对策 |
|------|------|------|
| `uiautomator dump` 速度慢（~1-2s） | Agent 交互循环变慢 | Phase 2 用 adbclawd 异步 dump；考虑缓存上一帧 + diff |
| 某些设备 `/dev/input` 无 shell 写权限 | sendevent 不可用 | 自动检测权限，降级到 `adb shell input` |
| `uiautomator dump` 在动画中可能失败 | UI 树为空 | 重试机制 + 返回上一次成功的 UI 树 |
| 不同设备 sendevent 协议差异大 | 适配工作量 | `device probe` 自动探测 + 内置多种协议模板 |
| Unicode 文字输入的剪贴板方案在 Android 10+ 受限 | 后台 App 无法写剪贴板 | adbclawd 作为前台 shell 进程执行，不受限制 |

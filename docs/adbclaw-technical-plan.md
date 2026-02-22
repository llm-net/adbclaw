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
- 内置 sendevent 反侦测输入引擎，支持三级降级
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
| 语言 | Python | Python | Go + Java（设备端） |
| LLM 依赖 | llama-index (必须) | LangGraph + LangChain (必须) | **无** |
| 输入注入 | `adb shell input`（可检测） | `adb shell input`（可检测） | **sendevent /dev/input**（高隐蔽，三级降级） |
| 反侦测 | 无 | 无 | **stealth / deep-stealth** |
| 手机端依赖 | Portal APK + AccessibilityService | UIAutomator2 server | **无 APK**，推送 JAR 通过 app_process 运行 |
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
- 文本输入走 Portal IME — 改用 ADB `input text` + 对特殊字符的转义处理，或用剪贴板方案

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

### 2.3 从 scrcpy 学到的

**采纳：**
- **app_process 启动方式** — 推送 JAR 到设备，通过 `app_process` 运行 Java 代码，兼容性极高（scrcpy 已大规模验证）
- **SurfaceControl 反射截屏** — 借鉴其按 Android 版本分支反射的技巧
- **剪贴板管理** — 借鉴其 ClipboardManager 反射方式

**不采纳：**
- InputManager.injectInputEvent() 输入注入 — 该方式 deviceId=-1，可被 App 检测，与 adbclaw 反侦测目标冲突
- 视频流编码/音频采集/摄像头 — MVP 不需要，后续可选集成
- scrcpy-server 整体 — 代码量 ~15000 行，adbclaw 只需 ~10% 的功能，自己写更可控

### 2.4 竞品共同缺陷 → adbclaw 的机会

| 缺陷 | adbclaw 的解决方案 |
|------|-------------------|
| 输入用 `adb shell input`，deviceId=-1 可检测 | sendevent 引擎，deviceId 真实 |
| 无时间抖动/轨迹模拟 | 人类化引擎（高斯抖动 + 贝塞尔曲线） |
| 依赖 Python 运行时 | Go 单 binary |
| 需要安装 APK / UIAutomator2 | 推送 JAR 到 `/data/local/tmp/`，通过 app_process 运行，不装 APK |
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
│  │           Connection Strategy Layer                 │ │
│  │  ┌────────────┐ ┌────────────┐ ┌────────────────┐  │ │
│  │  │ 连接探测    │ │ 能力协商    │ │ 通道选择       │  │ │
│  │  │ USB/WiFi   │ │ adbclawd   │ │ 传输/输入/截屏  │  │ │
│  │  └────────────┘ └────────────┘ └────────────────┘  │ │
│  └────────────────────┬───────────────────────────────┘ │
└───────────────────────┼─────────────────────────────────┘
                        │ USB / WiFi TCP
┌───────────────────────▼─────────────────────────────────┐
│  Android 设备                                            │
│                                                          │
│  /data/local/tmp/adbclawd.jar (设备端 Java 服务)          │
│  运行方式: app_process -Djava.class.path=... / Server    │
│  ├── 触摸注入 (/dev/input write 或 InputManager 降级)     │
│  ├── 截屏 (SurfaceControl 反射 → JPEG 压缩)              │
│  ├── UI 树 (uiautomator dump + XML 解析)                 │
│  ├── 设备能力探测 (getevent -pl 解析)                     │
│  └── 设备信息查询                                         │
│                                                          │
│  触屏驱动 (/dev/input/eventX) ← 内核自带，无需安装         │
└─────────────────────────────────────────────────────────┘
```

---

## 4. 核心模块设计

### 4.1 设备端辅助程序 `adbclawd`

#### 为什么需要设备端程序

DroidRun 每次 tap 都要启动一个 `adb shell input tap` 进程（~50ms 开销）。sendevent 更糟——一次 tap 需要 6-8 条 sendevent 命令，每条都是一个进程（~400ms）。

adbclawd 是一个**常驻在手机端的 Java 服务**，通过 ADB forward 的 socket 与 Mac 端通信，消除进程启动开销：

```
Mac (adbclaw) ──ADB forward──> Android (adbclawd: LocalSocket)

一次 tap 耗时对比:
  adb shell input tap      ~65ms   (每次启动进程)
  adb shell sendevent ×8   ~400ms  (8次进程启动)
  adbclawd socket write()  ~5ms    (常驻进程，直接 write)
```

#### 为什么用 Java JAR 而不是 C binary

原方案设计 adbclawd 为 C 语言编写、NDK 交叉编译的 native binary。经过分析，**改为 Java JAR + app_process 方式**：

| | C binary (原方案) | Java JAR (现方案) |
|---|---|---|
| 设备兼容性 | 中等 — SELinux 可能阻止 native 执行 | **极高** — scrcpy 已在数百款设备验证 |
| 多架构 | 需编译 arm64/arm32/x86 三个版本 | **无需** — Java 字节码跨架构 |
| JPEG 截屏 | 需嵌入 stb_image 等第三方库 | **系统内置** — `Bitmap.compress(JPEG)` |
| 截屏方式 | popen("screencap") 外部命令 | **SurfaceControl 反射** — 更快更直接 |
| 写 /dev/input | open() + write() | FileOutputStream — 同样可行 |
| 启动速度 | 即时 | ~200ms（一次性开销，启动后常驻） |

**关键认知**：Android 原生 `screencap` 命令**不支持 JPEG 输出**（只有 PNG 和 raw），JPEG 压缩必须由 adbclawd 自己完成。Java 方案天然拥有 `Bitmap.compress(JPEG)` 能力，无需引入第三方库。

#### adbclawd 运行方式

```bash
# 推送 JAR 到设备
adb push adbclawd.jar /data/local/tmp/

# 通过 app_process 启动（和 scrcpy 相同的方式）
adb shell CLASSPATH=/data/local/tmp/adbclawd.jar \
    app_process / com.adbclaw.Server
```

`app_process` 是 Android 系统自带的 Java 进程启动器，每台设备都有。

#### adbclawd 职责

```
adbclawd.jar (Java, 通过 app_process 运行)
├── /input     — 触摸注入 (sendevent 或 InputManager 降级)
├── /probe     — 探测 /dev/input 设备节点、触屏参数、写权限
├── /screen    — SurfaceControl 截屏 → JPEG 压缩 → socket 传回
├── /uitree    — uiautomator dump + XML 解析 → JSON 返回
├── /clipboard — 剪贴板读写 (用于 Unicode 文字输入)
└── /info      — 设备信息 (分辨率/DPI/Android 版本/电量)
```

#### 为什么不装 APK

- APK 安装需要用户确认，自动化场景不友好
- 已安装 App 列表可被其他 App 枚举（`PackageManager.getInstalledPackages`）
- AccessibilityService 可被检测
- `/data/local/tmp` 下的文件对普通 App 不可见

### 4.2 输入注入引擎

#### 三级降级策略

不同设备对 `/dev/input` 写入的支持程度不同。adbclawd 在启动时**自动探测**设备能力，选择最优注入方式，并支持自动降级：

```
adbclawd 启动时自动探测:

  尝试 open("/dev/input/eventX", O_WRONLY)
       │
       ├── 成功 → Level 1: sendevent 模式 (高隐蔽)
       │          ✓ deviceId 真实 (不是 -1)
       │          ✓ getSource() = SOURCE_TOUCHSCREEN
       │          ✓ 无 FLAG_IS_ACCESSIBILITY_EVENT
       │          ✓ 可自定义压力值和触摸面积
       │
       └── 失败 (Permission denied / SELinux)
               │
               ├── 尝试 InputManager.injectInputEvent() 反射
               │       │
               │       ├── 成功 → Level 2: inject 模式 (中隐蔽)
               │       │          ✗ deviceId = -1 (可被检测)
               │       │          ✓ 配合人类化时间/轨迹模拟仍有价值
               │       │          (与 scrcpy 同等水平)
               │       │
               │       └── 失败 → Level 3: adb shell input (低隐蔽)
               │                  ✗ 纯命令行调用，兜底方案
               │
               └── 上报实际能力给 Mac 端
```

连接建立后上报的 `ConnectionProfile`：

```json
{
  "transport": "usb",
  "daemon_available": true,
  "input": {
    "method": "sendevent",
    "stealth_level": "high",
    "touchscreen": {
      "device": "/dev/input/event2",
      "protocol": "B",
      "abs_mt_position_x": { "min": 0, "max": 1079 },
      "abs_mt_position_y": { "min": 0, "max": 2399 },
      "abs_mt_pressure":   { "min": 0, "max": 255 },
      "abs_mt_touch_major": { "min": 0, "max": 255 },
      "max_slots": 10
    }
  },
  "screen": {
    "method": "surfacecontrol",
    "width": 1080,
    "height": 2400,
    "density": 440
  }
}
```

如果 `/dev/input` 不可写，profile 中 `input.method` 为 `"inject"` 或 `"adb_input"`，`stealth_level` 相应降低。

#### 设备兼容性预期

| 品牌 | app_process | SurfaceControl 截屏 | /dev/input 写入 | 说明 |
|------|:-----------:|:-------------------:|:---------------:|------|
| Pixel/AOSP | ✅ | ✅ | ✅ | 基线设备，完全支持 |
| 小米 HyperOS/MIUI | ✅ | ✅ | ✅ | 主流设备，通常支持 |
| OPPO ColorOS | ✅ | ✅ | ✅ | 通常支持 |
| Vivo OriginOS | ✅ | ✅ | ✅ | 通常支持 |
| 三星 One UI | ✅ | ✅ | ⚠️ | Knox 企业模式可能限制 |
| 华为 HarmonyOS | ✅ | ✅ | ⚠️ | 额外安全策略可能限制 |

**app_process + SurfaceControl**：99%+ 设备支持（scrcpy 已验证）。
**/dev/input 写入**：80-90% 设备支持（minitouch 有验证，但范围不如 scrcpy 广）。
不支持的设备**自动降级**到 Level 2/3，不会导致功能不可用。

#### 触屏驱动说明

手机触屏驱动是**内核自带**的，不需要安装任何驱动。`/dev/input/eventX` 设备节点由内核触屏驱动自动创建。不同手机的差异不是"驱动不同需要安装"，而是**参数不同需要探测**：

| 差异项 | 示例 | 探测方式 |
|-------|------|---------|
| 设备节点路径 | Pixel: event2, 小米: event4 | `getevent -pl` 找含 `ABS_MT_POSITION_X` 的设备 |
| 坐标范围 | 0-1079 或 0-32767 | 读取 `ABS_MT_POSITION_X/Y` 的 min/max |
| 压力范围 | 0-255 或 0-63，部分设备无此轴 | 读取 `ABS_MT_PRESSURE` 的 min/max |
| 协议类型 | Type A (旧) 或 Type B (主流) | 检查是否有 `ABS_MT_SLOT` 事件 |
| 写权限 | 大部分 shell 用户可写 | 尝试 open() 验证 |

全部通过 `getevent -pl` 一条命令 + 尝试 open() 自动完成，**零人工配置**。

#### 输入引擎架构

```
┌─────────────────────────────────────────────┐
│            Input Engine (Mac 端 Go)          │
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
│  │ 根据 probe 结果生成   │                    │
│  │ sendevent 序列       │                    │
│  │ (适配 TypeA/TypeB)   │                    │
│  └──────────┬──────────┘                    │
│             ▼                                │
│  Socket message → adbclawd                   │
└─────────────────────────────────────────────┘

        │ ADB forward socket
        ▼

┌─────────────────────────────────────────────┐
│            adbclawd (设备端 Java)             │
│                                              │
│  收到 event 序列                              │
│       │                                      │
│       ├── Level 1: FileOutputStream 写入      │
│       │   /dev/input/eventX (批量 write)      │
│       │                                      │
│       ├── Level 2: InputManager 反射注入      │
│       │   (降级方案，deviceId=-1)              │
│       │                                      │
│       └── Level 3: Runtime.exec("input tap") │
│           (兜底方案)                           │
└─────────────────────────────────────────────┘
```

### 4.3 连接策略层

连接方式直接影响效率和隐蔽性，需要作为一个整体考虑。

#### 连接方式对比

| | USB ADB | WiFi ADB | 说明 |
|---|---|---|---|
| 延迟 | ~15-30ms/命令 | ~5-100ms/命令 (波动大) | WiFi 受网络质量影响 |
| 稳定性 | 高 | 中 | WiFi 可能断连 |
| 截屏 (1080p PNG ~5MB) | ~100ms | ~300ms+ | 带宽瓶颈 |
| 截屏 (JPEG ~200KB) | ~20ms | ~50ms | adbclawd JPEG 压缩消除带宽差 |
| 隐蔽性 | 低 — `sys.usb.state` 含 "adb" | **高** — 不暴露 USB 状态 | 反侦测关键差异 |

#### 自适应连接策略

```
adbclaw 连接时根据 --stealth 级别自动选择:

  stealth=off   → 优先 USB（最快最稳）
  stealth=on    → 优先 WiFi（隐藏 USB 状态）
                   WiFi 延迟 > 100ms 时警告用户
  stealth=deep  → 强制 WiFi + Root 级隐藏

  无论底层传输是 USB 还是 WiFi:
    → ADB forward 建立 socket 隧道
    → 后续操作走 adbclawd socket，传输层差异被屏蔽
    → adbclawd 端 JPEG 压缩进一步减少 WiFi 带宽影响
```

#### 效率对比总表

| 操作 | 无 adbclawd (adb shell) | 有 adbclawd (socket) | 提升 |
|------|------------------------|---------------------|------|
| 单次 tap | ~65ms | ~5ms | **13x** |
| 20 点 swipe | ~300ms | ~10ms | **30x** |
| 截屏 (USB) | ~200ms (PNG) | ~80ms (JPEG) | **2.5x** |
| 截屏 (WiFi) | ~500ms (PNG) | ~120ms (JPEG) | **4x** |
| UI 树 dump | ~1500ms | ~1200ms | 1.2x (瓶颈在 uiautomator 自身) |

### 4.4 屏幕与 UI 观测

#### 截屏

Android 原生 `screencap` 命令**只支持 PNG 和 raw 格式**，不支持 JPEG。adbclawd 通过 Java 运行时内置的 `Bitmap.compress()` 实现 JPEG 压缩：

```java
// adbclawd 截屏流程
// 1. SurfaceControl 反射获取 Bitmap（借鉴 scrcpy 的方式）
Bitmap bmp = SurfaceControl.screenshot(displayId);

// 2. Android 内置 JPEG 编码，无需第三方库
ByteArrayOutputStream out = new ByteArrayOutputStream();
bmp.compress(Bitmap.CompressFormat.JPEG, quality, out);

// 3. 通过 socket 发回 Mac 端 (~200KB vs PNG ~5MB)
socket.write(out.toByteArray());
```

SurfaceControl 反射需要按 Android 版本适配（scrcpy 已验证各版本的反射路径）：

| Android 版本 | 反射路径 | 说明 |
|-------------|---------|------|
| 5.0 - 8.1 | `SurfaceControl.screenshot(display, w, h, rotation)` | 早期 API |
| 9.0 (P) | 引入 hidden API 限制，但 shell/app_process 豁免 | 无影响 |
| 10 - 11 | API 签名变化，需要不同的反射参数 | scrcpy 已处理 |
| 12+ | SurfaceControl 部分重构 | scrcpy 已处理 |
| 14+ | 进一步重构 | scrcpy 已处理 |

**截屏方案降级**：

| 方案 | 延迟 | 条件 |
|------|------|------|
| adbclawd SurfaceControl + JPEG | ~80ms | adbclawd 运行中（默认） |
| `adb exec-out screencap -p` (PNG) | ~200ms | adbclawd 不可用时的 fallback |

#### UI 树

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

#### 合并观测命令 `observe`

借鉴 mobile-use 并行获取截屏 + UI 树的做法：

```bash
adbclaw observe
# 同时返回截屏（base64 JPEG）和 UI 树，减少 agent 调用次数
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

### 4.5 文字输入

两个竞品都遇到的问题：`adb shell input text` 不支持 Unicode、不支持特殊字符。

**adbclaw 的方案：分层处理**

```
文字输入
├── ASCII 文本: adb shell input text "hello" (直接，无需额外依赖)
├── Unicode 文本: 剪贴板方案
│   1. adbclawd 通过 ClipboardManager 反射设置剪贴板内容
│   2. sendevent 模拟 Ctrl+V 粘贴
│   (adbclawd 作为 app_process 进程可访问 ClipboardManager)
└── 兜底: adb shell input text (仅 ASCII)
```

### 4.6 App 管理

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

adbclaw 不做 app name → package name 的 LLM 匹配（DroidRun 做了），只接受精确的 package name。

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
adbclaw device probe                       # 探测设备能力 (触屏协议/节点/参数范围/写权限)
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
adbclaw doctor                             # 环境检查 (ADB/设备/adbclawd/能力)
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
    "input_methods": ["sendevent", "inject", "adb_input"],
    "stealth_levels": ["off", "on", "deep"],
    "humanize": true,
    "unicode_input": true
  },
  "tools": [
    {
      "name": "observe",
      "description": "Capture screenshot + UI element tree + device state in one call.",
      "usage": "adbclaw observe [--quality 80]",
      "returns": "screenshot (base64 JPEG), UI elements with index/id/text/bounds, focused app, screen dimensions",
      "when_to_use": "Before deciding any action. After every action to verify the result."
    },
    {
      "name": "tap",
      "description": "Tap on a UI element or coordinate. Prefer --index (from observe) for reliability.",
      "usage": "adbclaw tap {--index N | --id RESOURCE_ID | --text TEXT | X Y}",
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
      "name": "app_current",
      "description": "Get the currently focused app and activity.",
      "usage": "adbclaw app current",
      "when_to_use": "To check which app is in the foreground."
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
| 设备端程序 | **Java** (app_process 启动) | scrcpy 验证的兼容性路径，内置 JPEG/SurfaceControl/Clipboard |
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
│   │   ├── engine.go       # 输入引擎接口 (Level 1/2/3)
│   │   ├── sendevent.go    # sendevent 序列生成 (Level 1)
│   │   ├── inject.go       # InputManager 注入 (Level 2, 设备端执行)
│   │   ├── adbinput.go     # adb shell input 实现 (Level 3)
│   │   ├── humanize.go     # 人类化模拟 (抖动/曲线/压力)
│   │   └── probe.go        # 触屏参数探测 + 写权限检测
│   │
│   ├── observe/            # 屏幕与 UI 观测
│   │   ├── screenshot.go   # 截屏 (adbclawd JPEG / screencap PNG 降级)
│   │   ├── uitree.go       # UI 树解析 + 编号
│   │   └── combined.go     # observe 合并命令
│   │
│   ├── connect/            # 连接策略层
│   │   ├── strategy.go     # 连接方式选择 (USB/WiFi, stealth 适配)
│   │   ├── profile.go      # ConnectionProfile 定义
│   │   └── heartbeat.go    # adbclawd 心跳 + 自动重连
│   │
│   ├── daemon/             # adbclawd 管理
│   │   ├── push.go         # 推送 JAR 到设备
│   │   ├── lifecycle.go    # app_process 启动/停止/心跳
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
├── device/                 # adbclawd 设备端程序 (Java)
│   └── src/
│       └── com/adbclaw/
│           ├── Server.java         # 入口 + socket server + 命令分发
│           ├── ScreenCapture.java  # SurfaceControl 反射截屏 + JPEG
│           ├── TouchInjector.java  # /dev/input 写入 + InputManager 降级
│           ├── UiTreeDumper.java   # uiautomator dump + XML 解析
│           ├── DeviceProbe.java    # getevent -pl 解析 + 权限检测
│           ├── Clipboard.java      # ClipboardManager 反射
│           └── DeviceInfo.java     # 设备信息查询
│
├── skill/                  # Skill 描述文件
│   └── skill.json          # Agent 读取的能力描述
│
├── go.mod
├── go.sum
├── Makefile                # 构建: make build / make device / make all
└── README.md
```

#### adbclawd.jar 构建

```bash
# 编译 Java 源码为 DEX（Android 可执行格式）
# 不需要 Android Studio 或 Gradle，只需要 javac + d8

# 1. 编译 Java → class
javac -source 8 -target 8 device/src/com/adbclaw/*.java -d build/classes/

# 2. class → DEX (Android 字节码)
d8 build/classes/com/adbclaw/*.class --output build/

# 3. 打包
cp build/classes.dex adbclawd.jar

# 体积: ~50-100KB
```

`d8` 是 Android SDK build-tools 自带的 DEX 编译器，无需完整 Android Studio。

---

## 8. 实现优先级

### Phase 1 — 基础可用（MVP）

最小可用版本，agent 能完成基本交互循环。**全部通过 adb shell 命令实现，不需要 adbclawd**。

目标：跑通 `observe → decide → act` 循环。

| 步骤 | 内容 | 具体实现 |
|------|------|---------|
| 1.1 | Go 项目骨架 | cobra CLI + JSON envelope 输出 + 全局选项 |
| 1.2 | `device list` / `device info` | `adb devices` + `getprop` 解析 |
| 1.3 | `screenshot` | `adb exec-out screencap -p` → PNG stdout 或文件 |
| 1.4 | `ui tree` | `adb shell uiautomator dump` → XML 解析 → JSON 带编号 |
| 1.5 | `observe` | 并行执行截屏 + UI 树，合并输出 |
| 1.6 | `tap X Y` / `tap --index N` | `adb shell input tap` + 从 UI 树查坐标 |
| 1.7 | `swipe` / `key` / `type` | `adb shell input swipe/keyevent/text` |
| 1.8 | `app list/current/launch/stop` | `pm` / `dumpsys` / `monkey` / `am` 命令 |
| 1.9 | `skill` | 输出 skill.json |
| 1.10 | `doctor` | 环境检查 |

**Phase 1 交付物**：一个可工作的 `adbclaw` Go binary，agent 可以用它完成基本的 Android 控制。输入用 `adb shell input`（可被检测，但功能完整）。

### Phase 2 — adbclawd 设备端服务

引入 adbclawd.jar，大幅提升效率和隐蔽性。

| 步骤 | 内容 | 具体实现 |
|------|------|---------|
| 2.1 | adbclawd Java 服务骨架 | app_process 启动 + LocalSocket server |
| 2.2 | `device probe` | getevent -pl 解析 + /dev/input 权限检测 |
| 2.3 | sendevent 输入注入 | FileOutputStream → /dev/input (Level 1) |
| 2.4 | InputManager 降级 | 反射注入 (Level 2) |
| 2.5 | SurfaceControl 截屏 | 反射截屏 + JPEG 压缩 |
| 2.6 | Mac 端 daemon 管理 | push JAR + 启动/停止/心跳 |
| 2.7 | 连接策略层 | ConnectionProfile + 自动能力探测 |
| 2.8 | `--stealth on` | 默认使用 sendevent + 自动降级 |

**Phase 2 交付物**：adbclawd.jar 设备端服务，sendevent 高隐蔽输入，JPEG 截屏，效率提升 10-30x。

### Phase 3 — 人类化 + 高级功能

| 功能 | 说明 |
|------|------|
| `--humanize` | 高斯抖动 + 贝塞尔曲线 + 压力模拟 |
| `gesture` | 复杂手势序列 |
| `note` | 持久化 kv 存储 |
| `mcp serve` | MCP Server |
| `interactive` | 长连接模式 |
| Unicode 文字输入 | ClipboardManager 反射 + Ctrl+V |
| 无线 ADB 优先 | `device connect` + stealth 自适应 |

### Phase 4 — deep-stealth + 高级观测

| 功能 | 说明 |
|------|------|
| `--stealth deep` | Root 设备的深度反侦测 |
| scrcpy 视频流集成 | 高频截屏场景（可选） |
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
      │  ◄─ {ok: true, x: 540, y: 1200,       │
      │       method: "sendevent"}             │
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
| 某些设备 `/dev/input` 无 shell 写权限 | sendevent 不可用 | **三级降级**：sendevent → InputManager.inject → adb shell input |
| `uiautomator dump` 在动画中可能失败 | UI 树为空 | 重试机制 + 返回上一次成功的 UI 树 |
| 不同设备触屏参数差异 | 坐标/压力范围不同 | `device probe` 自动探测 getevent -pl，零人工配置 |
| SurfaceControl API 版本差异 | 截屏反射失败 | 按 Android 版本分支反射（借鉴 scrcpy），降级到 screencap |
| WiFi ADB 连接不稳定 | 操作中断 | adbclawd 心跳机制 + 自动重连 + 操作超时重试 |
| 三星 Knox / 华为安全策略 | /dev/input 或 app_process 受限 | 自动检测 + 降级 + doctor 命令提示用户 |
| Unicode 文字输入的剪贴板方案限制 | Android 10+ 后台剪贴板受限 | adbclawd 作为 app_process 进程有更高权限，不受限 |

---

## 11. Phase 1 编码计划

Phase 1 目标：**不依赖 adbclawd**，纯 adb shell 命令实现所有基础功能，跑通 agent 交互循环。

### 11.1 前置条件

- Go 1.21+
- 系统已安装 `adb` 命令（Phase 1 先依赖 adb binary，Phase 2 考虑自实现协议）
- 一台已开启 USB 调试的 Android 设备

### 11.2 编码顺序

```
Step 1: 项目初始化
  go mod init github.com/anthropics/adbclaw  (或实际 org)
  引入 cobra 依赖
  实现 root.go (全局选项 + JSON/text output envelope)

Step 2: ADB 封装层
  pkg/adb/shell.go — 封装 exec.Command("adb", "shell", ...)
  统一处理: 超时、设备选择 (-s serial)、错误解析

Step 3: device 子命令
  cmd/device.go — list (解析 adb devices)、info (getprop 合集)

Step 4: screenshot
  cmd/observe.go — adb exec-out screencap -p → PNG
  pkg/observe/screenshot.go — 支持输出到 stdout(base64) / 文件

Step 5: ui tree
  pkg/observe/uitree.go — adb shell uiautomator dump → 解析 XML → JSON 带编号
  cmd/ui.go — tree / find 子命令

Step 6: observe (合并)
  pkg/observe/combined.go — 并行截屏+UI树，合并输出

Step 7: 输入命令
  cmd/input.go — tap (坐标/index/id/text) / swipe / key / type
  pkg/input/adbinput.go — adb shell input 封装
  tap --index 需要先调 ui tree 获取坐标

Step 8: app 子命令
  cmd/app.go — list / current / launch / stop

Step 9: skill + doctor
  cmd/skill.go — 输出 skill.json
  cmd/doctor.go — 检查 adb 可用 / 设备连接 / 基本功能

Step 10: 集成测试
  用一台真机跑通: observe → tap → observe 循环
```

### 11.3 Phase 1 不做的事

- 不实现 adbclawd（Phase 2）
- 不实现 sendevent（Phase 2）
- 不实现人类化模拟（Phase 3）
- 不实现 MCP server（Phase 3）
- 不实现无线 ADB 优先策略（Phase 3）
- 不实现 note 子命令（Phase 3）
- 不自实现 ADB 协议（Phase 2+ 考虑）

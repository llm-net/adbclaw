# adbclaw — 带反侦测能力的 Android 设备控制 CLI

> 用于 OpenClaw 平台的命令行工具，通过 ADB 控制连接在 Mac 上的 Android 手机，同时尽可能避免被手机上的 App 侦测。

---

## 1. App 侦测手段全景

Android App 有多层手段来检测设备是否被外部控制：

### 1.1 ADB 连接层检测

| 检测方法 | API / 属性 | 说明 |
|----------|-----------|------|
| USB Debugging 开关 | `Settings.Global.ADB_ENABLED` = 1 | 最常见检查，银行/支付类 App 必查 |
| 开发者选项开关 | `Settings.Global.DEVELOPMENT_SETTINGS_ENABLED` = 1 | 很多 App 同时检查此项 |
| ADB 守护进程状态 | `getprop init.svc.adbd` = "running" | 进程级检测 |
| USB 连接状态 | `getprop sys.usb.state` 包含 "adb" | USB 模式字符串检测 |
| 可调试构建 | `getprop ro.debuggable` = 1 | 只在 userdebug/eng 构建中为 1 |

### 1.2 输入事件检测（关键）

| 检测方法 | 技术细节 |
|----------|---------|
| `MotionEvent.getDeviceId() == -1` | Android 对所有注入事件设置 deviceId = -1（VIRTUAL_KEYBOARD_ID），真实触屏事件 deviceId > 0。**这是最可靠的检测手段** |
| `MotionEvent.getSource()` | 注入事件可能报告 SOURCE_UNKNOWN，真实触摸报告 SOURCE_TOUCHSCREEN |
| `FLAG_IS_ACCESSIBILITY_EVENT` (0x800) | 无障碍服务生成的事件带此标志 |
| 压力值和触摸面积 | `adb shell input tap` 压力固定为 1.0，无触摸面积数据；真实触摸有连续变化的压力和面积 |
| 时间分析 | 真实人类触摸的时间间隔呈自然方差分布，自动化输入时间间隔过于规律 |

### 1.3 行为分析层

- **速度与轨迹分析**：真实滑动有加速/减速曲线，线性插值可被检测
- **传感器交叉验证**：部分 App 检查陀螺仪/加速度计数据是否与触摸事件关联——点击时设备应有微振动
- **多点触控模式**：真实交互包含微调整、误触修正等自然噪声
- **无障碍服务枚举**：App 调用 `AccessibilityManager.getEnabledAccessibilityServiceList()` 检测自动化服务

### 1.4 商业 SDK 检测

Appdome、Promon、Guardsquare 等商业 SDK 综合数十个信号，结合 ML 分析，同时使用 Play Integrity API 进行硬件级认证。几乎无法完美绕过。

---

## 2. 输入注入方式对比

| 方法 | 免 Root | deviceId 真实 | Flags 干净 | 行为逼真度 | 隐蔽性 |
|------|---------|--------------|-----------|-----------|--------|
| `adb shell input tap/swipe` | Yes | No (-1) | No | 低 | **很低** |
| `sendevent /dev/input` | Yes | **Yes** | **Yes** | 中 | **高** |
| scrcpy UHID 模式 | Yes | **Yes** | **Yes** | 中 | **高**（仅键鼠） |
| Accessibility Service | Yes | 部分 | No (有 flag) | 中 | 中低 |
| UIAutomator | Yes | No | No | 低 | 低 |
| Frida hook 注入 | 需 Root | 可配置 | 可配置 | 高 | 高 |

### 核心结论

**`sendevent /dev/input`** 是免 Root 下隐蔽性最高的方案：
- 事件走真实硬件输入路径
- `getDeviceId()` 返回真实触屏设备 ID（不是 -1）
- `getSource()` 报告 `SOURCE_TOUCHSCREEN`
- 可以自定义压力值和触摸面积
- ADB shell 用户有 `/dev/input/eventX` 的写权限

### sendevent 示例

```bash
# 在 (500, 800) 处模拟一次触摸
sendevent /dev/input/event2 3 57 0    # ABS_MT_TRACKING_ID
sendevent /dev/input/event2 3 53 500  # ABS_MT_POSITION_X
sendevent /dev/input/event2 3 54 800  # ABS_MT_POSITION_Y
sendevent /dev/input/event2 3 48 5    # ABS_MT_TOUCH_MAJOR (触摸面积)
sendevent /dev/input/event2 3 58 50   # ABS_MT_PRESSURE (压力值)
sendevent /dev/input/event2 0 0 0     # SYN_REPORT
# ... 抬起事件
sendevent /dev/input/event2 3 57 -1   # ABS_MT_TRACKING_ID = -1 (释放)
sendevent /dev/input/event2 0 0 0     # SYN_REPORT
```

---

## 3. 架构设计

```
adbclaw CLI (Mac 端)
│
├── 连接管理层
│   ├── USB ADB 连接
│   ├── Wireless ADB (Android 11+，避免 sys.usb.state 暴露 "adb")
│   └── 设备发现与多设备管理
│
├── 屏幕采集层
│   ├── scrcpy server（推送到手机端，高效截屏/视频流）
│   ├── screencap fallback（单帧截图）
│   └── 图像输出：base64 / 文件 / stdout pipe
│
├── 输入注入层（分级策略）
│   ├── Level 1 (高隐蔽): sendevent /dev/input
│   │   ├── 自动探测触屏设备节点 (getevent -pl 解析)
│   │   ├── 适配 Type A / Type B multitouch 协议
│   │   ├── 真实 pressure / touch_major 值（高斯分布）
│   │   └── 人类化时间抖动（gaussian jitter + 贝塞尔曲线轨迹）
│   ├── Level 2 (中隐蔽): scrcpy UHID (键盘/鼠标)
│   └── Level 3 (兜底): adb shell input (简单场景)
│
├── 反侦测层
│   ├── stealth 模式 (免 Root)
│   │   ├── 无线 ADB 连接
│   │   ├── sendevent 输入注入
│   │   └── 关闭开发者选项（ADB 仍保持工作）
│   └── deep-stealth 模式 (需 Root)
│       ├── Magisk + LSPosed + DevOptsHide
│       ├── strongR-frida（反检测版 Frida）
│       └── Settings.Global hook（ADB_ENABLED 返回 0）
│
├── 设备信息层
│   ├── 屏幕尺寸/分辨率/DPI 查询
│   ├── 已安装 App 列表
│   ├── 当前 Activity / UI 层级
│   └── 设备状态（电量、网络、方向）
│
└── 对外接口层
    ├── CLI Commands（人类直接使用）
    ├── MCP Server（供 Claude / AI Agent 调用）
    ├── JSON stdout 模式（供脚本 pipe）
    └── gRPC / REST API（可选，供平台集成）
```

---

## 4. CLI 命令设计（草案）

```bash
# 设备管理
adbclaw devices                          # 列出连接的设备
adbclaw connect <ip>:<port>              # 无线连接
adbclaw use <device-id>                  # 选择默认设备

# 屏幕
adbclaw screenshot                       # 截屏到 stdout (PNG)
adbclaw screenshot -o ./shot.png         # 截屏到文件
adbclaw stream                           # 实时视频流 (供 AI 视觉)
adbclaw screen-info                      # 屏幕尺寸/分辨率/DPI

# 输入（默认使用 sendevent 高隐蔽模式）
adbclaw tap 500 800                      # 点击坐标
adbclaw tap 500 800 --humanize           # 人类化点击（抖动+压力变化）
adbclaw swipe 100 500 400 500 --duration 300ms
adbclaw type "hello world"               # 文字输入
adbclaw key HOME | BACK | ENTER          # 按键
adbclaw gesture <gesture-file.json>      # 复杂手势（从文件加载）

# 输入模式控制
adbclaw tap 500 800 --method sendevent   # 强制 sendevent
adbclaw tap 500 800 --method input       # 强制 adb shell input
adbclaw tap 500 800 --method uhid        # 强制 scrcpy UHID

# App 管理
adbclaw app list                         # 已安装 App
adbclaw app current                      # 当前前台 App/Activity
adbclaw app launch com.example.app       # 启动 App
adbclaw app stop com.example.app         # 停止 App

# UI 信息
adbclaw ui dump                          # 当前 UI 层级 (XML)
adbclaw ui tree                          # 可访问性树 (结构化)
adbclaw ui find --text "Login"           # 按文字查找元素
adbclaw ui find --id "btn_submit"        # 按 ID 查找元素

# 反侦测
adbclaw stealth status                   # 检查当前反侦测状态
adbclaw stealth enable                   # 启用 stealth 模式
adbclaw stealth deep                     # 启用 deep-stealth (需 Root)

# 设备状态
adbclaw info                             # 设备综合信息
adbclaw shell <command>                  # 执行 shell 命令

# AI Agent 接口
adbclaw mcp serve                        # 启动 MCP Server
adbclaw api serve --port 8400            # 启动 REST API
```

---

## 5. 关键技术要点

### 5.1 sendevent 协议适配

不同手机的触屏使用不同的 multitouch 协议：

- **Type B**（主流，Android 4.0+）：每个触摸点有 `ABS_MT_SLOT` 和 `ABS_MT_TRACKING_ID`
- **Type A**（旧设备）：不使用 slot，通过 `SYN_MT_REPORT` 分隔触摸点

adbclaw 需要在首次连接时通过 `getevent -pl` 探测：
- 触屏设备节点路径（`/dev/input/eventX`）
- 支持的事件类型和坐标范围（`ABS_MT_POSITION_X min/max`）
- 协议类型（Type A vs Type B）

### 5.2 人类化输入模拟

```
真实人类触摸特征：
├── 坐标：目标点 ± 2-5px 高斯偏移
├── 压力：40-80 范围内的连续变化曲线
├── 触摸面积：3-8 范围内随压力变化
├── 按下时间：80-200ms (正态分布, μ=120ms, σ=30ms)
├── 点击间隔：200-800ms (对数正态分布)
├── 滑动轨迹：三阶贝塞尔曲线 (非线性插值)
└── 滑动速度：先加速后减速 (ease-in-out)
```

### 5.3 无线 ADB 策略

```bash
# 初始通过 USB 启用无线 ADB
adb tcpip 5555
adb connect <phone-ip>:5555
# 断开 USB，后续全程无线
# 此时 sys.usb.state 不含 "adb"，部分检测失效
```

Android 11+ 原生支持无线调试（Settings > Developer Options > Wireless Debugging），无需先 USB 连接。

### 5.4 设备端辅助程序

推送一个轻量 binary 到手机 `/data/local/tmp/`（类似 scrcpy 的做法），用于：
- 高速执行 sendevent 序列（避免每次 `adb shell sendevent` 的进程启动开销）
- 批量写入 `/dev/input/eventX`（直接 `write()` 比逐条 sendevent 快 100x）
- 实时截屏并通过 socket 传回 Mac 端

### 5.5 反侦测分层策略

#### stealth 模式（免 Root，覆盖大多数普通 App）

| 措施 | 绕过的检测 |
|------|-----------|
| 无线 ADB | `sys.usb.state` 不含 "adb" |
| 关闭开发者选项 | `DEVELOPMENT_SETTINGS_ENABLED` = 0 |
| sendevent 输入 | `deviceId != -1`、无注入 flag |
| 人类化模拟 | 时间/轨迹行为分析 |

**无法绕过**：`ADB_ENABLED` 仍为 1、`init.svc.adbd` 仍为 running

#### deep-stealth 模式（需 Root，覆盖银行级 App）

在 stealth 基础上增加：

| 措施 | 绕过的检测 |
|------|-----------|
| Magisk + DenyList | Root 检测、Play Integrity |
| LSPosed + DevOptsHide | `ADB_ENABLED` 返回 0 |
| strongR-frida | 运行时 hook，隐藏 Frida 自身 |
| SettingsFirewall | 按 App 粒度屏蔽 Settings 读取 |

---

## 6. 竞品与可复用项目

### 6.1 直接竞品

| 项目 | Stars | 说明 | 与 adbclaw 的差异 |
|------|-------|------|------------------|
| [DroidRun](https://github.com/droidrun/droidrun) | 7.7k | LLM 驱动 Android 控制 CLI，多模型支持 | **无反侦测能力**，Python 实现 |
| [mobile-use](https://github.com/minitap-ai/mobile-use) | 2.2k | AI agent 控制真机，AndroidWorld 基准第一 | **无反侦测**，面向测试场景 |
| [DroidClaw](https://github.com/unitedbyai/droidclaw) | 875 | 自然语言→ADB 操作 | **无反侦测**，TypeScript，名字已占用 |
| [agent-device](https://github.com/callstackincubator/agent-device) | 720 | Callstack 的 AI agent 设备控制 CLI | **无反侦测**，TypeScript |

### 6.2 可复用组件

| 项目 | Stars | 可复用点 |
|------|-------|---------|
| [scrcpy](https://github.com/Genymobile/scrcpy) | 136k | server jar（截屏/视频流）、UHID 模式参考 |
| [minitouch](https://github.com/openstf/minitouch) | 669 | 低级 socket 触摸注入协议设计参考 |
| [CyAndroEmu](https://github.com/hansalemaos/cyandroemu) | 127 | 无 ADB 方案参考、人类化鼠标模拟 |
| [strongR-frida](https://github.com/hzzheyang/strongR-frida-android) | 1.6k | deep-stealth 模式的 Frida 组件 |
| [DeviceSpoofLab](https://github.com/yubunus/DeviceSpoofLab-Hooks) | 新 | 设备指纹伪装（126+ 属性） |
| [py-scrcpy-client](https://github.com/leng-yue/py-scrcpy-client) | 419 | scrcpy 编程接口参考 |

### 6.3 差异化定位

**现有项目都没有做反侦测** — 这是 adbclaw 的核心差异点。

adbclaw = AI Agent 控制能力 + 反侦测能力 + CLI 工具链

---

## 7. 可行性评估

| 维度 | 评估 | 说明 |
|------|------|------|
| 基本设备控制 | **完全可行** | ADB 技术成熟，大量参考实现 |
| sendevent 高隐蔽输入 | **可行** | 免 Root，shell 用户有 /dev/input 写权限 |
| 免 Root 反侦测 | **部分可行** | 能绕过输入检测 + 部分 Settings 检测，但 `ADB_ENABLED` 无法隐藏 |
| 有 Root 反侦测 | **高度可行** | Magisk + LSPosed 组合拳成熟，但需同时隐藏 Root |
| 骗过商业安全 SDK | **极难** | ML + 硬件认证，接近不可能完美绕过 |
| MCP / AI Agent 集成 | **完全可行** | MCP 协议成熟，已有多个参考实现 |
| Mac 端 CLI 开发 | **完全可行** | Go 或 Rust 均可，跨平台编译方便 |

---

## 8. 技术选型建议

| 组件 | 推荐 | 原因 |
|------|------|------|
| CLI 主体语言 | **Go** | 与项目技术栈一致，交叉编译方便，单 binary 分发 |
| 设备端辅助程序 | **C**（交叉编译到 ARM） | 直接操作 /dev/input，性能关键 |
| ADB 协议 | 复用 scrcpy 的 server jar + 自实现 ADB 协议 | 避免依赖 adb 二进制 |
| 截屏/视频 | scrcpy server | 已有成熟方案，支持 H.265 编码 |
| MCP Server | Go 实现 | 与 CLI 同进程，简化部署 |

---

## 9. 命名说明

- **adbclaw**: 本项目名称，ADB + Claw（OpenClaw 体系）
- **OpenClaw**: 注意 github.com/openclaw/openclaw 已被占用（10万+ stars 的 AI 助手平台），如果 OpenClaw 是你自己的项目体系名称则无冲突，但 GitHub org 名需要另选

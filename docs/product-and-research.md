# adb-claw — 产品目标与技术调研

## 1. 产品定位

adb-claw 是一个**纯 CLI 工具层**的 Android 设备控制工具，类似 `kubectl` 之于 Kubernetes、`gh` 之于 GitHub：

- 提供完整的 Android 设备操作 CLI，所有命令输出结构化 JSON
- 通过 `adb` 命令与 Android 设备通信，无需在手机上安装任何 App 或服务
- 附带 Skill 描述文件，让 AI agent 知道自己能调什么
- Go 单 binary 分发，仅依赖系统已安装的 `adb`
- **不包含** LLM 调用、prompt、agent loop、任务规划

### 发布渠道

- **Claude Code 插件**：通过 Plugin Marketplace 安装，赋予 Claude Code Android 控制能力
- **OpenClaw Skill**：发布在 ClawHub 上，为 OpenClaw agent 提供 Android 控制能力

两个渠道共享同一套 CLI 工具和 App Profile。

### 与竞品的核心差异

| | DroidRun | mobile-use | adb-claw |
|---|---|---|---|
| 定位 | LLM Agent + 设备控制一体 | LLM Agent + 设备控制一体 | **纯设备控制工具** |
| 语言 | Python | Python | Go |
| LLM 依赖 | llama-index (必须) | LangGraph + LangChain (必须) | **无** |
| 手机端依赖 | Portal APK + AccessibilityService | UIAutomator2 server | **无**（纯 adb 命令） |
| Agent 集成 | 内嵌 agent | 内嵌 agent | **Skill 文件 + JSON stdout** |
| 分发 | pip install | pip install | **单 binary** |

**核心差异点**：纯工具层定位 + 零手机端依赖 + AI agent 友好的结构化输出。

---

## 2. 技术调研：App 侦测手段

> 以下为调研记录，供未来需要反侦测能力时参考。当前版本不实现反侦测。

### 2.1 ADB 连接层检测

| 检测方法 | API / 属性 | 说明 |
|----------|-----------|------|
| USB Debugging 开关 | `Settings.Global.ADB_ENABLED` = 1 | 最常见检查，银行/支付类 App 必查 |
| 开发者选项开关 | `Settings.Global.DEVELOPMENT_SETTINGS_ENABLED` = 1 | 很多 App 同时检查 |
| ADB 守护进程状态 | `getprop init.svc.adbd` = "running" | 进程级检测 |
| USB 连接状态 | `getprop sys.usb.state` 包含 "adb" | USB 模式字符串检测 |

### 2.2 输入事件检测

| 检测方法 | 技术细节 |
|----------|---------|
| `MotionEvent.getDeviceId() == -1` | `adb shell input` 注入事件 deviceId = -1（VIRTUAL_KEYBOARD_ID），真实触屏 deviceId > 0。**最可靠的检测手段** |
| `MotionEvent.getSource()` | 注入事件可能报告 SOURCE_UNKNOWN，真实触摸报告 SOURCE_TOUCHSCREEN |
| 压力值和触摸面积 | `adb shell input tap` 压力固定 1.0，无触摸面积；真实触摸有连续变化 |
| 时间分析 | 真实人类触摸时间间隔呈自然方差分布，自动化输入过于规律 |

### 2.3 行为分析层

- **速度与轨迹分析**：真实滑动有加速/减速曲线，线性插值可被检测
- **传感器交叉验证**：部分 App 检查陀螺仪/加速度计数据是否与触摸事件关联
- **无障碍服务枚举**：`AccessibilityManager.getEnabledAccessibilityServiceList()` 检测自动化服务

### 2.4 商业 SDK 检测

Appdome、Promon、Guardsquare 等商业 SDK 综合数十个信号 + ML 分析 + Play Integrity API 硬件级认证。几乎无法完美绕过。

---

## 3. 技术调研：输入注入方式

| 方法 | 免 Root | deviceId 真实 | 隐蔽性 | 说明 |
|------|---------|--------------|--------|------|
| `adb shell input tap/swipe` | Yes | No (-1) | 低 | **当前使用**，简单可靠 |
| `sendevent /dev/input` | Yes | **Yes** | **高** | 需设备端程序批量写入才实用 |
| scrcpy UHID 模式 | Yes | **Yes** | **高** | 仅键鼠，不支持触摸 |
| Frida hook 注入 | 需 Root | 可配置 | 高 | 复杂度高 |

当前使用 `adb shell input`，简单可靠，覆盖绝大多数自动化场景。sendevent 方案需要设备端常驻程序配合，作为未来研究方向保留。

---

## 4. 架构

```
┌─────────────────────────────────────────────────────────┐
│  Claude Code / OpenClaw Agent / 其他 Bot                  │
│  (读 Skill 描述知道 adb-claw 能做什么，                      │
│   调用 CLI 命令，解析 JSON stdout)                         │
└────────────────────┬────────────────────────────────────┘
                     │ 子进程调用
┌────────────────────▼────────────────────────────────────┐
│  adb-claw CLI (Go binary, Mac/Linux)                      │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐ │
│  │ device   │ │ input    │ │ observe  │ │ app        │ │
│  │ 设备管理  │ │ 输入操作  │ │ 屏幕/UI  │ │ App 管理   │ │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └─────┬──────┘ │
│  ┌────▼────────────▼────────────▼──────────────▼──────┐ │
│  │         Commander 接口 (pkg/adb)                    │ │
│  │  Shell() / ExecOut() / RawCommand()                 │ │
│  └────────────────────┬───────────────────────────────┘ │
└───────────────────────┼─────────────────────────────────┘
                        │ exec adb shell / adb exec-out
┌───────────────────────▼─────────────────────────────────┐
│  Android 设备 (通过 USB / WiFi ADB 连接)                  │
│  ├── adb shell input tap/swipe/keyevent/text            │
│  ├── adb exec-out screencap -p                          │
│  ├── adb shell uiautomator dump                         │
│  ├── adb shell pm / am / dumpsys                        │
│  └── adb shell getprop                                  │
└─────────────────────────────────────────────────────────┘
```

所有操作通过标准 `adb` 命令完成，无需在设备上安装或运行任何额外程序。

---

## 5. 竞品分析与借鉴

### 5.1 从 DroidRun 借鉴

- **采纳**：`DeviceDriver` 抽象接口 → adb-claw 的 Commander 接口；UI 元素编号 → `ui tree` 带 index；`tap --index` 按编号点击
- **不采纳**：Portal APK 依赖、llama-index 编排

### 5.2 从 mobile-use 借鉴

- **采纳**：Target 三级回退（bounds → resource_id → text）；截图 + UI 树并行获取 → `observe` 命令；截图缩放
- **不采纳**：UIAutomator2 server、LangGraph 状态图

---

## 6. 未来研究方向

以下为可能的演进方向，不在近期开发计划内：

### 6.1 MCP Server 集成

将 adb-claw 的所有命令暴露为 MCP tools，让 Claude Desktop / OpenClaw 通过标准 MCP 协议调用，替代子进程方式。

### 6.2 Unicode 文字输入

`adb shell input text` 不支持 Unicode（中文等）。可能的方案：通过 `adb shell` 操作剪贴板 + 模拟粘贴，或利用 App 的深度链接绕过输入。

### 6.3 sendevent 高隐蔽输入

如果未来有反侦测需求，可研究 sendevent 方案：需要设备端辅助程序批量写入 `/dev/input/eventX`，deviceId 真实，绕过输入事件检测。需配合触屏参数探测（`getevent -pl`）和 Type A/B 协议适配。

### 6.4 人类化输入模拟

在 sendevent 基础上叠加人类行为模拟：高斯坐标抖动、贝塞尔曲线轨迹、压力/面积连续变化、时间间隔自然方差。

---

## 7. 已知限制

| 限制 | 说明 |
|------|------|
| `uiautomator dump` 速度慢（~1-2s） | Android 系统限制，暂无更好方案 |
| `uiautomator dump` 在动画中可能失败 | 先暂停动画（如点击暂停视频），再 dump |
| 文本输入仅支持 ASCII | `adb shell input text` 的限制，CJK 字符需用深度链接绕过 |
| 输入事件可被 App 检测（deviceId=-1） | `adb shell input` 的固有特征，当前不做反侦测 |
| 截屏为 PNG 格式（~5MB） | 通过 `--width` 缩放减少体积，未来可考虑 JPEG 转换 |

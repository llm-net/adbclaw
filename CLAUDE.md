# adbclaw

Android 设备控制 CLI，供 AI agent 自动化调用。纯工具层，不含 LLM/Agent 逻辑。

## 两种发布渠道

adbclaw 同时作为两个平台的 Skill 发布，共享同一套 CLI 工具和 App Profile：

### 1. Claude Code 插件（Plugin Marketplace）

通过 Claude Code 插件市场安装，让 Claude Code 获得 Android 设备控制能力。

```
.claude-plugin/              # 插件配置
├── plugin.json              # 插件元数据（名称、版本、作者）
└── marketplace.json         # 市场发布配置
skills/android-control/      # Claude Code Skill 定义
└── SKILL.md                 # Skill 描述（触发条件、命令文档、工作流模式）
```

- **触发方式**：用户在 Claude Code 中提到 Android 控制/自动化时自动激活
- **二进制分发**：SessionStart hook 自动下载 adbclaw，或手动运行 `scripts/setup.sh`
- **Skill 路径**：`${CLAUDE_PLUGIN_ROOT}/bin/adbclaw`

### 2. OpenClaw Skill（ClawHub）

作为 [OpenClaw](https://github.com/openclaw/openclaw) 平台的 Skill 发布在 ClawHub 上，为 OpenClaw agent 提供 Android 设备控制能力。

```
skills/adb-claw/             # OpenClaw Skill 定义
└── SKILL.md                 # YAML frontmatter（元数据 + 安装依赖）+ Skill 文档
```

- **Skill 格式**：`SKILL.md` 的 YAML frontmatter 包含 `metadata.openclaw`，定义 OS 要求、依赖二进制、安装脚本
- **安装方式**：OpenClaw 自动检测 `requires.bins`，缺失时按 `install` 配置自动安装
- **依赖**：`adbclaw`（curl 安装）+ `adb`（brew 安装）

### 两者的关系

| | Claude Code 插件 | OpenClaw Skill |
|---|---|---|
| Skill 定义 | `skills/android-control/SKILL.md` | `skills/adb-claw/SKILL.md` |
| 配置文件 | `.claude-plugin/` | YAML frontmatter |
| 安装方式 | setup.sh / hook 自动下载 | OpenClaw 自动安装 |
| 内容差异 | 两者 Skill 文档内容基本一致，格式适配各自平台规范 |

### 共享资源

```
skills/apps/                 # App Profile（两个渠道共享）
├── README.md                # Profile 编写规范
└── douyin.md                # 抖音操作知识（深度链接、布局、已知问题）
```

App Profile 是针对具体 App 的操作知识库，包含深度链接、已知布局、设备差异和常见问题。两个平台的 Skill 都引用同一套 Profile，新增 Profile 后两边同时受益。

## 项目结构

```
src/                    # Go 代码根目录（go.mod 在此）
├── main.go             # 入口
├── Makefile            # 构建脚本
├── cmd/                # Cobra CLI 命令
│   ├── root.go         # 根命令 + 全局 flags（-s, -o, --timeout, --verbose）
│   ├── device.go       # device list / info
│   ├── observe.go      # observe / screenshot
│   ├── ui.go           # ui tree / find
│   ├── input.go        # tap / long-press / swipe / key / type
│   ├── app.go          # app list / current / launch / stop
│   ├── doctor.go       # 环境检查
│   ├── skill.go        # 输出 skill.json (go:embed)
│   └── skill.json      # 嵌入的 AI agent 能力描述
└── pkg/
    ├── adb/shell.go        # Commander 接口 + Client 实现
    ├── input/adbinput.go   # Tap/Swipe/LongPress/Key/Type
    ├── output/envelope.go  # JSON 响应 envelope + Writer
    └── observe/
        ├── screenshot.go   # 截屏 + 缩放
        ├── uitree.go       # XML 解析 + UI 树索引
        └── combined.go     # 并行 observe（截屏 + UI 树）

docs/                   # 技术文档
├── product-and-research.md    # 产品目标 + 技术调研 + 未来方向
└── development-roadmap.md     # 开发计划
website/                # React + Vite + Tailwind 官网
```

## 构建

```bash
cd src
make build   # 产物 → bin/adbclaw（项目根目录）
make test    # go test ./...
make lint    # go vet
make clean
```

Go 1.24，依赖 cobra v1.10.2 + golang.org/x/image v0.36.0。构建产物在项目根目录 `bin/`（已 gitignore）。

## 本地开发加载

开发完成后，编译并加载到 Claude Code：

```bash
cd src && make build   # 编译到 bin/adbclaw
claude --plugin-dir .  # 在项目根目录启动，加载当前目录为插件
```

- `make build` 产物输出到项目根目录 `bin/adbclaw`，与插件 SKILL.md 和 `setup.sh` 引用的路径一致
- SessionStart hook 检测到 `bin/adbclaw` 已存在会跳过下载，直接使用本地编译版本
- 已有会话中修改代码后，重新 `make build` + 重启 Claude Code 即可生效

## 架构要点

- **Commander 接口** (`pkg/adb/shell.go`) — 所有 pkg 通过 `Commander` 接口调用 ADB，测试用 mock。包含 `Shell()`、`ExecOut()`（二进制安全）、`RawCommand()` 三个方法
- **JSON Envelope** (`pkg/output/envelope.go`) — 统一 `{ok, command, data, error, duration_ms, timestamp}`。error 含 `{code, message, suggestion}`。支持 json/text/quiet 三种输出模式
- **UI 树过滤** (`pkg/observe/uitree.go`) — 只索引有 text/resource-id/content-desc 或 clickable/scrollable 的节点，减少 agent 噪音。Element 带 index/bounds/center
- **输入为顶级命令** — `adbclaw tap` 而非 `adbclaw input tap`
- **observe 部分失败容忍** — 截屏和 UI 树并行（sync.WaitGroup），互不阻塞
- **输入命令支持元素定位** — tap/long-press 支持 `--index`/`--id`/`--text` 直接定位 UI 元素
- **文本输入安全** — `type` 命令转义 shell 特殊字符，拒绝非 ASCII 字符

## 命令树

```
adbclaw
├── device list                    # 列出已连接设备
├── device info                    # 设备详情（型号/版本/屏幕尺寸/密度）
├── observe [--width px]           # 截屏 + UI 树并行
├── screenshot [--file path] [--width px]  # 截屏（base64 或文件）
├── ui tree                        # UI 元素树（带 index）
├── ui find --text/--id/--index    # 查找 UI 元素
├── tap <x> <y> | --index/--id/--text     # 点击
├── long-press <x> <y> [--duration ms]    # 长按
├── swipe <x1> <y1> <x2> <y2> [--duration ms]  # 滑动
├── key <HOME|BACK|ENTER|...>      # 按键（20+ 别名）
├── type <text>                    # 输入文本（仅 ASCII）
├── app list [--all]               # 应用列表（默认三方）
├── app current                    # 当前前台应用
├── app launch <package>           # 启动应用
├── app stop <package>             # 停止应用
├── skill                          # 输出 skill.json
└── doctor                         # 环境检查（adb/设备/能力）
```

## 全局 Flags

```
-s, --serial <id>      # 目标设备（多设备时）
-o, --output <format>  # json（默认）| text | quiet
--timeout <ms>         # 命令超时（默认 30000）
--verbose              # 调试输出到 stderr
```

## 代码约定

- 新命令放 `src/cmd/`，新包放 `src/pkg/`
- 所有 ADB 调用必须通过 `Commander` 接口，不直接 exec
- 命令输出必须使用 `output.Writer` 写 JSON envelope
- 测试文件与源码同目录，用 `_test.go` 后缀
- 错误码用大写下划线格式，如 `ELEMENT_NOT_FOUND`、`DEVICE_NOT_FOUND`
- 两个 Skill 的 `SKILL.md` 内容保持同步，修改一个时需同步更新另一个

## 技术方案

所有操作通过标准 `adb` 命令完成（`adb shell input`、`adb exec-out screencap`、`adb shell uiautomator dump` 等），无需在设备上安装或运行任何额外程序。产品目标与技术调研见 `docs/product-and-research.md`，开发计划见 `docs/development-roadmap.md`。

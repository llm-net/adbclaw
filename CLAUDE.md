# adbclaw

Android 设备控制 CLI，供 AI agent 自动化调用。纯工具层，不含 LLM/Agent 逻辑。

## 发布渠道

adbclaw 同时作为两个平台的 Skill 发布，**共用一份 `skills/adb-claw/SKILL.md`**：

- **Claude Code**：通过插件市场安装（`.claude-plugin/`），按 `## Triggers` 触发，`## Binary` 指示二进制位置
- **OpenClaw**：通过 ClawHub 安装，读取 YAML frontmatter 中的 `metadata.openclaw`（OS 要求、依赖、安装脚本）

两个平台读同一个文件，Claude Code 忽略 frontmatter，OpenClaw 忽略 `## Binary` 段落。

```
.claude-plugin/              # Claude Code 插件配置
├── plugin.json              # 插件元数据
└── marketplace.json         # 市场发布配置
skills/
├── adb-claw/SKILL.md        # Skill 定义（两个平台共用）
└── apps/                    # App Profile 知识库（运行时按需加载）
    ├── README.md            # Profile 编写规范
    └── douyin.md            # 抖音（深度链接、布局、已知问题）
```

### App Profile

App Profile 是针对具体 App 的操作知识库，Agent 运行时按需读取：

1. `adbclaw app current` → 获取当前 App 包名
2. 检查 `skills/apps/` 下有无对应 Profile
3. 有 → 按 Profile 操作（深度链接、已知布局）；无 → 常规 observe 探索

新增 App 支持只需往 `skills/apps/` 丢一个 `.md` 文件。

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
- `skills/adb-claw/SKILL.md` 同时服务 Claude Code 和 OpenClaw，修改时需兼顾两个平台的格式要求

## 技术方案

所有操作通过标准 `adb` 命令完成（`adb shell input`、`adb exec-out screencap`、`adb shell uiautomator dump` 等），无需在设备上安装或运行任何额外程序。产品目标与技术调研见 `docs/product-and-research.md`，开发计划见 `docs/development-roadmap.md`。

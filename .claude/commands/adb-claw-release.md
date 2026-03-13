---
description: "正式发布指定版本到 GitHub Releases + ClawHub + 官网"
allowed-tools: "Bash, Read, Edit, Grep, Glob, Write, Agent"
---

# 正式发布 $ARGUMENTS

按以下步骤**完整执行**发布流程。发布涉及三个平台：

- **GitHub Releases** — CI 自动构建 4 平台二进制（darwin/linux × arm64/amd64）
- **ClawHub** — OpenClaw 技能市场（https://clawhub.ai/dionren/adb-claw）
- **GitHub Pages** — 官网 adb-claw.llm.net

---

## Step 1: 同步版本号

将以下 4 处版本号更新为 `$ARGUMENTS`：

```
.claude-plugin/plugin.json      → "version": "$ARGUMENTS"
skills/adb-claw/SKILL.md        → version: $ARGUMENTS (frontmatter 第 3 行)
                                 → metadata.openclaw.version: "$ARGUMENTS" (第 11 行)
src/cmd/root.go                  → 注释中的 ldflags 示例版本
```

## Step 2: 更新文档

检查以下文件是否需要更新（如有未提交的功能变更）：

- `skills/adb-claw/SKILL.md` — 命令文档、Getting Started、Troubleshooting
- `README.md` — 与 SKILL.md 对齐（Features、命令树、Usage、Architecture、App Profiles 表）
- `CLAUDE.md` — 命令树、项目结构
- `skills/apps/*.md` — App Profile 变更

## Step 3: ClawHub 安全审查

**发布前必须审查 `skills/adb-claw/SKILL.md`**，确保不会触发 ClawHub 安全扫描告警。

逐项检查以下规则：

### 3.1 Install Mechanism（最常触发告警）

- **禁止** `kind: "script"` + `curl | bash` 模式 — 这会被标记为 Suspicious
- **必须**使用 `kind: "download"` 直接下载二进制，或 `kind: "brew"` 包管理器安装
- 所有 download URL 必须指向 `github.com/llm-net/adb-claw`（不是 AdbClaw 或其他旧 URL）
- 确认 frontmatter `homepage` 字段为 `https://github.com/llm-net/adb-claw`

```
# 合规示例
{ "kind": "download", "url": "https://github.com/llm-net/adb-claw/releases/latest/download/adb-claw-darwin-arm64" }
{ "kind": "brew", "formula": "android-platform-tools" }

# 违规示例（会触发告警）
{ "kind": "script", "script": "curl -fsSL ... | bash" }
```

### 3.2 Purpose & Capability

- `name` / `description` 必须与运行时功能一致
- `requires.bins` 只列必需的二进制（adb-claw, adb）

### 3.3 Instruction Scope

- SKILL.md 正文只包含 adb-claw/adb 命令指引
- 不能有读取无关本地文件、网络请求、数据外传的指令

### 3.4 Credentials

- 不请求任何环境变量、密钥、配置文件路径

### 3.5 Persistence & Privilege

- 不设置 `always: true`
- 不请求系统级配置修改

用 Grep 检查：
```bash
# 必须无结果
grep -n "curl.*bash\|kind.*script\|always.*true\|credential\|secret\|api.key\|AdbClaw" skills/adb-claw/SKILL.md
```

如发现问题，修复后再继续。

## Step 4: 运行测试 & 构建

```bash
export PATH="/Users/dionren/go-sdk/go/bin:$PATH"
cd src && make test && make build
```

测试全部通过、构建成功后才能继续。

## Step 5: 提交 & 推送

```bash
git add <所有变更文件>
git commit -m "feat: v$ARGUMENTS — 简要描述"
git push origin main
```

提交信息使用 `feat: vX.Y.Z — 简要描述` 格式。

## Step 6: 打 tag 触发 GitHub Release

```bash
git tag v$ARGUMENTS
git push origin v$ARGUMENTS
```

推送 tag 后 GitHub Actions 自动执行：

- **Release** workflow：test → 交叉编译 4 平台 → 创建 GitHub Release（含 6 个 assets）
- **Deploy Website** workflow：构建 website/ → 部署到 GitHub Pages（adb-claw.llm.net）

## Step 7: 发布到 ClawHub

**GitHub Release 只覆盖二进制分发，ClawHub 必须单独发布：**

```bash
clawhub publish skills/adb-claw --version $ARGUMENTS --changelog "变更摘要"
```

- ClawHub 上的技能路径：https://clawhub.ai/dionren/adb-claw
- 登录账号：`dionren`（`clawhub whoami` 验证）
- 发布后有安全扫描，通常几分钟后上线
- **同版本号不可重复发布**，如需重发必须 bump 版本

## Step 8: 验证

```bash
# GitHub CI 进度
gh run list --repo llm-net/adb-claw --limit 2

# GitHub Release assets（应有 6 个文件）
gh release view v$ARGUMENTS --repo llm-net/adb-claw

# ClawHub 版本确认（安全扫描中会暂时 hidden）
clawhub inspect adb-claw

# 官网部署确认
gh run list --repo llm-net/adb-claw --workflow deploy-website.yml --limit 1
```

向用户汇报三个平台的发布状态。

## 注意事项

- Git remote：`origin → llm-net/adb-claw`（主仓库，CI 和 Release 在此）
- 每步执行前确认上一步成功，不要跳步
- 如果 `clawhub publish` 报 "Version already exists"，说明该版本已发布过，需要 bump 版本号
- 如果 ClawHub 安全扫描标记为 Suspicious，检查 Step 3 的规则并修复后 bump 版本重发

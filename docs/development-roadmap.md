# adbclaw — 开发计划

## 已完成

核心功能已实现，跑通 `observe → decide → act` 循环：

- Go CLI 骨架（cobra + JSON envelope + 全局选项）
- `device list` / `device info`
- `screenshot`（PNG，支持 `--width` 缩放）
- `ui tree` / `ui find`（XML 解析 + 元素编号 + 过滤）
- `observe`（截屏 + UI 树并行）
- `tap` / `long-press` / `swipe` / `key` / `type`（支持 `--index`/`--id`/`--text` 元素定位）
- `app list` / `current` / `launch` / `stop`
- `skill` / `doctor`
- Claude Code 插件发布配置（`.claude-plugin/` + `skills/android-control/`）
- OpenClaw Skill 定义（`skills/adb-claw/`）
- App Profile 机制 + 抖音 Profile

---

## 近期计划

### 更多 App Profile

扩充 `skills/apps/` 目录，覆盖常用 App：

| App | 文件 | 关键内容 |
|-----|------|---------|
| 微信 | `wechat.md` | 发消息、小程序、朋友圈的深度链接和布局 |
| 淘宝 | `taobao.md` | 搜索商品、下单流程 |
| 小红书 | `xiaohongshu.md` | 搜索、浏览笔记 |
| B站 | `bilibili.md` | 搜索、播放视频 |

每个 Profile 需基于真机验证，标注 Phone/Pad 差异。

### Skill 完善

- 同步两个 SKILL.md 的内容，确保命令文档、工作流模式、App Profile 引用一致
- 补充更多 Troubleshooting 条目（基于实际使用中遇到的问题）
- 优化 Skill 触发条件描述

### CLI 功能增强

| 任务 | 说明 |
|------|------|
| 截屏体积优化 | 支持 `--format jpeg` 输出（Mac 端 PNG→JPEG 转换），减少 base64 体积 |
| observe 重试机制 | `uiautomator dump` 失败时自动重试，返回上一次成功的 UI 树 |
| 更多 key 别名 | 补充常用按键别名（APP_SWITCH、MENU 等） |
| `device connect` | 无线 ADB 连接（`adb tcpip` + `adb connect`） |
| `device shell` | 直接执行 shell 命令并返回 JSON envelope |

### 测试覆盖

- 补充 `pkg/observe/uitree_test.go` 的边界情况（空 XML、超大 UI 树）
- 补充 `pkg/input/adbinput_test.go` 的特殊字符转义测试
- 集成测试：真机上跑通完整交互循环

### 发布与分发

- 设置 GitHub Actions CI（build + test + lint）
- 配置 GoReleaser 自动构建多平台 binary
- 发布到 Claude Code Plugin Marketplace
- 发布到 ClawHub

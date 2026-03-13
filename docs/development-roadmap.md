# adb-claw — 开发计划

## 已完成

### 第一轮迭代 — 核心功能

跑通 `observe → decide → act` 循环：

- Go CLI 骨架（cobra + JSON envelope + 全局选项）
- `device list` / `device info`
- `screenshot`（PNG，支持 `--width` 缩放）
- `ui tree` / `ui find`（XML 解析 + 元素编号 + 过滤）
- `observe`（截屏 + UI 树并行）
- `tap` / `long-press` / `swipe` / `key` / `type`（支持 `--index`/`--id`/`--text` 元素定位）
- `app list` / `current` / `launch` / `stop`
- `skill` / `doctor`
- Claude Code 插件发布配置（`.claude-plugin/` + `skills/adb-claw/`）
- OpenClaw Skill 定义（YAML frontmatter + `metadata.openclaw`）
- App Profile 机制 + 抖音 Profile

### 第二轮迭代 — 高级命令 (v0.2.0)


| Phase | 内容 | 状态 |
|-------|------|------|
| 1 | `clear-field` + key 别名扩展（PASTE/COPY/CUT/WAKEUP/SLEEP 等） | ✅ |
| 2 | `open`（深度链接/URL）+ `scroll`（智能滚动，支持方向/页数/元素内滚动） | ✅ |
| 3 | `wait`（等待 UI 元素/Activity 出现或消失）+ `screen`（状态/亮灭/解锁/旋转） | ✅ |
| 4 | `app install` / `uninstall` / `clear` | ✅ |
| 5 | `shell`（原始 shell 命令）+ `file push` / `pull` | ✅ |

同步更新：
- SKILL.md 命令文档完善（两平台共用）
- 抖音 App Profile 基于 Phone + Pad 真机验证更新
- plugin.json / marketplace.json 版本和描述更新

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

### CLI 功能增强

| 任务 | 说明 |
|------|------|
| 截屏体积优化 | 支持 `--format jpeg` 输出，减少 base64 体积 |
| observe 重试机制 | `uiautomator dump` 失败时自动重试 |
| `device connect` | 无线 ADB 连接（`adb tcpip` + `adb connect`） |

### 测试覆盖

- 补充 `pkg/observe/uitree_test.go` 的边界情况（空 XML、超大 UI 树）
- 补充 `pkg/input/adbinput_test.go` 的特殊字符转义测试
- 集成测试：真机上跑通完整交互循环

### 发布与分发

- 设置 GitHub Actions CI（build + test + lint）
- 配置 GoReleaser 自动构建多平台 binary（darwin-arm64/amd64, linux-arm64/amd64）
- 发布到 Claude Code Plugin Marketplace
- 发布到 ClawHub

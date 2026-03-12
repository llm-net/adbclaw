# adbclaw — 第二轮迭代：高级命令

> 目标：增加高层命令，减少 AI agent 交互轮次，解锁中文输入等核心能力。

---

## Phase 1 — 文本输入增强

增强文本输入能力。

### 1.1 `clear-field`

```bash
adbclaw clear-field                    # 清空当前焦点输入框
adbclaw clear-field --index 5          # 聚焦到元素后清空
```

**实现方式**：
- `KEYCODE_MOVE_END` → `Ctrl+A`（全选）→ `KEYCODE_DEL`
- 或 `KEYCODE_MOVE_HOME` + `Shift+KEYCODE_MOVE_END` + `KEYCODE_DEL`

### 1.2 key 别名扩展

补充 `PASTE`(279)、`COPY`(278)、`CUT`(277)、`FORWARD_DEL`(112)、`MOVE_HOME`(122)、`MOVE_END`(123)、`PAGE_UP`(92)、`PAGE_DOWN`(93)、`WAKEUP`(224)、`SLEEP`(223) 等常用别名。

---

## Phase 2 — 导航增强

简化最常见的导航操作。

### 2.1 `open` — 打开 URL/深度链接

```bash
adbclaw open "https://www.baidu.com"           # 打开网页
adbclaw open "weixin://dl/scan"                # 微信扫一扫
adbclaw open "taobao://s.taobao.com?q=手机"    # 淘宝搜索
```

**实现方式**：
- `adb shell am start -a android.intent.action.VIEW -d '<uri>'`
- 返回：启动结果 + 目标 package/activity（如果系统返回）

### 2.2 `scroll` — 智能滚动

```bash
adbclaw scroll down                    # 自动计算坐标，向下滚一屏
adbclaw scroll up                      # 向上滚
adbclaw scroll left / right            # 水平滚动
adbclaw scroll down --index 5          # 在 scrollable 元素内滚动
adbclaw scroll down --pages 3          # 连续滚 N 屏
adbclaw scroll down --distance 500     # 指定滑动像素距离
```

**实现方式**：
1. `adb shell wm size` 获取屏幕尺寸
2. 计算中心点，swipe 屏幕高度的 60%（可配置）
3. `--index` 时获取 UI 树，在 scrollable 元素 bounds 内滚动
4. `--pages` 时循环执行，每次间隔 300ms

**新增文件**：
- `src/cmd/open.go` — open 命令
- `src/cmd/scroll.go` — scroll 命令
- `src/pkg/input/scroll.go` — 滚动计算逻辑

---

## Phase 3 — 状态感知

让 agent 可以等待 UI 变化，而不是反复 observe 轮询。

### 3.1 `wait` — 等待 UI 状态

```bash
adbclaw wait --text "加载完成" --timeout 10000         # 等文字出现
adbclaw wait --text "Loading" --gone --timeout 10000    # 等文字消失
adbclaw wait --id "com.app:id/content" --timeout 5000   # 等元素出现
adbclaw wait --activity ".MainActivity" --timeout 8000   # 等特定 Activity
```

**实现方式**：
- 内部轮询 uiautomator dump（间隔 ~800ms）
- 匹配条件：text/id/content-desc 出现或消失
- activity 模式用 `dumpsys window` 检查
- 超时返回错误 `WAIT_TIMEOUT`，成功返回匹配到的元素
- 默认超时 10000ms

### 3.2 `screen` — 屏幕状态管理

```bash
adbclaw screen status                  # 返回：亮屏/息屏、锁屏状态、旋转方向
adbclaw screen on                      # 亮屏 (KEYCODE_WAKEUP)
adbclaw screen off                     # 息屏 (KEYCODE_SLEEP)
adbclaw screen unlock                  # 亮屏 + 上滑解锁（无密码场景）
adbclaw screen rotation auto|0|1|2|3   # 设置屏幕旋转
```

**实现方式**：
- `status`: `dumpsys power | grep 'Display Power'` + `dumpsys window policy | grep 'showing='`
- `on/off`: `input keyevent WAKEUP/SLEEP`
- `unlock`: WAKEUP + `input swipe`（从底部上滑）
- `rotation`: `settings put system accelerometer_rotation 0/1` + `settings put system user_rotation N`

**新增文件**：
- `src/cmd/wait.go` — wait 命令
- `src/cmd/screen.go` — screen 子命令
- `src/pkg/device/screen.go` — 屏幕状态查询实现

---

## Phase 4 — App 管理增强

补全应用生命周期管理。

```bash
adbclaw app install ./app.apk             # 安装
adbclaw app install ./app.apk --replace   # 覆盖安装
adbclaw app uninstall com.example.app     # 卸载
adbclaw app clear com.example.app         # 清除应用数据
```

**实现方式**：
- `install`: `adb install [-r] <path>`
- `uninstall`: `adb uninstall <package>`
- `clear`: `adb shell pm clear <package>`

**修改文件**：
- `src/cmd/app.go` — 新增 install/uninstall/clear 子命令

---

## Phase 5 — 通用能力

逃生舱口和文件操作。

### 5.1 `shell` — 执行任意 ADB shell 命令

```bash
adbclaw shell "settings get system screen_brightness"
adbclaw shell "pm grant com.app android.permission.CAMERA"
```

**实现方式**：
- 直接转发到 `adb shell`，输出包在 JSON envelope 中
- stdout/stderr 分别捕获

### 5.2 `file push/pull`

```bash
adbclaw file push ./test.txt /sdcard/Download/
adbclaw file pull /sdcard/log.txt ./local.txt
```

**实现方式**：
- `adb push <local> <remote>`
- `adb pull <remote> <local>`

**新增文件**：
- `src/cmd/shell.go` — shell 命令
- `src/cmd/file.go` — file 子命令

---

## 实施顺序

| Phase | 内容 | 核心价值 |
|-------|------|---------|
| 1 | clear-field + key 别名 | 增强文本输入 |
| 2 | open + scroll | 简化导航操作 |
| 3 | wait + screen | 减少 observe 轮次 |
| 4 | app install/uninstall/clear | 补全 App 管理 |
| 5 | shell + file | 通用逃生舱 |

每个 Phase 完成后：`make test` + `make build` 验证，更新 SKILL.md 命令文档。

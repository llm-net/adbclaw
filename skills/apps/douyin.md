# 抖音 (Douyin)

- 包名: `com.ss.android.ugc.aweme`
- Scheme: `snssdk1128://`

## 深度链接

优先使用深度链接，可跳过多步 UI 操作，且天然支持中文参数（adb input text 不支持中文）。

```bash
# 通用调用方式
adb shell am start -a android.intent.action.VIEW -d '{link}'
```

| 动作 | 链接 | 参数说明 |
|------|------|----------|
| 搜索 | `snssdk1128://search/result?keyword={keyword}&type={type}` | type: 0=综合, 1=直播, 2=视频, 3=用户 |
| 用户主页 | `snssdk1128://user/profile/{user_id}` | user_id 为数字 ID |
| 直播间 | `snssdk1128://live?room_id={room_id}` | |
| 视频详情 | `snssdk1128://detail/{video_id}` | |

## 已知布局

### 首页（推荐 Feed）

```
┌─────────────────────────────────────────────────┐
│ [侧栏]  经验│免费看剧│游戏│...│推荐(默认)  [搜索] │  ← 顶部导航
├─────────────────────────────────────────────────┤
│                                                   │
│                 视频内容区                          │
│                                                   │
│                              [头像]               │
│                              [关注]               │
│                              [点赞]               │
│                              [评论]               │
│                              [收藏]               │
│                              [分享]               │
│  @用户名                                          │
│  视频描述文字...                                   │
├─────────────────────────────────────────────────┤
│    首页  │  朋友  │  [拍摄]  │  消息  │  我        │  ← 底部导航
└─────────────────────────────────────────────────┘
```

**关键元素定位**：
- 搜索按钮: `content_desc="搜索"`，右上角
- 侧栏按钮: `content_desc` 含 "侧边栏"，左上角
- 底部导航项: `resource_id` 含 `0tn`，text 为 "首页"/"朋友"/"消息"/"我"
- 用户头像: `content_desc` 含用户名，resource_id 含 `user_avatar`
- 关注按钮: text="关注"
- 顶部 Tab: resource_id 含 `5j4`，水平可滚动

### 搜索结果页

```
┌─────────────────────────────────────────────────┐
│ [<]  搜索关键词                          [搜索]   │
├─────────────────────────────────────────────────┤
│  综合 │ 直播 │ 商品 │ 视频 │ 图文 │ 用户 │ ...    │  ← 分类 Tab
├─────────────────────────────────────────────────┤
│                                                   │
│                 搜索结果列表                        │
│                                                   │
└─────────────────────────────────────────────────┘
```

## 设备差异

### Pad（短边 >= 1200px）

- 横屏为主要使用模式，顶部导航 Tab 更多可见
- 搜索结果页：直播 Tab 下为单列大卡片，纵向滚动浏览
- 首页推荐 Feed 视频区域在中间，右侧互动按钮（点赞/评论/分享）
- 底部导航栏居中显示，不铺满全屏宽度
- 屏幕坐标系：横屏时宽 > 高（如 3200x2136）

### Phone（短边 < 1200px）

- 竖屏为主要使用模式
- 搜索结果页：直播 Tab 下可能为双列小卡片
- 首页为全屏沉浸式视频
- 底部导航栏铺满屏幕宽度
- 屏幕坐标系：宽 < 高（如 1080x2400）

## 常见工作流

### 搜索内容

```
1. adb shell am start -a android.intent.action.VIEW \
   -d 'snssdk1128://search/result?keyword=遥控车&type=0'
2. 等待 3 秒加载
3. 根据需要点击分类 Tab（直播/视频/用户等）
```

不要尝试用 `adbclaw type` 输入中文，直接用深度链接。

### 搜索直播

```
1. adb shell am start -a android.intent.action.VIEW \
   -d 'snssdk1128://search/result?keyword={关键词}&type=1'
2. 等待 3 秒
3. 页面直接展示直播结果，纵向滚动浏览
```

### 浏览推荐 Feed

```
1. adbclaw app launch com.ss.android.ugc.aweme
2. 等待 3 秒
3. 上滑切换下一个视频:
   - Pad:  adbclaw swipe 1600 1500 1600 400 --duration 300
   - Phone: adbclaw swipe 540 1800 540 600 --duration 300
```

### 获取当前视频信息

```
1. adbclaw tap {屏幕中心}     → 暂停视频（重要！否则 UI dump 会失败）
2. adbclaw ui tree             → 获取 UI 元素
3. 查找 content_desc 含用户名的元素、desc 含视频描述的元素
```

## 已知问题

### UI dump 在视频播放时失败

**现象**: `adbclaw ui tree` 或 `adbclaw observe` 返回 `UI_DUMP_FAILED`。

**原因**: 视频播放动画导致 uiautomator dump 超时。

**解决**: 先 tap 屏幕中心暂停视频，等待 1 秒后再 dump。

```bash
adbclaw tap 1600 1046    # Pad 屏幕中心（根据 device info 调整）
sleep 1
adbclaw ui tree
```

### 中文输入不可用

**现象**: `adbclaw type "中文"` 报错或输入乱码。

**原因**: `adb shell input text` 不支持非 ASCII 字符。

**解决**: 所有涉及中文输入的场景，使用深度链接代替手动输入。

### 搜索框残留文本

**现象**: 打开搜索页时输入框可能有预填的推荐文字。

**解决**:
1. 如果通过深度链接搜索则无此问题
2. 如果必须手动输入，先点击清除按钮（输入框右侧 X 图标），或全选后删除:
   ```bash
   adb shell input keyevent KEYCODE_MOVE_END
   adb shell input keyevent --longpress KEYCODE_DEL KEYCODE_DEL ...
   ```

### 新安装首次启动有引导页

**现象**: 首次打开抖音会有登录/权限引导页，不是直接进入 Feed。

**解决**: 观察屏幕，查找 "跳过"、"以后再说"、"同意" 等按钮并点击。

---

> 测试信息: Xiaomi Pad (25097RP43C), Android 16 (SDK 36), 屏幕 2136x3200, 密度 440dpi。App 版本: 2026 年 3 月安装。

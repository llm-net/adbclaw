export default {
  nav: {
    features: '特性',
    install: '安装',
    commands: '命令',
    usage: '用法',
  },
  hero: {
    title: '你在 Android 上的',
    titleHighlight: '眼、手与耳',
    description:
      '为 Agent、Claw、Bot 和 LLM 而生。30+ 结构化命令 — 截屏观察、元素精准定位、视频播放中读取直播弹幕、采集系统音频用于语音识别、全生命周期应用管理。JSON 输入，JSON 输出。连接物理世界的桥梁。',
    getStarted: '安装 Skill',
    seeExamples: '查看示例',
    versionNote: '音频采集 + monitor + 持续迭代中',
  },
  features: {
    label: '核心能力',
    title: 'Agent 的超级能力',
    description: '纯工具层。不含 LLM 逻辑，不含 Agent 框架。只有可靠的、结构化的命令，赋予你在 Android 设备上的感知与操控能力。',
    items: [
      {
        title: '直播流智能感知',
        description:
          'monitor 直连 Android 无障碍框架，实时读取所有 UI 文本 — 即使在视频播放和直播期间（此时 uiautomator dump 会超时挂起）。弹幕、字幕、动态内容，其他工具看不到的你都能看到。',
        icon: 'zap',
      },
      {
        title: '系统音频采集',
        description:
          'audio capture 通过 REMOTE_SUBMIX（Android 11+）录制设备音频，WAV 流输出到 stdout。管道接入 ASR 即可语音转文字。配合 monitor，视觉文本 + 音频双通道全感知。',
        icon: 'layers',
      },
      {
        title: '结构化 JSON API',
        description:
          '每条命令返回 {ok, command, data, error, duration_ms, timestamp}。可靠解析，无需猜测。错误包含错误码和可操作的建议。',
        icon: 'json',
      },
      {
        title: '智能元素定位',
        description:
          '通过元素索引、resource ID 或文本内容进行点击、长按或滚动。无需猜测坐标。UI 树已索引，包含边界和中心点，支持像素级精确操作。',
        icon: 'grid',
      },
      {
        title: '深度链接导航',
        description:
          '通过 URI 直接跳转到任意应用页面。打开网页、触发微信扫一扫、搜索淘宝 — 完全跳过多步导航。一条命令，即刻到达。CJK 文本输入的关键方案。',
        icon: 'link',
      },
      {
        title: '等待状态变化',
        description:
          '阻塞等待 UI 元素出现或消失，或 Activity 加载完成。无需在 Agent 代码中编写轮询循环。可配置超时和间隔。成功时返回匹配的元素。',
        icon: 'clock',
      },
      {
        title: '完整设备控制',
        description:
          '30+ 条命令覆盖屏幕观察、输入注入、智能滚动、应用生命周期、屏幕管理、shell 访问和文件传输。一切通过标准 ADB。无需安装 APK。',
        icon: 'package',
      },
      {
        title: 'App 知识档案',
        description:
          '预置热门应用（抖音、美团等）的操作档案，包含深度链接、UI 布局和已知问题。加载一次，跳过试错。每次发版都有新 Profile。',
        icon: 'book',
      },
      {
        title: '持续进化',
        description:
          '新能力持续交付中。monitor、音频采集、App 档案 — 每次发版都扩展你的感知和控制范围。安装一次，自动获得新能力。为 Agent 打造，由 Agent 驱动。',
        icon: 'bot',
      },
    ],
  },
  install: {
    label: '安装',
    title: '几秒即可上手',
    description: '提供 macOS 和 Linux 预编译二进制文件。无需 Go 工具链。下载即用。',
    recommended: '推荐',
    oneLiner: '一键安装',
    oneLinerDesc: '自动检测操作系统和架构。将最新二进制文件下载到',
    manual: '手动安装',
    downloadBinary: '下载二进制文件',
    downloadBinaryDesc: '从 GitHub Releases 下载适合你平台的预编译二进制文件。',
    fromSource: '从源码构建',
    buildWithGo: '使用 Go 构建',
    buildWithGoDesc: '克隆仓库并构建。需要 Go 1.24+。',
    prerequisite: '前置条件：',
    prerequisiteText: 'ADB（Android Debug Bridge）必须已安装并在 PATH 中。',
  },
  howItWorks: {
    label: '架构',
    title: '工作原理',
    description: '命令从 AI Agent 经由 adb-claw 流向设备。每个响应都是结构化 JSON。',
    agentLoop: '推荐的 Agent 循环',
    architectureSteps: [
      {
        label: 'AI Agent',
        sublabel: 'Claude / OpenClaw / LLM',
        description: '读取技能描述，发送结构化命令，解析 JSON 响应以决定下一步操作',
      },
      {
        label: 'adb-claw',
        sublabel: 'Go CLI · v1.5.4',
        description: '将 30+ 条命令转换为 ADB 操作。返回带有错误码和建议的结构化 JSON',
      },
      {
        label: 'ADB',
        sublabel: 'USB / WiFi',
        description: '将 shell 命令、截屏和文件传输传送到 Android 设备',
      },
      {
        label: '设备',
        sublabel: 'Android',
        description: '执行操作 — UI 转储、截屏、输入事件、应用管理、文件 I/O',
      },
    ],
    agentWorkflow: [
      { step: '01', action: '观察', detail: '运行 observe 命令，一次调用同时获取截屏和 UI 树' },
      { step: '02', action: '决策', detail: 'AI Agent 分析屏幕状态并规划下一步操作' },
      { step: '03', action: '执行', detail: '点击、滚动、打开深度链接或输入文字 — 通过元素索引' },
      { step: '04', action: '等待', detail: '使用 wait 命令阻塞等待 UI 状态变化，然后重新观察' },
    ],
  },
  codeDemo: {
    label: '用法',
    title: '30+ 条命令，一个二进制',
    description: '观察、导航、等待、管理 — 全部作为顶级命令，输出结构化 JSON。优先使用元素索引而非坐标。',
    jsonEnvelope: 'JSON 信封',
    everyCommand: '每条命令都返回此格式',
    examples: [
      {
        title: '观察与检查',
        commands: [
          { cmd: 'adb-claw observe', comment: '截屏 + UI 树' },
          { cmd: 'adb-claw screenshot --width 720', comment: '缩放截屏' },
          { cmd: 'adb-claw ui tree', comment: '索引化元素' },
          { cmd: 'adb-claw ui find --text "Login"', comment: '按文本查找' },
        ],
      },
      {
        title: '输入与导航',
        commands: [
          { cmd: 'adb-claw tap --index 5', comment: '按元素索引点击' },
          { cmd: 'adb-claw type "hello world"', comment: '输入文本' },
          { cmd: 'adb-claw scroll down --pages 3', comment: '智能滚动' },
          { cmd: 'adb-claw open "weixin://dl/scan"', comment: '深度链接' },
          { cmd: 'adb-claw clear-field --index 2', comment: '清空输入框' },
        ],
      },
      {
        title: '等待与屏幕',
        commands: [
          { cmd: 'adb-claw wait --text "Done"', comment: '等待元素出现' },
          { cmd: 'adb-claw wait --text "Loading" --gone', comment: '等待元素消失' },
          { cmd: 'adb-claw screen status', comment: '亮屏/灭屏/锁定/旋转' },
          { cmd: 'adb-claw screen unlock', comment: '唤醒 + 滑动解锁' },
        ],
      },
      {
        title: '监控与音频',
        commands: [
          { cmd: 'adb-claw monitor --stream', comment: '实时 UI 文本（视频安全）' },
          { cmd: 'adb-claw monitor --duration 30000', comment: '30 秒定时采集' },
          { cmd: 'adb-claw audio capture --file out.wav', comment: '录制系统音频' },
          { cmd: 'adb-claw audio capture --stream | asrclaw transcribe', comment: '管道接 ASR' },
        ],
      },
      {
        title: '应用与系统',
        commands: [
          { cmd: 'adb-claw app launch com.example', comment: '启动应用' },
          { cmd: 'adb-claw app install ./app.apk', comment: '安装 APK' },
          { cmd: 'adb-claw shell "pm list packages"', comment: '原始 shell' },
          { cmd: 'adb-claw file pull /sdcard/log.txt .', comment: '拉取文件' },
        ],
      },
    ],
  },
  commandTree: {
    label: '参考',
    title: '完整命令参考',
    description: '每条命令返回结构化 JSON。所有命令支持',
    commands: [
      {
        category: '观察',
        items: [
          { cmd: 'observe', desc: '并行截屏 + UI 树', flags: '--width' },
          { cmd: 'screenshot', desc: '屏幕截图（base64 或文件）', flags: '--file, --width' },
          { cmd: 'ui tree', desc: '索引化 UI 元素树' },
          { cmd: 'ui find', desc: '按文本/ID/索引查找元素', flags: '--text, --id, --index' },
        ],
      },
      {
        category: '输入',
        items: [
          { cmd: 'tap', desc: '按坐标或元素点击', flags: '--index, --id, --text' },
          { cmd: 'long-press', desc: '长按（可设时长）', flags: '--duration' },
          { cmd: 'swipe', desc: '在坐标间滑动', flags: '--duration' },
          { cmd: 'key', desc: '按键（30+ 别名）', flags: 'HOME, BACK, ENTER...' },
          { cmd: 'type', desc: '输入 ASCII 文本' },
          { cmd: 'clear-field', desc: '清空聚焦的输入框', flags: '--index, --id, --text' },
        ],
      },
      {
        category: '导航',
        items: [
          { cmd: 'scroll', desc: '任意方向智能滚动', flags: '--pages, --distance, --index' },
          { cmd: 'open', desc: '打开 URI / 深度链接' },
        ],
      },
      {
        category: '状态',
        items: [
          { cmd: 'wait', desc: '等待 UI 元素或 Activity', flags: '--text, --id, --gone, --timeout' },
          { cmd: 'screen status', desc: '屏幕亮灭、锁定、旋转状态' },
          { cmd: 'screen on/off', desc: '亮屏或灭屏' },
          { cmd: 'screen unlock', desc: '唤醒 + 滑动解锁' },
          { cmd: 'screen rotation', desc: '设置旋转模式', flags: 'auto, 0-3' },
        ],
      },
      {
        category: '应用',
        items: [
          { cmd: 'app list', desc: '已安装应用列表', flags: '--all' },
          { cmd: 'app current', desc: '前台应用包名/Activity' },
          { cmd: 'app launch', desc: '按包名启动应用' },
          { cmd: 'app stop', desc: '强制停止应用' },
          { cmd: 'app install', desc: '安装 APK', flags: '--replace' },
          { cmd: 'app uninstall', desc: '卸载应用' },
          { cmd: 'app clear', desc: '清除应用数据' },
        ],
      },
      {
        category: '感知',
        items: [
          { cmd: 'monitor', desc: '无障碍实时 UI 文本（视频安全）', flags: '--duration, --interval, --stream' },
          { cmd: 'audio capture', desc: '系统音频 → WAV 流（Android 11+）', flags: '--file, --duration, --rate, --stream' },
        ],
      },
      {
        category: '系统',
        items: [
          { cmd: 'device list', desc: '已连接设备' },
          { cmd: 'device info', desc: '型号、屏幕、SDK 版本' },
          { cmd: 'shell', desc: '执行原始 ADB shell 命令' },
          { cmd: 'file push', desc: '推送文件到设备' },
          { cmd: 'file pull', desc: '从设备拉取文件' },
          { cmd: 'doctor', desc: '环境健康检查' },
          { cmd: 'skill', desc: '输出 skill.json 供 Agent 使用' },
        ],
      },
    ],
  },
  relatedProjects: {
    label: '生态',
    title: '相关项目',
    description: 'Android 自动化和 AI Agent 领域的其他工具。',
    skillPlatform: 'Skill 平台',
    items: [
      {
        name: 'OpenClaw',
        url: 'https://github.com/openclaw/openclaw',
        description: '本地优先的个人 AI 助手平台。adb-claw 作为 OpenClaw Skill 发布在 ClawHub 上。',
        stars: '',
        highlight: true,
      },
      {
        name: 'mobile-use',
        url: 'https://github.com/anthropics/mobile-use',
        description: 'Anthropic 的 AI Agent，用于控制真实移动设备，AndroidWorld 基准测试表现顶尖。',
        stars: '2.2k',
      },
      {
        name: 'DroidRun',
        url: 'https://github.com/droidrun/droidrun',
        description: '基于 LLM 的 Android 设备控制框架，支持多模型。',
        stars: '7.7k',
      },
      {
        name: 'scrcpy',
        url: 'https://github.com/Genymobile/scrcpy',
        description: 'Android 屏幕镜像和控制的标杆工具。',
        stars: '136k',
      },
    ],
  },
  footer: {
    description: 'Android 设备控制 CLI，用于 AI Agent 自动化。30+ 条 ADB 命令。纯工具层 — 不含 LLM/Agent 逻辑。',
    project: '项目',
    documentation: '文档',
    issues: '问题',
    releases: '版本',
    availableOn: '可用平台',
    claudeCodePlugin: 'Claude Code 插件',
    openClawClawHub: 'OpenClaw / ClawHub',
    standaloneCli: '独立 CLI',
    stack: '技术栈',
  },
}

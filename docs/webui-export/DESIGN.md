# AniFeed Web UI Design Export

版本：v0.03.1

## 1. 设计定位

AniFeed 是一个自托管动漫 RSS → 115 离线管理工具。Web UI 采用轻量后台管理风格：

- 浅色为主，支持深色主题
- 左侧固定导航 + 顶部标题区
- 卡片式内容区
- 高密度表格与表单
- 移动端底部导航
- 强调状态、错误和操作反馈

## 2. 技术栈

- 纯 HTML / CSS / Vanilla JavaScript
- 无前端框架
- 静态资源由 Go `embed` 打包进单二进制
- API 同源访问
- 登录保护、CSP、X-Frame-Options 由 Go 服务端设置

## 3. 文件结构

```text
assets/
  index.html   页面骨架、侧边栏、移动端导航
  styles.css   设计令牌、布局、组件、响应式
  app.js       路由、页面渲染、API 交互、弹窗

screenshots/
  dashboard-desktop.png
  feeds-desktop.png
  tasks-desktop.png
  settings-desktop.png
  ai-desktop.png
  logs-desktop.png
  dashboard-mobile.png
  ai-mobile.png
  settings-mobile.png
```

注意：`index.html` 通过 `/app.js` 和 `/styles.css` 加载资源，API 由 AniFeed Go 服务提供，不能作为纯静态文件直接双击运行。

## 4. 设计令牌

### 浅色主题

```text
背景        #f4f6fb
卡片        #ffffff
次级表面    #f8fafc
边框        #e3e8f2
正文        #0f172a
次级文字    #334155
弱文字      #64748b
主色        #4f46e5
成功        #10b981
警告        #f59e0b
危险        #ef4444
信息        #0ea5e9
```

### 深色主题

通过 `data-theme="dark"` 切换，主要令牌：

```text
背景        #0a0f1c
卡片        #111a2c
三级表面    #1a2438
边框        #223049
正文        #e9eef9
主色        #7c83ff
```

### 字体

```text
界面：Inter, Noto Sans SC, PingFang SC, Microsoft YaHei
等宽：JetBrains Mono, SFMono-Regular, Consolas
```

### 圆角与阴影

```text
大圆角      14px
小圆角      10px
卡片阴影    0 1px 2px rgba(15,23,42,.04), 0 10px 28px -16px rgba(15,23,42,.22)
```

## 5. 布局

### 桌面

```text
+----------------------+--------------------------------------+
| Sidebar 216px        | Topbar                               |
| Logo                 +--------------------------------------+
| 概览                 | Page title / subtitle                |
| 订阅管理             | Toolbar                              |
| 离线任务             | Content cards / tables               |
| 历史记录             |                                      |
| 运行日志             |                                      |
| 设置                 |                                      |
| AI 助手              |                                      |
+----------------------+--------------------------------------+
```

### 移动端

- 隐藏侧边栏
- 顶部标题保留
- 底部固定导航：

```text
概览 / 订阅 / 任务 / 日志 / 设置 / AI
```

- 键盘弹出时隐藏底部导航，避免被顶到键盘上方
- AI 页面进入时不自动聚焦输入框

## 6. 页面

### 概览

- 服务状态
- 番剧卡片
- 最近任务
- 运行统计

### 订阅管理

- 搜索、来源筛选、状态筛选
- 订阅海报卡片
- 批量管理
- 编辑弹窗：基本 / 备用 RSS / 自定义
- 自动备用 RSS、TMDB 自动偏移开关

### 离线任务

- 缓存优先渲染
- 刷新按钮拉取最新数据
- 默认最近 20 条
- 加载更多
- 状态筛选
- 清理任务

### 历史记录

- 按订阅筛选
- 时间、标题、订阅、Hash
- 分页

### 运行日志

- 列表 / 终端模式
- 级别筛选
- 搜索
- 自动刷新
- 深色终端风格

### 设置

分组：

```text
连接与账号
下载器
整理与命名
字幕组
通知
高级
关于
```

### AI 助手

- 固定高度聊天区
- 消息区独立滚动
- 推荐问题
- 清空对话
- 只读 / 完整权限
- 修改操作二次确认

## 7. 组件规范

### 卡片

- `.card`
- 标题 `.card-head`
- 内容 `.section-body`
- 分组标题 `.section-title`

### 按钮

- `.btn` 默认
- `.btn-primary` 主操作
- `.btn-danger` 危险操作
- `.btn-sm` 小尺寸
- `.btn-ghost` 弱操作
- `.btn-icon-only` 图标按钮

### 表单

- `.field`
- `.form-grid`
- `.input`
- `.select`
- `.textarea`
- `.switch`
- `.check`
- 帮助按钮 `?` 使用 `.help-dot`

### 状态

- `.badge-success`
- `.badge-warning`
- `.badge-danger`
- `.badge-info`
- `.notice`
- `.banner`

### 表格

- `.table`
- `.table-wrap`
- 移动端自动切换为卡片
- 使用 `data-label` 显示字段名

## 8. 响应式断点

```text
700px：移动端导航、表格卡片化、设置标签横向滚动
640px：概览最近任务卡片化、任务表卡片化
```

## 9. 交互约定

- 页面切换使用 hash 路由
- 使用 `state.routeEpoch` 丢弃过期异步结果
- 已加载页面优先渲染缓存，再静默刷新
- 危险操作必须二次确认
- 设置保存统一写入 `config.toml`
- 错误通过 toast 或表单错误展示
- AI 修改/删除操作必须二次确认

## 10. 导出说明

本目录为设计源码导出，不包含：

- 后端 Go 代码
- `config.toml`
- Cookies / AI Key / 密码
- 用户数据
- Docker / FPK 打包文件

要查看实际效果，请运行 AniFeed 服务并访问：

```text
http://服务器地址:80
```

或 FPK 安装后的：

```text
http://fnOS地址:8115
```

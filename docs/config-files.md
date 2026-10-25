# 配置文件说明

`anifeed` 支持统一的 TOML 配置文件和传统的 JSON 配置文件。新用户推荐使用 `config.toml` 统一配置，现有用户可以继续使用原有配置文件，无需修改。

## 统一 TOML 配置（推荐）

### 配置文件位置

`config.toml` 查找顺序：

1. 当前工作目录的 `config.toml`
2. `~/.config/anifeed/config.toml`

### 完整配置示例

```toml
# 认证配置
[auth]
# 方式 1: 指定 cookies 文件路径（相对路径基于 config.toml 所在目录）
cookies_file = ".cookies"

# 方式 2: 直接配置 cookies 字符串（不推荐，建议使用 cookies_file）
# cookies = "UID=xxx; CID=xxx; SEID=xxx; KID=xxx"

# 服务器配置
[server]
host = "0.0.0.0" # 监听地址，默认 0.0.0.0（所有网卡）
port = 8115  # HTTP 服务器端口，默认 8115
# access_token = "change-me"  # 可选：Web UI / API 访问令牌

# 数据库配置
[database]
path = "db.sqlite"  # SQLite 数据库文件路径，默认 db.sqlite

# 115 云存储配置
[p115]
disable_cache = false    # 是否禁用缓存，默认 false
chunk_delay = 2          # 分块延迟（秒），默认 2
chunk_size = 200         # 分块大小，默认 200
cooldown_min_ms = 1000   # 冷却最小时间（毫秒），默认 1000
cooldown_max_ms = 1100   # 冷却最大时间（毫秒），默认 1100

# 代理配置
[proxy]
http = "http://127.0.0.1:10809"  # HTTP 代理地址，默认 http://127.0.0.1:10809

# 订阅源与网络
[subscription]
mikan_host = "https://mikanani.me" # Mikan 站点/镜像地址
rss_timeout_seconds = 20           # RSS 抓取超时，默认 20 秒

# 默认离线设置
[offline]
default_cid = ""                   # 默认 115 目录 ID
default_savepath = "番剧/${title}" # 默认保存路径模板，留空使用 115 默认目录
global_exclude = ["720p", "/\\.5$/"] # 全局排除，每项支持 /正则/
skip_x5 = true                     # 跳过 12.5 这类 x.5 集
retry = 2                          # 115 提交失败重试次数
auto_clear_completed = false       # 全部/自动运行后清理已完成任务（不删文件）
dual_folder = false                 # 离线后转移到保存文件夹：完成后移动到 final_cid
final_cid = ""                      # 离线后转移到保存文件夹的最终目录 cid
one_per_episode = false             # 默认只识别记录；开启后同一集只提交评分最高的一个
keep_only_video = false             # 离线后转移到保存文件夹只移动视频和字幕文件
auto_language = ["chs", "bilingual", "cht"]
auto_subtitle = ["softsub", "hardsub", "external"]
auto_container = ["mkv", "mp4"]
auto_resolution = ["r1080", "r2160", "r720"]

# 重命名
[rename]
enabled = false
template = "[${subgroup}] ${title} S${seasonFormat}E${episodeFormat}"
season = 1
offset = 0

# TMDB
[tmdb]
api = "https://api.themoviedb.org"
api_key = ""
image = "https://image.tmdb.org"

# 自动运行
[auto]
enabled = false                    # 是否定时运行所有启用的订阅
interval_minutes = 30              # 间隔（分钟）
run_on_start = false               # 服务启动后立即运行一次

# 通知（运行摘要推送到 Webhook）
[notify]
type = "auto"                       # auto | generic | wecom
webhook = ""
content = ["summary", "details", "errors", "time"]
on_success = true
on_error = true

# RSS 订阅配置（可配置多个）
[[rss]]
site = "mikanani.me"
name = "测试订阅"
source = "mikan"                   # mikan | anibt | animegarden | other
enabled = true                     # 是否参与自动/全部运行
filter = "简体内嵌"
url = "https://mikanani.me/RSS/Bangumi?bangumiId=2739&subgroupid=12"

[[rss]]
site = "share.dmhy.org"
name = "水星的魔女"
savepath = "文件夹名称"
filter = "简日双语"
exclude = "720p"            # 可选：排除规则（黑名单），匹配的标题不会被离线
url = "https://share.dmhy.org/topics/rss/rss.xml?keyword=..."
cid = "123456"           # 可选：115 目录 ID
expiration = 3600        # 可选：缓存过期时间（秒）

# 站点特定配置
[sites."share.dmhy.org"]
https_agent = true  # 是否使用代理

[sites."nyaa.si"]
https_agent = true

[sites."115.com".headers]
Cookie = "UID=xxx; CID=xxx; SEID=xxx; KID=xxx"

# 兼容迁移开关（Web UI 首次保存时自动写入）
[compat]
disable_legacy = true  # 不再回退读取 rss.json / node-site-config.json / .cookies
```

### 配置字段说明

#### [auth] 认证配置

- `cookies_file`: cookies 文件路径，支持相对路径（相对于 config.toml 所在目录）和绝对路径
- `cookies`: 直接配置 cookies 字符串（不推荐，建议使用 cookies_file）

#### [server] 服务器配置

- `host`: 监听地址，默认 `0.0.0.0`；设为 `127.0.0.1` 时只允许本机访问
- `port`: HTTP 服务器端口，有效范围 1-65535，默认 8115
- `access_token`: 可选。设置后 Web UI 和 API 需要携带 `X-Access-Token` 请求头

#### [compat] 兼容配置

- `disable_legacy`: 设为 `true` 后不再回退读取 `rss.json`、`node-site-config.json`、`.cookies`。Web UI 首次保存时会自动写入，避免删空订阅后旧配置「复活」。

#### [database] 数据库配置

- `path`: SQLite 数据库文件路径，支持相对路径（相对于 config.toml 所在目录）和绝对路径，默认 db.sqlite

#### [p115] 115 云存储配置

- `disable_cache`: 是否禁用缓存，默认 false
- `chunk_delay`: 分块延迟（秒），默认 2
- `chunk_size`: 分块大小，默认 200
- `cooldown_min_ms`: 冷却最小时间（毫秒），默认 1000
- `cooldown_max_ms`: 冷却最大时间（毫秒），默认 1100

#### [proxy] 代理配置

- `http`: HTTP 代理地址，默认 http://127.0.0.1:10809

#### [subscription] 订阅源与网络

- `mikan_host`: Mikan 站点地址，Web UI 生成 Mikan RSS 时使用，默认 `https://mikanani.me`
- `rss_timeout_seconds`: 抓取 RSS 的超时秒数，默认 20

#### [offline] 默认离线设置

- `default_cid`: 订阅未填 `cid` 时使用的 115 目录 ID
- `default_savepath`: 订阅未填 `savepath` 时使用的保存路径模板。支持 `${title}`、`${site}`、`${year}`、`${month}`、`${day}`、`${season}`；替换值会自动清理路径非法字符
- `global_exclude`: 全局排除规则数组。每项支持普通字符串或 `/正则/`，对所有订阅生效
- `skip_x5`: 是否自动跳过标题中的 x.5 集（例如 `12.5`）
- `retry`: 115 离线提交失败后的重试次数，默认 0
- `auto_clear_completed`: 全部/自动运行结束后是否调用 115 的「清除已完成任务」（不会删除网盘文件）
- `dual_folder`: 是否启用离线后转移到保存文件夹。开启后任务先下载到 `default_cid`，完成并重命名后移动到 `final_cid`
- `final_cid`: 离线后转移到保存文件夹的目标目录 ID；`dual_folder = true` 时必填
- `one_per_episode`: 是否每集只提交一个版本，默认 `false`；关闭时只识别并记录，不实际跳过

#### [rename] 重命名

- `enabled`: 是否对 115 已完成离线任务重命名
- `template`: 重命名模板，默认 `[${subgroup}] ${title} S${seasonFormat}E${episodeFormat}`
- `season`: 季数，默认 1；用于 `${season}` / `${seasonFormat}`
- `offset`: 集数偏移，默认 0；例如标题第 00 集实际为第 01 集时填 1
- 模板支持 `${tmdbid}`

#### [tmdb] TMDB

- `api`: TMDB API 地址，默认 `https://api.themoviedb.org`
- `api_key`: TMDB API Key
- `image`: TMDB 图片地址，默认 `https://image.tmdb.org`

模板变量：`${title}`、`${subgroup}`、`${season}`、`${seasonFormat}`、`${episode}`、`${episodeFormat}`、`${itemTitle}`、`${resolution}`。

启用重命名或离线后转移后，后台每 20 秒检查一次已完成任务；解析不到集数的任务不会移动，并记录错误。自动清理会暂停，避免任务在被处理前清除。

#### [auto] 自动运行

- `enabled`: 是否启用定时运行
- `interval_minutes`: 运行间隔（分钟），默认 30
- `run_on_start`: 服务启动后是否立即运行一次

自动运行只处理 `enabled = true` 的订阅；115 未登录或已有运行中的作业时会跳过并记录日志。

#### [notify] 通知

- `type`: `auto`（自动识别企业微信）、`generic`（通用 JSON Webhook）或 `wecom`（企业微信机器人）
- `webhook`: 运行结束后 POST 的地址，留空关闭；企业微信填机器人 Webhook URL
- `content`: 纯文本通知包含的内容，可选 `summary`（运行摘要）、`details`（订阅明细）、`errors`（错误）、`time`（时间）
- `on_success`: 成功时是否推送
- `on_error`: 失败时是否推送

Webhook payload 包含 `event`、`title`、`status`、`feed_count`、`matched`、`added`、`error`、`timestamp` 和 `results`。

#### [[rss]] RSS 订阅配置

每个 `[[rss]]` 块定义一个订阅，必填字段：

- `site`: 站点域名（用于分组）
- `name`: 订阅名称
- `url`: RSS 订阅地址

可选字段：

- `source`: 订阅来源，可选 `mikan`、`anibt`、`animegarden`、`other`；留空时按 URL 域名自动识别
- `enabled`: 是否参与自动运行 / 全部运行，默认 `true`
- `cid`: 115 云存储目录 ID
- `savepath`: 保存路径
- `filter`: 内容过滤规则（白名单，匹配的标题会被离线）
- `exclude`: 内容排除规则（黑名单，匹配的标题不会被离线）
- `expiration`: 缓存过期时间（秒）

#### [sites."<域名>"] 站点配置

- `https_agent`: 是否为该站点启用代理
- `headers`: 自定义 HTTP 请求头

### 配置优先级规则

配置值按以下优先级合并（从高到低）：

1. **命令行参数**：`--cookies`、`--rss`、`--port` 等
2. **TOML 配置**：`config.toml` 中的配置
3. **传统配置文件**：`rss.json`、`node-site-config.json`、`.cookies`
4. **程序默认值**：内置默认配置

示例：如果同时存在 `config.toml` 中的 `[server].port = 8115` 和命令行参数 `--port 9000`，则使用命令行参数的值 9000。

### 相对路径解析

`[auth].cookies_file` 和 `[database].path` 中的相对路径基于 `config.toml` 所在目录解析：

```toml
# 如果 config.toml 位于 ~/.config/anifeed/config.toml
[auth]
cookies_file = ".cookies"  # 解析为 ~/.config/anifeed/.cookies

[database]
path = "data/app.db"  # 解析为 ~/.config/anifeed/data/app.db
```

## 传统配置文件（向后兼容）

现有用户可以继续使用传统配置文件，无需迁移到 `config.toml`。

### 默认查找顺序

`rss.json`：

1. 当前工作目录的 `rss.json`
2. `~/.config/anifeed/rss.json`

`.cookies`：

1. 当前工作目录的 `.cookies`
2. `~/.config/anifeed/.cookies`

`node-site-config.json`：

1. 当前工作目录的 `node-site-config.json`
2. `~/.config/anifeed/node-site-config.json`
3. `~/node-site-config.json`

`~/node-site-config.json` 是旧行为的兼容路径。

### 显式指定 RSS 配置

命令行参数 `--rss` / `-r` 保持原有语义：传入什么路径就读取什么路径，不走默认查找顺序。

```bash
anifeed --rss /path/to/rss.json
```

### 示例目录

```text
~/.config/anifeed/
  .cookies
  rss.json
  node-site-config.json
```

### 数据库文件

`db.sqlite` 默认优先读取当前工作目录的既有文件；如果当前目录不存在，会继续查找 `~/.config/anifeed/db.sqlite`。两处都不存在时，会在当前工作目录创建新的 `db.sqlite`。

如果使用 TOML 配置，可以通过 `[database].path` 指定数据库文件路径。相对路径按 `config.toml` 所在目录解析。

## 配置迁移指南

从传统配置文件迁移到 `config.toml` 可以逐步进行：

### 使用 Linux 安装脚本时

如果通过 [`install-release.sh`](./linux-systemd-install.md) 安装，脚本会自动创建推荐的 systemd 运行目录和基础配置：

```text
/var/lib/anifeed/config.toml
/var/lib/anifeed/.cookies
/var/lib/anifeed/db.sqlite
```

服务的 `WorkingDirectory` 是 `/var/lib/anifeed`，所以这里的 `config.toml` 会被优先读取。脚本生成的基础配置如下：

```toml
[auth]
cookies_file = ".cookies"

[server]
port = 8115

[database]
path = "db.sqlite"
```

从旧配置迁移到脚本安装目录时，可以先把传统文件放到 `/var/lib/anifeed`，保持兼容读取：

```bash
sudo install -m 600 /path/to/.cookies /var/lib/anifeed/.cookies
sudo install -m 600 /path/to/rss.json /var/lib/anifeed/rss.json
sudo install -m 600 /path/to/node-site-config.json /var/lib/anifeed/node-site-config.json
sudo install -m 600 /path/to/db.sqlite /var/lib/anifeed/db.sqlite
```

`db.sqlite` 只在需要保留历史 RSS 任务记录时迁移。完成文件迁移或修改 `config.toml` 后，重启服务：

```bash
sudo systemctl restart anifeed
```

后续可以继续按下面步骤把 `rss.json` 和 `node-site-config.json` 的内容逐步合并到 `/var/lib/anifeed/config.toml`。

### 步骤 1：创建基础 config.toml

普通命令行或手动部署场景，创建 `~/.config/anifeed/config.toml`，先配置认证部分：

如果使用 Linux 安装脚本，则编辑 `/var/lib/anifeed/config.toml`，不需要再创建 `~/.config/anifeed/config.toml`。

```toml
[auth]
cookies_file = ".cookies"
```

此时 RSS 和站点配置仍从 `rss.json` 和 `node-site-config.json` 读取。

### 步骤 2：迁移 RSS 订阅

将 `rss.json` 中的订阅添加到 `config.toml`：

```toml
[[rss]]
site = "mikanani.me"
name = "订阅名称"
url = "https://..."
```

添加 `[[rss]]` 配置后，`rss.json` 将不再被读取。

### 步骤 3：迁移站点配置

将 `node-site-config.json` 中的站点配置添加到 `config.toml`：

```toml
[sites."share.dmhy.org"]
https_agent = true
```

添加 `[sites]` 配置后，`node-site-config.json` 将不再被读取。

### 步骤 4：清理（可选）

完成迁移后，可以删除不再使用的传统配置文件：

- `rss.json`
- `node-site-config.json`

保留 `.cookies` 文件，或将其内容直接配置到 `config.toml` 的 `[auth].cookies` 字段（不推荐）。

## 命令行参数

所有命令行参数的优先级最高，会覆盖配置文件中的设置：

- `--cookies`: 指定 cookies 字符串
- `--rss` / `-r`: 指定 RSS 配置文件路径
- `--port`: 指定服务器端口
- `--disable-cache`: 禁用 115 缓存
- `--chunk-delay`: 设置分块延迟
- `--chunk-size`: 设置分块大小
- `--cooldown-min-ms`: 设置冷却最小时间
- `--cooldown-max-ms`: 设置冷却最大时间

示例：

```bash
# 使用命令行参数覆盖配置文件
anifeed --port 9000 --cookies "UID=xxx; CID=xxx; SEID=xxx; KID=xxx"
```

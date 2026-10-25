# AniFeed

自托管动漫 RSS → 115 网盘离线 Web UI。

## 快速部署

```bash
docker pull huohen92/anifeed:0.02

docker run -d \
  --name anifeed \
  --restart unless-stopped \
  -p 8115:8115 \
  -v "$PWD/data:/data" \
  -e TZ=Asia/Shanghai \
  huohen92/anifeed:0.02
```

打开 `http://服务器IP:8115`，首次创建 Web UI 登录账号，然后在「设置 → 115 账号」填写 Cookies 或扫码登录。

## Docker Compose

```yaml
services:
  anifeed:
    image: huohen92/anifeed:0.02
    container_name: anifeed
    restart: unless-stopped
    ports:
      - "8115:8115"
    volumes:
      - ./data:/data
    environment:
      - TZ=Asia/Shanghai
```

```bash
mkdir -p data
docker compose pull
docker compose up -d
```

## 数据

```text
data/config.toml
data/.cookies
data/db.sqlite
```

## 说明

- 当前镜像：`linux/amd64`
- 基础镜像：`scratch`
- 静态 musl 二进制，约 8 MB
- 源码：https://github.com/huohen92/anifeed
- fnOS FPK：见 GitHub Releases

> 首次初始化前不要把未设置登录的 8115 端口直接暴露到公网。

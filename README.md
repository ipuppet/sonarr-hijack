# sonarr-hijack

## 配置文件

配置文件路径需为 `config/config.json`

示例：

```json
{
    "jackett_scheme": "https",
    "jackett_host": "yourjackett.com",
    "jackett_ip": "127.0.0.1",
    "jackett_port": "443",
    "tmdb_api_key": "your tmdb api key",
    "tmdb_search_url": "https://api.themoviedb.org/3/search/tv"
}
```

- `jackett_host` 可为 IP 地址
- `jackett_ip` 为可选项，项目搭建在本地时，用于跳过域名解析步骤

## 编译运行

项目运行需指定 log 参数

```shell
go mod tidy
make
chmod +x ./sonarr-hijack
mkdir var storage
./sonarr-hijack -log=./var
```

对于 Windows 系统，使用

`make build_win`

## 仅运行

`make run`

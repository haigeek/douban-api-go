# douban-api-go

将 [douban-api-rs](https://github.com/cxfksword/douban-api-rs) 迁移为 Go 版本（Gin），包含电影与图书接口。

## 运行

```bash
go run ./cmd/server
```

可选参数：

- `--host` 监听地址，默认 `0.0.0.0`
- `--port` 监听端口，默认 `8080`
- `--limit` Jellyfin 请求默认搜索条数，默认 `3`，可用 `DOUBAN_API_LIMIT_SIZE` 覆盖
- `--cookie` 豆瓣登录 cookie，可用 `DOUBAN_COOKIE` 覆盖
- `--debug` 开启 debug 日志
- `--basic-user` Basic Auth 用户名（与 `--basic-pass` 同时设置时生效）
- `--basic-pass` Basic Auth 密码（与 `--basic-user` 同时设置时生效）

## Docker

```bash
docker build -t douban-api-go .
docker run -d --name douban-api-go --restart=unless-stopped -p 5000:80 douban-api-go
```

## API

```text
/movies?q={movie_name}                  # 搜索电影
/movies?q={movie_name}&type=full        # 搜索电影并获取详细信息（仅 type=full）
/movies/{sid}                           # 获取指定电影信息
/movies/{sid}/celebrities               # 获取演员列表
/celebrities/{cid}                      # 获取演员信息
/photo/{sid}                            # 获取电影壁纸
/proxy?url={image_url}                  # 图片代理

/v2/book/search?q={book_name}&count=2   # 搜索书籍，count 默认 2，最大 20
/v2/book/id/{sid}                       # 获取指定 id 的书籍
/v2/book/isbn/{isbn}                    # 获取指定 isbn 的书籍
```

### movies 接口 type 参数说明

- `type=full`：返回电影详情列表（会进一步抓取每个结果的详情）
- 不传 `type`：返回基础搜索结果列表
- 其他值：按基础搜索结果列表处理（当前仅 `full` 有特殊行为）

## 返回结果示例

搜索：

```
[
    {
        "cat": "电影",
        "sid": "26862259",
        "name": "乘风破浪 ",
        "rating": "6.8",
        "img": "https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2408407697.jpg",
        "year": " 2017"
    },
    {
        "cat": "电影",
        "sid": "34894589",
        "name": "乘风破浪的姐姐 第一季 ",
        "rating": "6.8",
        "img": "https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2608297477.jpg",
        "year": "2020"
    }
]
```


获取电影信息：

```
{
    "sid": "26862259",
    "name": "乘风破浪",
    "rating": "6.8",
    "img": "https://img1.doubanio.com/view/photo/s_ratio_poster/public/p2408407697.jpg",
    "year": "2017",
    "intro": "赛车手阿浪（邓超 饰）一直对父亲（彭于晏 饰）反对自己的赛车事业耿耿于怀，在向父亲证明自己的过程中，阿浪却意外卷入了一场奇妙的冒险。他在这段经历中结识了一群兄弟好友，一同闯过许多奇幻的经历，也对自己的身世有了更多的了解。",
    "director": "导演",
    "writer": "编剧",
    "actor": "主演",
    "genre": "类型",
    "site": "",
    "country": "制片国家/地区",
    "language": "语言",
    "screen": "上映日期",
    "duration": "片长",
    "episodes": "集数",
    "subname": "上映日期",
    "imdb": "IMDb",
    "celebrities": [
        {
            "id": "1275307",
            "img": "https://img3.doubanio.com/view/celebrity/raw/public/p42220.jpg",
            "name": "韩寒",
            "role": "导演"
        }
    ]
}
```

获取演员信息：

```
{
    "id": "1274235",
    "img": "https://img2.doubanio.com/icon/u183170142-13.jpg",
    "name": "邓超 Chao Deng",
    "role": "演员 / 导演 / 配音 / 主持人",
    "intro": "1979年，邓超出生在一个重新组合的小康家庭，爸爸是博物...",
    "gender": "男",
    "constellation": "水瓶座",
    "birthdate": "1979年02月08日",
    "birthplace": "中国,江西,南昌",
    "nickname": "",
    "imdb": "nm2874732",
    "family": "孙俪(妻)"
}
```

## 腾讯云 `ddns` 插件
> 起因：想用家用宽带 + `ipv6` 做个文件分享服务器

参考 `config.json` 中修改相关内容
```json
[{
  "type": "TENCENT",
  "key": "AKID86qMe*****7NzMTQ8bpoY2JL9j4B",
  "secret": "2rRnun*****g3C7oLejYCajcGgrmR9QS",
  "domains": [{
    "domain": "****.com",
    "sub": "liwd",
    "type": "AAAA"
  }]
}]
```
# 鉴权网关

- [x] 根据登录接口返回UIN判断是否登录/退出成功
- [x] 自动添加跨域处理
- [x] 会话重启即失效
- [x] TLS鉴权
- [x] 负载均衡
    - [x] 随机

## Config

```yaml
# 开启验证
enable_verify: false
# 验证证书
ca_path: "./conf/ca.pem"
module_cert_pem: "./conf/server.pem"
module_key_pem: "./conf/server.key"
# 请求将会通过UIN传输到后端
uin_header_name: "uin"
# 登录页面
sign_page: "https://account.dxkite.cn/signin"
# 跨域配置
cors_config:
  allow_origin:
    - https://account.dxkite.cn
  allow_method:
    - GET
    - POST
  allow_header:
    - Content-Type
  allow_credentials: true
sign_info:
  redirect_name: "redirect_uri"
  redirect_url: "https://dxkite.cn"
# 会话数据
session:
  name: "session"
  domain: "dxkite.cn"
  expires_in: 86400
  secure: true
  http_only: true
  path: "/"
# 路由配置
routes:
  - pattern: "/user/signin"
    signin: true # 登录接口
    backend:
      - http://127.0.0.1:2334?trim_prefix=/user
  - pattern: "/user/signout"
    signout: true # 登出接口
    backend:
      - http://127.0.0.1:2334?trim_prefix=/user
  - pattern: "/user/captcha"
    sign: false #不需要登录
    backend:
      - http://127.0.0.1:2334?trim_prefix=/user
  - pattern: "/user/verify_captcha"
    sign: false #不需要登录
    backend:
      - http://127.0.0.1:2334?trim_prefix=/user
  - pattern: "/user"
    sign: true #需要登录才能访问
    backend:
      - http://127.0.0.1:2334?trim_prefix=/user
```

## Nginx 开启双向认证

将nginx作为一个模块服务，可用gateway连接作为一个模块

```
# 服务器证书
ssl_certificate /cert/fullchain.pem;
ssl_certificate_key /cert/privkey.pem;
# Client验证CA
ssl_client_certificate /cert/ca.pem;
# 验证Client
ssl_verify_client on;
```

## Nginx 反向代理

走http

```
location / {
  proxy_pass http://127.0.0.1:2333;
}
```

走 tls 带证书验证模块（需要Nginx提供模块证书，只允许加载了证书的nginx访问gateway）

```
location / {
  proxy_pass https://127.0.0.1:2333;
  proxy_ssl_name gateway;
  proxy_ssl_certificate /gateway/nginx.pem;
  proxy_ssl_certificate_key /gateway/nginx.key;
  proxy_ssl_trusted_certificate /gateway/conf/ca.pem;
  proxy_ssl_server_name on;
  proxy_ssl_verify on;
}
```

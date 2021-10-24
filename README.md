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
# 会话Cookie名称
cookie_name: "session"
# 登录SESSION位置 (leveldb) 废弃
session_path: "./conf/session"
# 会话过期时间（秒）
session_expire_in: 3600
# 配置热加载（秒）
hot_load: 60
# 登录页面
sign_page: "https://dxkite.cn"
# 跨域配置
cors_config:
  allow_origin:
    - "http://127.0.0.1:2333"
  allow_method:
    - GET
    - POST
sign_info:
  redirect_name: "redirect_uri"
  redirect_url: "https://dxkite.cn"
# 路由配置
routes:
  - pattern: "/signin.php"
    signin: true # 登录接口
    backend:
      - http://127.0.0.1:8088
  - pattern: "/signout.php"
    sign: true
    signout: true # 登出接口
    backend:
      - http://127.0.0.1:8088
  - pattern: "/user"
    sign: true #需要登录才能访问
    backend:
      - https://114.132.243.178?server_name=nginx
  - pattern: "/" # 普通非鉴权接口
    backend:
      - https://114.132.243.178?server_name=nginx
```

## Nginx 开启双向认证

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
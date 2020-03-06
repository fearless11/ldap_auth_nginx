
### beego结合nginx实现LDAP认证登录

   框架： beego 参考： [nginxinc](https://github.com/nginxinc/ldap_auth_nginx/blob/master/backend-sample-app.py)

#### 基本思路

- 存储：cookie存储base64后的用户密码、原始url
- 生产者： 登录页面存储cookie，编码用户密码，返回200跳转到原始url
- 消费者： 认证页面读取cookie，解码用户密码LDAP认证，成功返回200，否则返回401

#### 登录模块 [backend-sample-app.py](https://github.com/nginxinc/ldap_auth_nginx/blob/master/backend-sample-app.py)

- 表单认证：获取nginx传递X-Target，获取到则表单渲染，否则直接500返回内部错误

- 表单渲染：传递X-Target到表单，渲染出包含用户密码和隐藏X-Target的表单登录页

- 表单提交： 获取表单页提交的用户、密码、隐藏X-Target；判断不为空，则base64编码用户密码到cookie，返回200跳转到X-Target；判断为空，跳到表单渲染

#### 认证逻辑  [ldap_auth_nginx-daemon.py](https://github.com/nginxinc/ldap_auth_nginx/blob/master/ldap_auth_nginx-daemon.py)
 
- cookie不存在： 直接返回401时设置`Cache-Control:no-cache`
- cookie存在：用base64解码用户密码，ldap认证，成功返回200，失败返回401时设置`Cache-Control:no-cache`   

#### nginx逻辑 [ldap_auth_nginx.conf](https://github.com/nginxinc/ldap_auth_nginx/blob/master/backend-sample-app.py)

- auth_request模块：首先内部认证

- error_page指令： 内部跳转到登录页

- 缓存模块：认证通过后进行缓存，否则一直需要进行LDAP认证会很慢。

  [proxy-cache实现回源服务器缓存](https://blog.csdn.net/dengjiexian123/article/details/53386586)

  ```sh
  # 编译 auth_request 和 ngx_cache_purge模块
  wget http://labs.frickle.com/files/ngx_cache_purge-2.3.tar.gz
  ./configure --prefix=/usr/local/openresty --with-http_auth_request_module  --add-module=modules/ngx_cache_purge-2.3
  make && make install
  ```

  ```sh
  # 配置
  http {
      # 缓存设置
      # 缓存路径 cache/
      # keys_zone 设置缓存名字和共享内存大小
      # levels 设置缓存文件目录层次；levels=1:2 表示两级目录
      # inactive  在指定时间内没人访问则被删除
      # max_size 最大缓存空间，如果缓存空间满，默认覆盖掉缓存时间最长的资源
  	proxy_cache_path cache/  keys_zone=auth_cache:10m levels=1:2 inactive=7d max_size=1000g;
  	
  	server {
  		listen       81;
  		server_name  192.168.56.101;
  		
  		location / {
                 # 先调整到内部认证/auth,不改变url
  			   auth_request /auth;
  			   # 内部认证返回401，则再次内部跳转到登录页
                 error_page 401 =200 /login;
                 proxy_pass  http://10.51.1.31:5601/;
          }
          
          location /auth {
                # 设置为内部调用，不会改变原始请求url
                internal;
                proxy_pass http://127.0.0.1:8081/auth;
                	
                proxy_pass_request_body off;
             	  proxy_set_header Content-Length "";
             	  
                ### 十分钟内不用重新向后端认证
                ### 当后端设置Cache-Control:no-cache时，nginx不会缓存
             	# 使用auth_cache对应缓存配置
                proxy_cache auth_cache;
                # 对httpcode为200…的缓存10分钟
                proxy_cache_valid 200 10m;
                # 缓存唯一key来进行hash存取，这里的cookie中nginxauth保存base64的用户密码字段
  			   proxy_cache_key "$http_authorization$cookie_nginxauth";
          }
          
         location /login {
  	       # 因为都是内部认证,所以url为原始rul
               proxy_set_header X-Target $request_uri;
               proxy_pass http://127.0.0.1:8081/login;
        }

        # 清理缓存 curl http://192.168.56.101:81/cleancache/auth_cache
        location /cleancache/ {
               # allow 127.0.0.1;
               # deny all;
               proxy_cache_purge auth_cache "$http_authorization$cookie_nginxauth";
        }
      }
  }
  ```

  
#### beego开发  [ldap_auth_nginx](https://gitee.com/fearless11/project/tree/master/ldap_auth_nginx)

bee 工具
```
bee new ldap_auth_nginx
bee run
bee pack 
```

`app.conf`

```ini
; LDAP
addr = "xxx:389"
binddn = "xxx"
bindpass = "xxx"
basedn = "ou=xxx,dc=aaa,dc=com"
tls = false
starttls = false

; 自定义用户白名单
whitelist = "aa;bb"
```

`github.com/astaxie/beego/config.go`

```go
// 增加配置
type Config struct {
    ... 
	LDAPConfig          LDAPConfig
	WhiteList           string
	WhiteMap            map[string]bool
}

// 初始化配置
func assignConfig(ac config.Configer) error {
	for _, i := range []interface{}{BConfig, &BConfig.Listen, &BConfig.WebConfig, &BConfig.Log, &BConfig.LDAPConfig, &BConfig.WebConfig.Session} {
		assignSingleConfig(i, ac)
	}
	...
	// set whitelist
	BConfig.WhiteMap = make(map[string]bool)
	if BConfig.WhiteList != "" {
		list := strings.Split(BConfig.WhiteList, ";")
		for _, name := range list {
			BConfig.WhiteMap[name] = true
		}
	}
    ...
}
```

#### 有个问题
 
- 在表单页要点击两下才能成功跳转，我不清楚原因，如果你知道原因，告诉我下，十分感谢 :）
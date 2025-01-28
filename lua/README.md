# Discord invites

install dependencies

```bash
opm get ledgetech/lua-resty-http
```

setup script

```bash
mkdir -p /usr/local/openresty/nginx/conf/lua.d
```

place `discord_invite.lua` to `/usr/local/openresty/nginx/conf/lua.d/`

In `nginx.conf` add

```conf
http {
  lua_package_path '/usr/local/openresty/nginx/conf/lua.d/?.lua;;';
}
```

place `discord_invite.conf` in `/usr/local/openresty/nginx/conf/conf.d/`
and change `server_name` and other settings like tls

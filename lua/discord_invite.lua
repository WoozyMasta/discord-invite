local http = require "resty.http"
local cjson = require "cjson"

local _M = {}

local function generate_discord_invite(channel_id, discord_token)
    local httpc = http.new()
    local res, err = httpc:request_uri("https://discord.com/api/v9/channels/" .. channel_id .. "/invites", {
        method = "POST",
        headers = {
            ["Authorization"] = "Bot " .. discord_token,
            ["Content-Type"] = "application/json"
        },
        body = cjson.encode({
            max_age = 3600,  -- Invite life time (in seconds)
            max_uses = 1,    -- Max times invite usage
            unique = true    -- Unique invite (one-time)
        })
    })

    if not res then
        ngx.log(ngx.ERR, "failed to request Discord API: ", err)
        return nil, err
    end

    if res.status ~= 200 and res.status ~= 201 then
        ngx.log(ngx.ERR, "failed to create Discord invite: ", res.status, " ", res.reason)
        return nil, "failed to create Discord invite"
    end

    local invite_data = cjson.decode(res.body)
    return invite_data.code
end

function _M.redirect_to_discord_invite(channel_id, discord_token)
    local invite_code, err = generate_discord_invite(channel_id, discord_token)
    if err then
        ngx.status = 500
        ngx.say("Failed to generate Discord invite: ", err)
        ngx.exit(ngx.HTTP_INTERNAL_SERVER_ERROR)
    end

    ngx.redirect("https://discord.gg/" .. invite_code, ngx.HTTP_MOVED_TEMPORARILY)
end

return _M

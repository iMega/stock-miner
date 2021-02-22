## Install webhook on VPS

```shell
wget https://github.com/adnanh/webhook/releases/download/2.8.0/webhook-linux-amd64.tar.gz
tar -zxf webhook-linux-amd64.tar.gz
mv webhook-linux-amd64/webhook /usr/bin/webhook
chmod +x /usr/bin/webhook
```

Copy webhook.service from repo to /usr/lib/systemd/system/webhook.service

Insert secret (see 'Create release hook') and copy hooks.json from repo to /etc/webhook/hooks.json

Copy webhook from repo to /etc/sysconfig/webhook

```shell
systemctl start webhook
tail -f /var/log/messages
```

## Add location in nginx config of VPS

```nginx
server {
    listen 80;
    server_name d.imega.ru;

    location /hooks/ {
        proxy_pass http://172.17.0.1:9000/hooks/;
    }
}
```

## Create release hook

Go to Repo - Settings - Hooks

-   Payload URL: https://d.imega.ru/hooks/stock-miner
-   Content type: application/json
-   Secret: Generate very strong password
-   Select event: Releases

## Generate github token

Go to https://github.com/settings/tokens. Need access to repo.

## Shell script for update app

copy redeploy-stock-miner.sh from Repo to ~/ on VPS

insert github token (see prev. paragraph)

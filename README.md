# Stock miner

[![codecov](https://codecov.io/gh/iMega/stock-miner/branch/master/graph/badge.svg?token=JFHLSRY9NS)](https://codecov.io/gh/iMega/stock-miner)

## SDK

https://tinkoffcreditsystems.github.io/invest-openapi/

https://github.com/TinkoffCreditSystems/invest-openapi-go-sdk

https://tinkoffcreditsystems.github.io/invest-openapi/swagger-ui/

curl -s https://query1.finance.yahoo.com/v10/finance/quoteSummary/AAPL?modules=price | json_pp | grep -C2 regularMarketPrice

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

# Cases

## добавление бумаги в белый список

-   указать ее идентификатор (с помощью sugget)
-   запрос текущей стоимости и вывод ее с возможностью исправить
-   добавить к белому списку бумаг с указанной ценой, являюшейся отправной
    точкой для расчета проведения сделки. Если стоимость на торгах будет ниже указанной,
    брокер приобретет эту акцию.

## Покупка бумаги

-   Если стоимость бумаги похожа на лестницу вниз, брокер смотрит опцию
    "Количество ступеней по лестнице вниз", указанное число в опции разрешает
    брокеру совершить равное количество раз (ступеней) покупки при уменьшении цены.

## Опции

-   Количество ступеней по лестнице вниз
-   Количество ступеней по лестнице вверх
-   Процент со сделки
-   НДФЛ
-   минимальная прибыль

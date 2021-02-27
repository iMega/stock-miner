#!/bin/bash
set -o allexport; source .env; set +o allexport

app_url=$(curl -s -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3.raw"  \
    https://api.github.com/repos/imega/stock-miner/releases/tags/$@ | \
    jq -r '.assets[] | select(.name|test("linux-amd64.tar.gz$")) | .url')

md5_url=$(curl -s -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3.raw"  \
    https://api.github.com/repos/imega/stock-miner/releases/tags/$@ | \
    jq -r '.assets[] | select(.name|test("linux-amd64.tar.gz.md5$")) | .url')

curl -H "Authorization: token $GITHUB_TOKEN" \
     -H "Accept: application/octet-stream" \
     -sL -o stock-miner.tar.gz \
     $app_url && echo done! || (echo failed to download stock-miner && exit 1)

curl -H "Authorization: token $GITHUB_TOKEN" \
     -H "Accept: application/octet-stream" \
     -sL -o stock-miner.tar.gz.md5 \
     $md5_url && echo done! || (echo failed to download stock-miner && exit 1)

echo $(cat stock-miner.tar.gz.md5) stock-miner.tar.gz | md5sum --quiet -c || (echo failed to check md5 sum && exit 1)

tar --overwrite -xvf stock-miner.tar.gz

rm stock-miner.tar.gz*

systemctl restart stock-miner

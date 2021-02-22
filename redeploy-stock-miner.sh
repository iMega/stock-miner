#!/bin/bash

curl -H "Authorization: token <INSERT-GITHUB-TOKEN>" \
     -H "Accept: application/octet-stream" \
     -sL -o app \
     $@ && echo done! || (echo failed to download stock-miner && exit 1)

chmod +x app

systemctl restart stock-miner

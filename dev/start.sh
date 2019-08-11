#!/usr/bin/env bash

/killbyname.sh gaiad

(cd .. && make install)

rm -rf cache

gaiad testnet --v 1 --output-dir cache --chain-id testchain --starting-ip-address 127.0.0.1<<EOF
12345678
EOF

sleep 1
gaiad start --home cache/node0/gaiad
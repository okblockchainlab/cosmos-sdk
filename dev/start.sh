#!/usr/bin/env bash

/killbyname.sh gaiad


rm -rf catch

gaiad testnet --v 1 --output-dir catch --chain-id testchain --starting-ip-address 127.0.0.1<<EOF
12345678
EOF

sleep 1
gaiad start --home catch/node0/gaiad
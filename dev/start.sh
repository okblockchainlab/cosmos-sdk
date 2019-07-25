#!/usr/bin/env bash

/killbyname.sh gaiad

(cd .. && make install2)

rm -rf catch

gaiad testnet --v 1 --output-dir cache --chain-id testchain --starting-ip-address 127.0.0.1<<EOF
12345678
EOF

sleep 1
gaiad start --home catch/node0/gaiad
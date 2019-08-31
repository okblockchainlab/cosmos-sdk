#!/usr/bin/env bash

gaiacli keys add $1<<EOF
12345678
EOF

gaiacli keys add --recover org <<EOF
12345678
any shell idle east pill fine ocean explain derive miracle helmet exit regret dune mountain deal strike gap drive tray super razor broken guess
EOF

gaiacli tx send $(gaiacli keys show $1 -a) 500stake --from=org --yes --chain-id testchain

sleep 1

gaiacli tx staking create-validator --amount 1000000stake \
    --pubkey $(gaiad tendermint show-validator) \
    --chain-id testchain --commission-rate 0.1 \
    --commission-max-rate 0.5 \
    --commission-max-change-rate 0.001 \
    --min-self-delegation 1 --from $1 --moniker $1<<EOF
y
EOF

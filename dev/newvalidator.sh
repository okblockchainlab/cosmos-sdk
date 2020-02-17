#!/usr/bin/env bash

#./recover.sh

index=$1

gaiacli keys add node${index} --home cache/node${index}/gaiacli<<EOF
12345678
EOF

toAddr=$(gaiacli keys show node${index} -a --home cache/node${index}/gaiacli)
valPubkey=$(gaiad tendermint show-validator --home cache/node${index}/gaiacli)


gaiacli tx send $toAddr 10000000000stake --from=307 --fees 2stake -yes -b block <<EOF
12345678
EOF

sleep 1

gaiacli tx staking create-validator --amount 2000000000stake \
    --pubkey ${valPubkey} \
    --chain-id testchain --commission-rate 0.1 \
    --commission-max-rate 0.5 \
    --commission-max-change-rate 0.001 \
    --min-self-delegation 1 --from node${index} --moniker node${index} -y \
    --home cache/node${index}/gaiacli -b block <<EOF
12345678
EOF

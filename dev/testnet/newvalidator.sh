#!/usr/bin/env bash

./recover.sh

index=$1

gaiacli keys add node${index} -y --home cache/node${index}/gaiacli
toAddr=$(gaiacli keys show node${index} -a --home cache/node${index}/gaiacli)
valPubkey=$(gaiad tendermint show-validator --home cache/node${index}/gaiacli)


gaiacli tx send $toAddr 10000000000stake --from=u3 --yes --chain-id testchain --node localhost:10057

sleep 10

gaiacli tx staking create-validator --amount 2000000000stake \
    --pubkey ${valPubkey} \
    --chain-id testchain --commission-rate 0.1 \
    --commission-max-rate 0.5 \
    --commission-max-change-rate 0.001 \
    --min-self-delegation 1 --from node${index} --moniker node${index} -y --node localhost:10057 \
    --home cache/node${index}/gaiacli

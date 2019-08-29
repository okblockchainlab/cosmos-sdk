#!/usr/bin/env bash

index=$1
HOME_PATH=cache/node${index}/gaiacli
ADDR=$(gaiacli keys show node${index} --bech val -a --home ${HOME_PATH})
gaiacli tx staking delegate ${ADDR} 500000000stake --from node${index} --home ${HOME_PATH} -y --chain-id testchain --node tcp://localhost:10057




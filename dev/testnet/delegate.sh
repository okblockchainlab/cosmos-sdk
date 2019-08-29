#!/usr/bin/env bash

if [ $# -lt 2 ]; then
    echo "missing params..."
    exit
fi

'''
for i in $@; do
    HOME_PATH=${GOPATH}/src/github.com/cmdpos/cmapp/dev/testnet/cache/node${i}/gaiacli
    ADDR=$(gaiacli keys show node${i} --bech val -a --home ${HOME_PATH})
    echo "gaiacli tx staking delegate ${ADDR} 1stake --from node${i} --home ${HOME_PATH} --chain-id testchain --node tcp://localhost:10${i}57"
    gaiacli tx staking delegate ${ADDR} 1stake --from node${i} --home ${HOME_PATH} --chain-id testchain --node tcp://localhost:10${i}57<<EOF
y
EOF
done
'''
NODE_ENV="--node tcp://localhost:10057"
HOME_PATH=${GOPATH}/src/github.com/cmdpos/cmapp/dev/testnet/cache/node${1}/gaiacli
ADDR=$(gaiacli keys show node${1} --bech val -a --home ${HOME_PATH})
echo "gaiacli tx staking delegate ${ADDR} ${2}stake --from node${1} --home ${HOME_PATH} --chain-id testchain ${NODE_ENV}"
gaiacli tx staking delegate ${ADDR} ${2}stake --from node${1} --home ${HOME_PATH} --chain-id testchain ${NODE_ENV} <<EOF
y
EOF


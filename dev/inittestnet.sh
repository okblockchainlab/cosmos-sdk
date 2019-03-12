#!/usr/bin/env bash

BASEPORT=20056

COSMOS_TOP=${GOPATH}/src/github.com/cosmos/cosmos-sdk
COSMOS_BIN=${GOPATH}/src/github.com/cosmos/cosmos-sdk/build/darwin

function init {
    cd ${COSMOS_TOP}
    make mac
    cd ${COSMOS_TOP}/dev
    /killbyname.sh gaiad
    rm -rf gentxs
    rm -rf node* *.log *.json
    ${COSMOS_BIN}/gaiad testnet --v $1 -o . --starting-ip-address 127.0.0.1 --base-port ${BASEPORT}  <<EOF
EOF
}

function start {
    for ((index=0; index<${1}; index++)) do
        let p2pport=${BASEPORT}+${index}*100
        let rpcport=${BASEPORT}+${index}*100+1

        echo "./gaiad --home ./node${index}/gaiad start --p2p.laddr tcp://0.0.0.0:${p2pport} --rpc.laddr tcp://0.0.0.0:${rpcport}"
        ${COSMOS_BIN}/gaiad --home ./node${index}/gaiad  start --p2p.laddr tcp://0.0.0.0:${p2pport} --rpc.laddr tcp://0.0.0.0:${rpcport} > g${index}.log 2>&1 &
    done
    echo "start node done"
}

init $1
start $1


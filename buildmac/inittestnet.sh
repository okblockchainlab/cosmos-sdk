#!/usr/bin/env bash

BASEPORT=20056

function init {
    cd ..
    make macgaiad
    cd buildmac
    /killbyname.sh gaiad
    rm -rf gentxs
    rm -rf node* *.log *.json
./gaiad testnet --v $1 -o . --starting-ip-address 127.0.0.1 --base-port ${BASEPORT}  <<EOF
EOF
}

function start {
    for ((index=0; index<${1}; index++)) do
        let p2pport=${BASEPORT}+${index}*100
        let rpcport=${BASEPORT}+${index}*100+1

        sleep 1
        echo "./gaiad --home ./node${index}/gaiad start --p2p.laddr tcp://0.0.0.0:${p2pport} --rpc.laddr tcp://0.0.0.0:${rpcport}"
        ./gaiad --home ./node${index}/gaiad  start --p2p.laddr tcp://0.0.0.0:${p2pport} --rpc.laddr tcp://0.0.0.0:${rpcport} > g${index}.log 2>&1 &
    done
    echo "start node done"
}

init $1
start $1


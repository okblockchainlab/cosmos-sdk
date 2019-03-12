#!/usr/bin/env bash

BASEPORT=20056



function start {
    for ((index=0; index<${1}; index++)) do
        let p2pport=${BASEPORT}+${index}*100
        let rpcport=${BASEPORT}+${index}*100+1
        echo "./gaiad --home ./node${index}/gaiad start --p2p.laddr tcp://0.0.0.0:${p2pport} --rpc.laddr tcp://0.0.0.0:${rpcport}"
        ./gaiad --home ./node${index}/gaiad  start --p2p.laddr tcp://0.0.0.0:${p2pport} --rpc.laddr tcp://0.0.0.0:${rpcport} > g${index}.log 2>&1 &
    done
    echo "start node done"
}

#start $1

/killbyname.sh gaiad
rm -rf *.log

./gaiad --home ./node0/gaiad  start --p2p.laddr tcp://0.0.0.0:20056 --rpc.laddr tcp://0.0.0.0:20057 > g0.log 2>&1 &
./gaiad --home ./node1/gaiad  start --p2p.laddr tcp://0.0.0.0:20156 --rpc.laddr tcp://0.0.0.0:20157 > g1.log 2>&1 &
./gaiad --home ./node2/gaiad  start --p2p.laddr tcp://0.0.0.0:20256 --rpc.laddr tcp://0.0.0.0:20257 > g2.log 2>&1 &


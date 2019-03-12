#!/usr/bin/env bash


index=$1
./gaiad --home ./node${index}/gaiad  start --p2p.laddr tcp://0.0.0.0:20${index}56 \
    --rpc.laddr tcp://0.0.0.0:20${index}57 > g${index}.log 2>&1 &


#!/bin/bash


CUR_DIR=`pwd`

COSMOS_TOP=${GOPATH}/src/github.com/cmdpos/cmapp
BIN_CLI_NAME=gaiacli
BIN_NAME=gaiad

CHAIN_ENV0="--node tcp://localhost:10057 --home cache/node${1}/gaiacli  --chain-id testchain"
CHAIN_HOME="--home ${CUR_DIR}/cache/node${1}/gaiacli"
CHAIN_NODE="--node tcp://localhost:10057"

printf "node${1}\n"
echo "==================================="

gaiacli query account $(gaiacli keys show node${1} -a ${CHAIN_HOME}) ${CHAIN_ENV0}

echo "==================================="
gaiacli query distr params ${CHAIN_ENV0}
echo "==================================="

gaiacli query distr validator-outstanding-rewards ${CHAIN_NODE} --chain-id testchain
echo "==================================="

gaiacli query distr commission $(gaiacli keys show node${1} --bech val -a ${CHAIN_HOME}) ${CHAIN_ENV0} --chain-id testchain
echo "==================================="

gaiacli query distr rewards $(gaiacli keys show node${1} -a ${CHAIN_HOME}) ${CHAIN_NODE} --chain-id testchain

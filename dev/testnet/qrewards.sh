#!/bin/bash


CUR_DIR=`pwd`

BIN_CLI_NAME=gaiacli


CHAIN_ENV0="--node tcp://localhost:10057 --home cache/node${1}/gaiacli  --chain-id testchain"
CHAIN_HOME="--home ${CUR_DIR}/cache/node${1}/gaiacli"
CHAIN_NODE="--node tcp://localhost:10057"

printf "node${1}\n"
gaiacli query distr rewards $(gaiacli keys show node${1} -a ${CHAIN_HOME}) ${CHAIN_NODE} --chain-id testchain

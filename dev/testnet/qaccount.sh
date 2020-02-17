#!/bin/bash


CUR_DIR=`pwd`

COSMOS_TOP=${GOPATH}/src/github.com/cmdpos/cmapp

CHAIN_NODE="--node tcp://localhost:10057"

query_account() {
    for (( i=0; i<${1}; i++)) do
        if [ -d cache/node${i} ]; then
            echo "============================================================="
            echo "node${i} account details:"
            CHAIN_HOME="--home ${CUR_DIR}/cache/node${i}/gaiacli"
            CHAIN_ENV0="--node tcp://localhost:10057 --home "${CUR_DIR}"/cache/node${i}/gaiacli --chain-id testchain"
        gaiacli query account $(gaiacli keys show node${i} -a ${CHAIN_HOME}) ${CHAIN_ENV0}
        fi
    done
}

query_account 10

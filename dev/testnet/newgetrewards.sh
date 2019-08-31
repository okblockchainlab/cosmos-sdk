#!/bin/bash


CUR_DIR=`pwd`

BIN_CLI_NAME=gaiacli


CHAIN_ENV0="--node tcp://localhost:10057 --home cache/node${1}/gaiacli  --chain-id testchain"
CHAIN_HOME="--home ${CUR_DIR}/cache/node${1}/gaiacli"
CHAIN_NODE="--node tcp://localhost:10057"

printf "node${1}\n"

./qrewards.sh ${1}
gaiacli tx distr withdraw-all-rewards --from node${1} ${CHAIN_NODE} ${CHAIN_HOME} --chain-id testchain -y

exit
node3 account details:
Account:
  Address:       cmdpos1rd2vs80sts3ykcrxrgh7tdxej6dal9674cruap
  Pubkey:        cmdpospub1addwnpepq0sejc5vglarj0ynl9vgd78u6d8ew4jjush37l003z78g86m05kwjyc45wp
  Coins:         240000.00000000btc,399720.00000000stake
  AccountNumber: 3
  Sequence:      3

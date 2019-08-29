#!/usr/bin/env bash

./killbyname.sh gaiad

rm -rf ~/.gaiad
rm -rf ~/.gaiacli

./recover.admin.sh

gaiad init --chain-id testchain

gaiacli config chain-id testchain
gaiacli config trust-node true
gaiacli config indent true

gaiad add-genesis-account cosmos1m3gmu4zlnv2hmqfu2jwr97r2653w9yshxkhfea 2000000stake
gaiad gentx --amount 1000000stake \
    --min-self-delegation 1 --commission-rate 0.1 \
    --commission-max-rate 0.5 --commission-max-change-rate 0.001 \
    --pubkey $(gaiad tendermint show-validator) --name admin

rm gentx/data/gentx-*
cp ~/.gaiad/config/gentx/gentx-*.json gentx/data


go run main.go

cp ~/.gaiad/config/genesis.json genesis.tmp.json

gaiad tendermint show-node-id
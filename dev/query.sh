#!/usr/bin/env bash


gaiacli query account cosmos1hg40dv5e237qy28vtyum52ygke32ez35hm307h --chain-id testchain --node localhost:26657
gaiacli query account cosmos1geyy4wtak2q9effnfhze9u4htd8yxxmagdw3q0 --chain-id testchain --node localhost:26657


curl http://localhost:1317/auth/accounts/cosmos1hg40dv5e237qy28vtyum52ygke32ez35hm307h
curl http://localhost:1317/auth/accounts/cosmos1geyy4wtak2q9effnfhze9u4htd8yxxmagdw3q0


gaiacli query staking validators

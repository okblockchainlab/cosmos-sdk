#!/bin/bash


gaiacli query account $(gaiacli keys show node3 --home cache/node3/gaiacli -a) --node localhost:10057 --chain-id testchain


curl http://localhost:1317/auth/accounts/cosmos1hg40dv5e237qy28vtyum52ygke32ez35hm307h
curl http://localhost:1317/auth/accounts/cosmos1geyy4wtak2q9effnfhze9u4htd8yxxmagdw3q0

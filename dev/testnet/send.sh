#!/bin/bash


gaiacli tx send $(gaiacli keys show node2 --home cache/node2/gaiacli -a) 100000000stake \
    --from node3 --node localhost:10057 --chain-id testchain -y --home cache/node3/gaiacli --fees 6000000stake



gaiacli tx send cosmos1hg40dv5e237qy28vtyum52ygke32ez35hm307h 2stake --from q0 --fees 2stake -yes --node localhost:10057



gaiacli tx send cosmos1geyy4wtak2q9effnfhze9u4htd8yxxmagdw3q0 26000000stake --from 307 --fees 6000000stake -yes --node localhost:10057
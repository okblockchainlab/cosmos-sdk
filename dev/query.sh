#!/usr/bin/env bash


gaiacli query account cosmos1hg40dv5e237qy28vtyum52ygke32ez35hm307h --chain-id testchain --node localhost:26657
gaiacli query account cosmos1geyy4wtak2q9effnfhze9u4htd8yxxmagdw3q0 --chain-id testchain --node localhost:26657


curl http://localhost:1317/auth/accounts/cosmos1hg40dv5e237qy28vtyum52ygke32ez35hm307h
curl http://localhost:1317/auth/accounts/cosmos1geyy4wtak2q9effnfhze9u4htd8yxxmagdw3q0


gaiacli query staking validators

# vf
# cosmos1qnkgg9h04v4avc79lzqj9tgdztzlw4e8454mvm
# cosmosvaloper1pjx74f0l6nvwx857e8m5x78fepph4rresakmn3

# my addr
# cosmos1fsvvrkwvlh7084mwlpjek4vjm04enljejpl7z4
gaiacli query distr rewards cosmos1qnkgg9h04v4avc79lzqj9tgdztzlw4e8454mvm cosmosvaloper1pjx74f0l6nvwx857e8m5x78fepph4rresakmn3 --node http://18.163.89.47:20181
gaiacli query distr rewards cosmos1qnkgg9h04v4avc79lzqj9tgdztzlw4e8454mvm  --node http://18.163.89.47:20181
gaiacli query distr rewards cosmos1fsvvrkwvlh7084mwlpjek4vjm04enljejpl7z4 cosmosvaloper1pjx74f0l6nvwx857e8m5x78fepph4rresakmn3 --node http://18.163.89.47:20181
gaiacli query distr rewards cosmos1fsvvrkwvlh7084mwlpjek4vjm04enljejpl7z4  --node http://18.163.89.47:20181
gaiacli query distr commission  cosmosvaloper1pjx74f0l6nvwx857e8m5x78fepph4rresakmn3 --node http://18.163.89.47:20181
gaiacli query distr community-pool   --node http://18.163.89.47:20181

gaiacli query staking delegations-to cosmosvaloper1pjx74f0l6nvwx857e8m5x78fepph4rresakmn3 --node http://18.163.89.47:20181

community-pool

gaiacli query account cosmos1fsvvrkwvlh7084mwlpjek4vjm04enljejpl7z4 --node http://18.163.89.47:20181

#######################################3
gaiacli query staking validators  --node http://18.163.89.47:20181
gaiacli query staking pool  --node http://18.163.89.47:20181
gaiacli query staking params  --node http://18.163.89.47:20181


## cosmos132q0hvhfjx84wl04ez9urnvqs3f7futq48atsw
## cosmosvaloper132q0hvhfjx84wl04ez9urnvqs3f7futqr6la5t
gaiacli query staking delegations cosmos1fsvvrkwvlh7084mwlpjek4vjm04enljejpl7z4  --node http://18.163.89.47:20181


gaiacli query account cosmos132q0hvhfjx84wl04ez9urnvqs3f7futq48atsw --node http://18.163.89.47:20181

gaiacli query staking delegations-to cosmosvaloper132q0hvhfjx84wl04ez9urnvqs3f7futqr6la5t --node http://18.163.89.47:20181


gaiacli query staking delegation cosmos1qpfqusq9atcmag6nmjhh8age3jaznw0nwrjg5j cosmosvaloper132q0hvhfjx84wl04ez9urnvqs3f7futqr6la5t --node http://18.163.89.47:20181



gaiacli query distr rewards cosmos1qpfqusq9atcmag6nmjhh8age3jaznw0nwrjg5j cosmosvaloper132q0hvhfjx84wl04ez9urnvqs3f7futqr6la5t --node http://18.163.89.47:20181












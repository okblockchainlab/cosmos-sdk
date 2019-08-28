#!/usr/bin/env bash


#gaiacli tx staking unbond cosmosvaloper1ru8xdlfwhdskkn4h87v5mxuvaddnse9wpsxp9z 1stake --from 307 --fees 2stake -yes -b block <<EOF

#gaiacli tx staking redelegate [src-validator-addr] [dst-validator-addr] [amount] [flags]

gaiacli tx staking delegate cosmosvaloper1ru8xdlfwhdskkn4h87v5mxuvaddnse9wpsxp9z 10000stake --from 307 --fees 2stake -yes -b block <<EOF
12345678
EOF
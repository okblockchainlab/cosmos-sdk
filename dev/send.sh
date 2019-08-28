#!/usr/bin/env bash

#gaiacli tx send cosmos15fypuqpfvkdgg590u6phdq9ts6ag50dqz38nxf 1stake \
#--from node0 --home catch/node0/gaiacli \
#--memo "hqeiurh qeiu hqeiur hqeiur qieuv gqeiurvqei uvqiuvqeifuvqiuhvqieuvwpiuvhw" \
#--chain-id testchain --fees 2stake -y <<EOF
#12345678
#EOF


gaiacli tx send cosmos1hg40dv5e237qy28vtyum52ygke32ez35hm307h 2stake --from q0 --fees 2stake -yes -b block <<EOF
12345678
EOF

gaiacli tx send cosmos1geyy4wtak2q9effnfhze9u4htd8yxxmagdw3q0 20stake --from 307 --fees 2stake -yes -b block <<EOF
12345678
EOF



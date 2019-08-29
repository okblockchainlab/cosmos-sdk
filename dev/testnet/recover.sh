#!/bin/bash



for ((index=0; index<4; index++)) do
    pkey=$(cat cache/node${index}/gaiacli/key_seed.txt)
    gaiacli keys add --recover u${index}  << EOF
$pkey
EOF

    gaiacli keys add --recover u${index}  << EOF
y
$pkey
EOF
done



gaiacli config chain-id testchain
gaiacli config trust-node true
gaiacli config output json
gaiacli config indent true


#mnemonicCaptain
#depart neither they audit pen exile fire smart tongue express blanket burden culture shove curve address together pottery suggest lady sell clap seek whisper
#prikey
#2c999c5afe7f0c902846e1b286fed29c5c5914998655d469568560955abe0d5d
#address
#cosmos1hg40dv5e237qy28vtyum52ygke32ez35hm307h

gaiacli keys add --recover 307  << EOF
depart neither they audit pen exile fire smart tongue express blanket burden culture shove curve address together pottery suggest lady sell clap seek whisper
EOF

gaiacli keys add --recover q0 << EOF
country live width exotic behind mad belt bachelor later outside forget rude pudding material orbit shoot kind curve endless prosper make exotic welcome maple
EOF

gaiacli keys add --recover 307  << EOF
y
depart neither they audit pen exile fire smart tongue express blanket burden culture shove curve address together pottery suggest lady sell clap seek whisper
EOF


gaiacli keys add --recover q0 << EOF
y
country live width exotic behind mad belt bachelor later outside forget rude pudding material orbit shoot kind curve endless prosper make exotic welcome maple
EOF

#!/bin/bash
CURDIR=`pwd`

cd $GOPATH/src
tar -zcvf ${CURDIR}/src.tar.gz \
 golang.org \
 google.golang.org \
 gopkg.in \
 github.com/gogo/protobuf \
 github.com/go-kit/kit/log \
 github.com/go-kit/kit/metrics \
 github.com/prometheus/client_golang/prometheus \
 github.com/btcsuite/btcutil/bech32 \
 github.com/go-logfmt/logfmt \
 github.com/btcsuite/btcd/btcec \
 github.com/gorilla/websocket \
 github.com/pkg/errors \
 github.com/rs/cors \
 github.com/spf13/afero \
 github.com/spf13/cast \
 github.com/spf13/jwalterweatherman \
 github.com/spf13/pflag \
 github.com/spf13/cobra \
 github.com/spf13/viper \
 github.com/syndtr/goleveldb/leveldb \
 github.com/tendermint/btcd/btcec \
 github.com/tendermint/go-amino \
 github.com/tendermint/iavl \
 github.com/rcrowley/go-metrics \
 github.com/prometheus/client_model/go \
 github.com/prometheus/common/expfmt \
 github.com/prometheus/common/model \
 github.com/prometheus/procfs \
 github.com/pelletier/go-toml \
 github.com/mitchellh/mapstructure \
 github.com/magiconair/properties \
 github.com/hashicorp/hcl \
 github.com/beorn7/perks/quantile \
 github.com/davecgh/go-spew/spew \
 github.com/fsnotify/fsnotify \
 github.com/golang \
 github.com/matttproud/golang_protobuf_extensions/pbutil \
 github.com/prometheus/common/internal/bitbucket.org/ww/goautoneg


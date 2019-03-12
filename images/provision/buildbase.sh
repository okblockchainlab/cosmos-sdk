#!/bin/bash

NAME=okchain/baseimage
RELEASE=`uname -m`-$1
DOCKERHUB_NAME=$NAME:$RELEASE

CURDIR=`dirname $0`


function downloadProtobuf() {
    if [ ! -f v3.0.2.tar.gz ]; then
        echo "download protobuf..."
        wget --quiet https://github.com/google/protobuf/archive/v3.0.2.tar.gz
    fi
}

(cd $CURDIR/../../images/base && downloadProtobuf)

docker inspect $DOCKERHUB_NAME 2>&1 > /dev/null
if [ "$?" == "0" ]; then
    echo "BUILD-CACHE: exists!"
    BASENAME=$DOCKERHUB_NAME
else
    echo "BUILD-CACHE: Pulling \"$DOCKERHUB_NAME\" from dockerhub.."
    docker pull $DOCKERHUB_NAME
    docker inspect $DOCKERHUB_NAME 2>&1 > /dev/null
    if [ "$?" == "0" ]; then
	echo "BUILD-CACHE: Success!"
	BASENAME=$DOCKERHUB_NAME
    else
	echo "BUILD-CACHE: WARNING - Build-cache unavailable, attempting local build"
	(cd $CURDIR/../../images/base && make docker DOCKER_TAG=localbuild)
	if [ "$?" != "0" ]; then
            echo "ERROR: Build-cache could not be compiled locally"
            exit -1
	fi
	BASENAME=$NAME:localbuild
    fi
fi

# Ensure that we have the baseimage we are expecting
docker inspect $BASENAME 2>&1 > /dev/null
if [ "$?" != "0" ]; then
   echo "ERROR: Unable to obtain a baseimage"
   exit -1
fi

# any further errors should be fatal
set -e

TMP=`mktemp -d`
DOCKERFILE=$TMP/Dockerfile

LOCALSCRIPTS=$TMP/scripts
REMOTESCRIPTS=/okchain/scripts/provision

mkdir -p $LOCALSCRIPTS
cp -R $CURDIR/* $LOCALSCRIPTS

echo "BASENAME: $BASENAME"
# extract the FQN environment and run our common.sh to create the :latest tag
cat <<EOF > $DOCKERFILE
FROM $BASENAME
`for i in \`docker run -i $BASENAME /bin/bash -l -c printenv\`;
do
   echo ENV $i
done`
COPY scripts $REMOTESCRIPTS
RUN $REMOTESCRIPTS/common.sh
RUN chmod a+rw -R /opt/gopath

EOF

[ ! -z "$http_proxy" ] && DOCKER_ARGS_PROXY="$DOCKER_ARGS_PROXY --build-arg http_proxy=$http_proxy"
[ ! -z "$https_proxy" ] && DOCKER_ARGS_PROXY="$DOCKER_ARGS_PROXY --build-arg https_proxy=$https_proxy"
[ ! -z "$HTTP_PROXY" ] && DOCKER_ARGS_PROXY="$DOCKER_ARGS_PROXY --build-arg HTTP_PROXY=$HTTP_PROXY"
[ ! -z "$HTTPS_PROXY" ] && DOCKER_ARGS_PROXY="$DOCKER_ARGS_PROXY --build-arg HTTPS_PROXY=$HTTPS_PROXY"
[ ! -z "$no_proxy" ] && DOCKER_ARGS_PROXY="$DOCKER_ARGS_PROXY --build-arg no_proxy=$no_proxy"
[ ! -z "$NO_PROXY" ] && DOCKER_ARGS_PROXY="$DOCKER_ARGS_PROXY --build-arg NO_PROXY=$NO_PROXY"
docker build $DOCKER_ARGS_PROXY -t $NAME:latest $TMP

echo "docker file:"
ls -l $TMP
cat $TMP/Dockerfile
rm -rf $TMP

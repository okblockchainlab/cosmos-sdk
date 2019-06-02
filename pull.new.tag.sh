#!/usr/bin/env bash
git remote add upstream https://github.com/cosmos/cosmos-sdk

git checkout remotes/upstream/$1

git tag $1

git push origin $1
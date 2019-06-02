#!/usr/bin/env bash
git remote add upstream https://github.com/cosmos/cosmos-sdk

git checkout remotes/upstream/$1


git checkout -b $1


git push --set-upstream origin $1

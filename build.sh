#!/bin/bash

path=$(cd `dirname $0`; pwd)

export GOPATH=$path

echo "GOPATH=$GOPATH"

cd $GOPATH/src/macd.me/blogci

echo "开始编译..."

go install

echo "编译成功"
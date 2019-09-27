#!/bin/bash

HOST=${1:-apigwdev01.its.auckland.ac.nz}
SOURCE=${2:-../$(basename $PWD)}
TARGET=${3:-~}

while true; do
  inotifywait -r -e modify,attrib,close_write,move,create,delete "$SOURCE"
  rsync -avz -e "ssh -o StrictHostKeyChecking=no"  "$SOURCE" $HOST:$TARGET
done

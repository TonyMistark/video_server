#! /bin/bash

# Build UI

cp -R ./template ./bin/
mkdir ./bin/videos

cd bin
nohup ./api &
nohup ./schduler &
nohup ./streamserver &
nohup ./web &

echo "deploy finish"
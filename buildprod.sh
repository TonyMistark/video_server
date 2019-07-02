#! /bin/bash

# Build web and other services

cd /Users/ice/projects/go/src/video_server/api
env GOOS=linux GOARCH=amd64 go build -o ../bin/api


cd /Users/ice/projects/go/src/video_server/scheduler
env GOOS=linux GOARCH=amd64 go build -o ../bin/scheduler


cd /Users/ice/projects/go/src/video_server/web
env GOOS=linux GOARCH=amd64 go build -o ../bin/web
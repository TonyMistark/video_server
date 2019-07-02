#! /bin/bash
# Build web UI
cd /Users/ice/projects/go/src/video_server/web
go install
cp /Users/ice/projects/go/bin/web /Users/ice/projects/go/bin/video_server_web_ui/
cp -R /Users/ice/projects/go/src/video_server/templates /Users/ice/projects/go/bin/video_server_web_ui/

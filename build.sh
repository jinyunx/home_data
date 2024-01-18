######## mac ########
pwd
mkdir -p release/mac
go build .
mv home_data release/mac/task
cp main.html release/mac/

mkdir -p release/mac/web/videoviewer
go build ./web/videoviewer/
mv videoviewer release/mac/web/videoviewer/
cp web/videoviewer/*.html release/mac/web/videoviewer/

mkdir -p release/mac/crawl/js
cp crawl/js/*.js release/mac/crawl/js

echo "killall task && nohup ./task >log.txt 2>&1 &" > release/mac/run.sh
echo "killall videoviewer && nohup ./videoviewer >log.txt 2>&1 &" > release/mac/web/videoviewer/run.sh
######## pi ########
mkdir -p release/pi
env GOOS=linux GOARCH=arm GOARM=6 go build .
mv home_data release/pi/task
cp main.html release/pi/

mkdir -p release/pi/web/videoviewer
env GOOS=linux GOARCH=arm GOARM=6 go build ./web/videoviewer/
mv videoviewer release/pi/web/videoviewer
cp web/videoviewer/*.html release/pi/web/videoviewer

mkdir -p release/pi/crawl/js
cp crawl/js/*.js release/pi/crawl/js

echo "sodu killall task && sodu nohup ./task &" > release/pi/run.sh
echo "sodu killall videoviewer && sodu nohup ./videoviewer &" > release/pi/web/videoviewer/run.sh
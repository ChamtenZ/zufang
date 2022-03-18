#! /bin/bash
sudo fdfs_trackerd /etc/fdfs/tracker.conf
sudo fdfs_storaged /etc/fdfs/storage.conf

#sudo fdfs_trackerd /etc/fdfs/tracker.conf stop
#sudo fdfs_storaged /etc/fdfs/storage.conf stop

sudo /usr/local/nginx/sbin/nginx
#sudo /usr/local/nginx/sbin/nginx -s stop

go run ~/work/go/src/zufang/service/captcha/main.go
go run ~/work/go/src/zufang/service/user/main.go
go run ~/work/go/src/zufang/web/main.go
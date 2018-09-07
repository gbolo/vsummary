
docker stop vsum-mysql
docker rm vsum-mysql

docker run -d --name vsum-mysql \
-e MYSQL_ROOT_PASSWORD=secret \
-e MYSQL_DATABASE=vsummary \
-e MYSQL_USER=vsummary \
-e MYSQL_PASSWORD=secret \
-p 3306:3306 \
mysql:5.7

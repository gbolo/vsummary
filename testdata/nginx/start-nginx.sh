docker rm -f vsummary-nginx

docker run --rm -d --name vsummary-nginx \
-v $(pwd)/testdata/nginx/conf.d:/etc/nginx/conf.d:ro \
-p 8081:8081 \
gbolo/nginx:alpine

FROM alpine:edge

MAINTAINER George Bolo <gbolo@linuxctl.com>

EXPOSE 80 443
VOLUME ["/data"]

# -----------------------------------------------------------------------------
# Set some ENV variables
# -----------------------------------------------------------------------------
ENV LANG="en_US.UTF-8" \
    LC_ALL="en_US.UTF-8" \
    LANGUAGE="en_US.UTF-8" \
    MYSQL_ROOT_USER="vsummary" \
    MYSQL_ROOT_PASSWORD="changeme" \
    TERM="xterm"

# -----------------------------------------------------------------------------
# Install required software
# -----------------------------------------------------------------------------
RUN apk add --no-cache bash supervisor mariadb mariadb-client \
    nginx php5-fpm php5-pdo php5-json php5-curl php5-pdo_mysql \
    python python-dev uwsgi uwsgi-python py2-pip \
    && pip install pyvmomi flask pymysql

# -----------------------------------------------------------------------------
# Prepare and configure
# -----------------------------------------------------------------------------
RUN mkdir -p /data/mysql/data && \
    mkdir -p /data/mysql/conf && \
    mkdir -p /data/nginx/www && \
    mkdir -p /data/nginx/conf && \
    mkdir -p /data/php-fpm/conf && \
    mkdir -p /data/flask

COPY ./docker/my.cnf /data/mysql/conf/my.cnf
COPY ./docker/nginx.conf /data/nginx/conf/nginx.conf
COPY ./docker/php-fpm.conf /data/php-fpm/conf/php-fpm.conf
COPY ./docker/supervisord.conf /etc/supervisord.conf

COPY ./sql/vsummary_mysql_schema.sql /data/mysql/conf/vsummary_mysql_schema.sql
COPY ./scripts/data-generator/gen_sample_data.php /data/gen_sample_data.php
COPY ./docker/wrapper.sh /wrapper.sh

COPY ./src /data/nginx/www/

COPY ./collectors/internal/python /data/flask/

RUN chmod +x /wrapper.sh
CMD ["/wrapper.sh"]

# -----------------------------------------------------------------------------
# run supervisord
# -----------------------------------------------------------------------------
#CMD ["/usr/bin/supervisord"]

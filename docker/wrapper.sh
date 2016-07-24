#!/bin/sh
set -xo pipefail

########################################
#  vsummary -- SET UP MySQL            #
########################################
mkdir -p /data/mysql/logs
mkdir -p /var/run/mysqld
chown -R mysql:mysql /data/mysql
chown -R mysql:mysql /var/run/mysqld

DATA_DIR="/data/mysql/data"
PID_FILE=/var/run/mysqld/mysqld.pid

if [ ! -d "$DATA_DIR/mysql" ]; then
  if [ -z "$MYSQL_ROOT_PASSWORD" ]; then
    echo >&2 'error: database is uninitialized and password option is not specified '
    echo >&2 'You need to specify MYSQL_ROOT_PASSWORD'
    exit 1
  fi

  rm -rf /etc/mysql/my.cnf
  ln -s /data/mysql/conf/my.cnf /etc/mysql/my.cnf

  mkdir -p "$DATA_DIR"
  chown -R mysql:mysql "$DATA_DIR"

  echo 'Initializing database'
  mysql_install_db --user=mysql --datadir="$DATA_DIR" --rpm
  echo 'Database initialized'

  mysqld_safe --pid-file=$PID_FILE --skip-networking --nowatch

  mysql_options='--protocol=socket -uroot'

  for i in `seq 30 -1 0`; do
    if mysql $mysql_options -e 'SELECT 1' &> /dev/null; then
      break
    fi
    echo 'MySQL init process in progress...'
    sleep 1
  done
  if [ "$i" = 0 ]; then
    echo >&2 'MySQL init process failed.'
    exit 1
  fi

  mysql $mysql_options <<-EOSQL
    -- What's done in this file shouldn't be replicated
    --  or products like mysql-fabric won't work
    SET @@SESSION.SQL_LOG_BIN=0;

    DELETE FROM mysql.user ;
    CREATE USER 'root'@'%' IDENTIFIED BY '${MYSQL_ROOT_PASSWORD}' ;
    GRANT ALL ON *.* TO 'root'@'%' WITH GRANT OPTION ;
    DROP DATABASE IF EXISTS test ;
    FLUSH PRIVILEGES ;
EOSQL

  mysql_options="$mysql_options -p${MYSQL_ROOT_PASSWORD}"

  # create vsummary database and user
  mysql $mysql_options -e "CREATE DATABASE IF NOT EXISTS vsummary ;"
  mysql $mysql_options -e "CREATE USER 'vsummary'@'%' IDENTIFIED BY 'changeme' ;"
  mysql $mysql_options -e "GRANT ALL ON vsummary.* TO 'vsummary'@'%' ;"
  mysql $mysql_options -e 'FLUSH PRIVILEGES ;'

  # create mysql schema
  mysql $mysql_options vsummary < /data/mysql/conf/vsummary_mysql_schema.sql

  pid="`cat $PID_FILE`"
  if ! kill -s TERM "$pid"; then
    echo >&2 'MySQL init process failed.'
    exit 1
  fi

  # make sure mysql completely ended
  sleep 2

  echo
  echo 'MySQL init process done. Ready for start up.'
  echo
fi

########################################
#  vsummary -- SET UP NGINX + PHP-FPM  #
########################################
mkdir -p /data/nginx/logs
chown -R nginx: /data/nginx
rm -rf /etc/nginx/nginx.conf
ln -s /data/nginx/conf/nginx.conf /etc/nginx/nginx.conf

mkdir -p /data/php-fpm/logs
chown -R nginx: /data/php-fpm
rm -rf /etc/php5/php-fpm.conf
ln -s /data/php-fpm/conf/php-fpm.conf /etc/php5/php-fpm.conf

########################################
#  vsummary -- START SUPERVISORD       #
########################################
/usr/bin/supervisord -c /etc/supervisord.conf

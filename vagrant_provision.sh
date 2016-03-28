#!/bin/bash

curl -s "https://www.postgresql.org/media/keys/ACCC4CF8.asc" | apt-key add -

echo "deb http://apt.postgresql.org/pub/repos/apt/ trusty-pgdg main" > /etc/apt/sources.list.d/postgresql.list

apt-get update && apt-get dist-upgrade -y
apt-get install -y postgresql-9.5

sudo -u postgres psql -c "CREATE USER vagrant WITH PASSWORD 'vagrant'"
sudo -u postgres psql -c "CREATE DATABASE ghmirror"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE ghmirror TO vagrant"

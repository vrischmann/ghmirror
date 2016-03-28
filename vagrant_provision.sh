#!/bin/bash

curl -s "https://www.postgresql.org/media/keys/ACCC4CF8.asc" | apt-key add -

echo "deb http://apt.postgresql.org/pub/repos/apt/ jessie-pgdg main" | sudo tee -a /etc/apt/sources.list.d/postgresql.list

apt-get update
apt-get install -y postgresql-9.5

sudo -u postgres psql -c "CREATE USER vagrant WITH PASSWORD 'vagrant'"
sudo -u postgres psql -c "CREATE DATABASE ghmirror"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE ghmirror TO vagrant"

sudo sed -i "s/#listen_addresses.*/listen_addresses = '*'/g" /etc/postgresql/9.5/main/postgresql.conf
sudo sed -i '$s/$/\n\nhost all all samenet md5/' /etc/postgresql/9.5/main/pg_hba.conf
sudo systemctl restart postgresql

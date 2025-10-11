#!/bin/bash

# PostgreSQL Setup Script for Flow Backend
# Run this script with: sudo ./setup_postgres.sh

echo "Setting up PostgreSQL for Flow Backend..."

# Update package list
apt update

# Install PostgreSQL
apt install -y postgresql postgresql-contrib

# Start PostgreSQL service
service postgresql start

# Enable PostgreSQL to start on boot
systemctl enable postgresql

# Switch to postgres user and configure database
sudo -u postgres psql << EOF
-- Create user with password
CREATE USER postgres WITH PASSWORD 'password';
ALTER USER postgres CREATEDB;

-- Create database
CREATE DATABASE flow OWNER postgres;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE flow TO postgres;

-- Exit
\q
EOF

# Configure PostgreSQL to accept connections
echo "Configuring PostgreSQL..."

# Backup original config
cp /etc/postgresql/*/main/postgresql.conf /etc/postgresql/*/main/postgresql.conf.backup
cp /etc/postgresql/*/main/pg_hba.conf /etc/postgresql/*/main/pg_hba.conf.backup

# Update postgresql.conf to listen on all addresses
sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" /etc/postgresql/*/main/postgresql.conf

# Update pg_hba.conf to allow password authentication
echo "host    all             all             127.0.0.1/32            md5" >> /etc/postgresql/*/main/pg_hba.conf

# Restart PostgreSQL
service postgresql restart

echo "PostgreSQL setup complete!"
echo "Database: flow"
echo "User: postgres"
echo "Password: password"
echo "Host: localhost"
echo "Port: 5432"
echo ""
echo "You can now run your Flow backend application!"

#!/bin/bash

# PostgreSQL Authentication Fix Script
# Run this script with: sudo ./fix_postgres_auth.sh

echo "Fixing PostgreSQL authentication for Flow Backend..."

# Method 1: Reset postgres user password
echo "Resetting postgres user password..."
sudo -u postgres psql << EOF
ALTER USER postgres PASSWORD 'password';
\q
EOF

echo "Password reset complete!"

# Method 2: Create a dedicated flow user (alternative)
echo ""
echo "Alternative: Creating dedicated flow user..."
sudo -u postgres psql << EOF
CREATE USER flowuser WITH PASSWORD 'password';
ALTER USER flowuser CREATEDB;
CREATE DATABASE flow OWNER flowuser;
GRANT ALL PRIVILEGES ON DATABASE flow TO flowuser;
\q
EOF

echo "Dedicated user created!"

# Test connection
echo ""
echo "Testing connection..."
psql -h localhost -U postgres -d flow -c "SELECT 'Connection successful!' as status;"

echo ""
echo "Authentication fix complete!"
echo "You can now run your Flow backend application."

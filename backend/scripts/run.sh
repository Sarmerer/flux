#!/bin/bash

# Flow Backend Run Script
# This script builds and runs the Flow backend application

echo "Building Flow backend..."
go build -o bin/api ./cmd/api

if [ $? -eq 0 ]; then
    echo "Build successful! Starting server..."
    echo "Make sure PostgreSQL is running on localhost:5432"
    echo "Press Ctrl+C to stop the server"
    echo ""
    ./bin/api
else
    echo "Build failed!"
    exit 1
fi

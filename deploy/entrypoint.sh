#!/bin/bash
set -e

echo "Running database migrations..."
migrator up

echo "Starting HTTP server..."
exec http-server
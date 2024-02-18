#!/bin/bash

set -xe

# Set environment variables
export POSTGRES_USER="gophermart"
export POSTGRES_PASSWORD="gophermart"
export POSTGRES_DB="gophermart"
export POSTGRES_HOST="postgres"  # Set your PostgreSQL host here
export POSTGRES_PORT=5432  # Set your PostgreSQL port here
export APP_PORT=8080  # Set the host port for the gophermart service


# Set environment variables
export JWT_SECRET_KEY=`openssl rand -base64 32`

export JAEGER_HOST="jaeger" # Example Jaeger host

# Reuse environment variables for database URI
export DATABASE_URI="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}"

# Run Docker Compose
docker-compose up -d

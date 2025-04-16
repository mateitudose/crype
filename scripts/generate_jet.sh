#!/bin/bash

# Path to the .env file
ENV_FILE="../.env"

# Check if .env file exists
if [ ! -f "$ENV_FILE" ]; then
    echo "Error: .env file not found at $ENV_FILE"
    exit 1
fi

# Load environment variables from .env file
source "$ENV_FILE"

# Build the PostgreSQL connection string
CONNECTION_STRING="postgres://${CRYPE_DB_USER}:${CRYPE_DB_PASSWORD}@${CRYPE_DB_HOST}:${CRYPE_DB_PORT}/${CRYPE_DB_NAME}?sslmode=disable"

echo "Running PostgreSQL with connection string: $CONNECTION_STRING"

# Run Jet with the connection string
# Modify the Jet command as needed
echo "Generating SQL Builder and Models using Jet..."
jet -dsn="$CONNECTION_STRING" -schema=public -path=../api/gen_jet
#!/bin/bash

export PGPASSWORD='postgres'; 

# Regular DB
psql -h postgres -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'pubcast'" | grep -q 1 || psql -h postgres -U postgres -c "CREATE DATABASE pubcast"

# Testing DB
psql -h postgres -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'pubcast_test'" | grep -q 1 || psql -h postgres -U postgres -c "CREATE DATABASE pubcast_test"

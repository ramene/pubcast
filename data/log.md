#!/bin/bash

set -ex

export GOPATH=/go && \
export PATH=$PATH:$GOPATH/bin

# install go 1.11+ via direct binary download
# migrate requires go v1.11 - see: https://github.com/golang-migrate/migrate/tree/master/cli
# Downloads: https://github.com/golang/dep
# https://golang.org/dl/
cd /opt && wget -qO- https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz | tar xvz -C . && cp go/bin/go /bin && export GOROOT=/opt/go

go version

wget -qO- https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz
tar -xvf go1.11.1.linux-amd64.tar.gz

wget -qO- \
    https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz | \
    tar xvz -C .

```sh
$ root@ap-sidecar:/opt/go/bin# go version
go: cannot find GOROOT directory: /usr/local/go
$ root@ap-sidecar:/opt/go/bin# export GOROOT=/opt/go
$ root@ap-sidecar:/opt/go/bin# go version
go version go1.11.1 linux/amd64
```

# Install go/dep
# TODO: Add to Dockerfile - see: https://github.com/golang/dep
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Build Pubcast
# root@ap-sidecar:/go/src/pubcast# dep ensure

# Run Make Test
# root@ap-sidecar:/go/src/pubcast# make test 

# Install migrate
# see: https://github.com/golang-migrate/migrate/tree/master/cli
go get -u github.com/lib/pq && \
go get -u -d github.com/golang-migrate/migrate/cli && \
cd $GOPATH/src/github.com/golang-migrate/migrate/cli && \
go build -tags 'postgres' -o /go/bin/migrate github.com/golang-migrate/migrate/cli


# Create DB
export PGPASSWORD='postgres' && \
psql -h postgres -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'pubcast'" | grep -q 1 || psql -h postgres -U postgres -c "CREATE DATABASE pubcast" && \
psql -h postgres -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'pubcast_test'" | grep -q 1 || psql -h postgres -U postgres -c "CREATE DATABASE pubcast_test"



# Migrate db, FROM our cloned repo
# root@ap-sidecar:/go/src/pubcast# 
cd $GOPATH/src/pubcast && \
migrate -source file://data/migrations -database postgres://postgres:postgres@postgres:5432/pubcast?sslmode=disable up && \
migrate -source file://data/migrations -database postgres://postgres:postgres@postgres:5432/pubcast_test?sslmode=disable up

# Verify
root@ap-sidecar:/go/src/pubcast# psql postgres://postgres:postgres@postgres:5432/pubcast
psql (10.6 (Debian 10.6-1.pgdg90+1), server 9.6.2)
Type "help" for help.

pubcast=# \dt
               List of relations
 Schema |       Name        | Type  |  Owner
--------+-------------------+-------+----------
 public | groups            | table | postgres
 public | organizations     | table | postgres
 public | podcasts          | table | postgres
 public | schema_migrations | table | postgres
(4 rows)

pubcast=# \q
root@ap-sidecar:/go/src/pubcast#
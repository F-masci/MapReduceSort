#!/bin/bash

# Master
go run master.go --idx 0 &
go run master.go --idx 1 &

# Mapper
go run worker.go --address localhost --proto tcp --port 45980 --map &
go run worker.go --address localhost --proto tcp --port 45981 --map &

# Reducer
go run worker.go --address localhost --proto tcp --port 45989 --reduce &
go run worker.go --address localhost --proto tcp --port 45990 --reduce &

# Client
# go run client.go --client fmasci --master-idx 0 &

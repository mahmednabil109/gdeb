#!/bin/bash

# run the bootstrap node
# go run main.go -port 8080 &
# sleep 1


# create the other nodes
port=8080
pk=0

go run ./main.go -port 8080 -first 1 -pk 0 &
sleep 1

for i in `seq 8`; do
     ((port=8080+i))
     ((pk=0+i))
     echo $pk $port
     go run ./main.go -port $port -first 0 -pk $pk -bootstrap 127.0.0.1:8080 &
     sleep 1
done
# sleep 2
# go run ./main.go -port 8091 -first 0 -send 1 -pk 9 -bootstrap 127.0.0.1:8080 &

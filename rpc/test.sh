#!/bin/bash
i=0
while [ $i != 10 ]
do
#go run localhost:1111 /tmp/testyaml/
nohup /home/liuchjlu/golang/go/bin/go  run client.go  localhost:1111 /tmp/testyaml/"$i".yaml &
let i++
sleep 1
done

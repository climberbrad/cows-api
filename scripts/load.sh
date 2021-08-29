#!/bin/sh


CNT=`cat data.json| jq ". | length"`

for ((i=0;i<=CNT-1;i++)); do
   COW=`cat data.json | jq ".[$i]"`
   POST=`curl -XPOST -H"Content-Type:application/json" localhost:8080/cows -d "$COW"`
done


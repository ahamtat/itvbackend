#!/bin/bash

# Valid resource
for i in `seq 1 3`;
do
  echo $i
  curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"method":"GET","url":"http://google.com"}' \
    http://localhost:8080/v1/requests/request

done
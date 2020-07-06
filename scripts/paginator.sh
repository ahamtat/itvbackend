#!/bin/bash

# Populate server storage
for i in `seq 1 10`;
do
  echo $i
  curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"method":"GET","url":"http://google.com"}' \
    http://localhost:8080/v1/requests/request
done

# Get page 2
curl --header "Content-Type: application/json" \
  --request GET \
  --data '{"page":2,"requestsPerPage":2}' \
  http://localhost:8080/v1/requests/list

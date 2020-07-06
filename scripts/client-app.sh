#!/bin/bash

# Valid resource
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"method":"GET","url":"http://google.com"}' \
  http://localhost:8080/v1/requests/request

# Non-existed resource
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"method":"GET","url":"http://kjfslkdjgfnsljkn.com"}' \
  http://localhost:8080/v1/requests/request

# Invalid resource URL
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"method":"GET","url":"invalidresourceurl"}' \
  http://localhost:8080/v1/requests/request

# Invalid request
curl http://localhost:8080/v1/requests/request

# Get all requests from server storage
curl http://localhost:8080/v1/requests/list

# Delete request from server storage
curl --header "Content-Type: application/json" \
  --request DELETE \
  --data '{"id":"d9c1ded7-bd0d-499e-b349-a0586a549562"}' \
  http://localhost:8080/v1/requests/request
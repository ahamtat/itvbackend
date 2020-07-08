#!/bin/bash

swagger validate --stop-on-error ../api/api-swagger.yaml
retVal=$?
if [ $retVal -ne 0 ]; then
  echo "Swagger YAML is not valid"
  exit $retVal
else
  echo "Generating GO code"
  mkdir -p ../internal/app/swagger/server
  swagger generate server -f ../api/api-swagger.yaml -t ../internal/app/swagger/server --exclude-main
fi

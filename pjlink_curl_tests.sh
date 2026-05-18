#!/bin/bash

# Set your environment variables before running
MICROSERVICE_URL="your.microservice.url"
DEVICE_FQDN="your.device.fqdn"

echo "Running PJLink Microservice Tests..."

# GET requests
curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/power"
sleep 1

curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/volume"
sleep 1

curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/videoroute/1"
sleep 1

curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/videomute/1"
sleep 1

curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/audiomute/1"
sleep 1

curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/audioandvideomute/1"
sleep 1

curl -X GET "http://$MICROSERVICE_URL/$DEVICE_FQDN/lamphours"
sleep 1

# PUT requests
curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/power" \
     -H "Content-Type: application/json" \
     -d '"on"'
sleep 10

curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/volume" \
     -H "Content-Type: application/json" \
     -d '"up"'
sleep 1

curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/videoroute/1" \
     -H "Content-Type: application/json" \
     -d '"56"'
sleep 1

curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/videomute" \
     -H "Content-Type: application/json" \
     -d '"true"'
sleep 1

curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/audiomute/1" \
     -H "Content-Type: application/json" \
     -d '"true"'
sleep 1

curl -X PUT "http://$MICROSERVICE_URL/$DEVICE_FQDN/audioandvideomute/1" \
     -H "Content-Type: application/json" \
     -d '"true"'
sleep 1

echo "Tests complete."

#!/bin/bash

# Use the first arg as the sleep time in seconds, default to 5s
time="${1:-5}"

while [ true ]
do
  echo '{"gh_username":"geoff-coppertop", "gateway_ip":"192.168.1.1", "username":"thomasga", "run_rpi_boot":"false", "os_choice":"RaspiOSLite"}'
  sleep $time
done

echo 'All done'
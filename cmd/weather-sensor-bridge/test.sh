#!/bin/bash

# Use the first arg as the sleep time in seconds, default to 5s
time="${1:-5}"

while [ true ]
do
  echo '{"time" : "2021-09-02 12:49:53", "model" : "Test Deivce", "device" : 12, "id" : 0, "batterylow" : 0, "avewindspeed" : 3, "gustwindspeed" : 7, "winddirection" : 340, "cumulativerain" : 54, "temperature" : 808, "humidity" : 86, "light" : 183, "uv" : 0, "mic" : "CRC"}'
  sleep $time
done

echo 'All done'
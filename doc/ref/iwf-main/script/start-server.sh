#!/bin/bash

CONFIG_TEMPLATE_PATH="${CONFIG_TEMPLATE_PATH:-/iwf/config/config_template.yaml}"
SRC_ROOT="${SRC_ROOT:-/iwf}"
HOST=''
TEMPORAL_SERVICE_NAME="${TEMPORAL_SERVICE_NAME:-temporal}"
CADENCE_SERVICE_NAME="${CADENCE_SERVICE_NAME:-cadence}"

if [[ -n "${BACKEND_DEPENDENCY}" && "${BACKEND_DEPENDENCY}" = "cadence" ]]; then
  HOST=$(echo ${CADENCE_HOST_PORT:-"${CADENCE_SERVICE_NAME}:7833"} | sed 's/:/ /')
else
  HOST=$(echo ${TEMPORAL_HOST_PORT:-"${TEMPORAL_SERVICE_NAME}:7233"} | sed 's/:/ /')
fi

RESULT=1
while [[ "${RESULT}" = "1" ]]
do
  nc -z ${HOST}
  RESULT=$?
  if [[ "${RESULT}" = "1" ]]; then
    sleep 3
    echo "Waiting for ${HOST} to be ready..."
  fi
done

for run in {1..60}; do
  sleep 1
  if nc -zv ${HOST}; then
    break
  fi
done

echo "now waiting 10s for server to be ready, so that another script will register namespace/search attributes."
sleep 10
"${SRC_ROOT}/iwf-server" --config "${CONFIG_TEMPLATE_PATH}" start "$@"

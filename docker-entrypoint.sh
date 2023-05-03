#!/usr/bin/env bash

ETP_RT="${ETP_RT}"
CREDENTIALS="${CREDENTIALS}"


if [[ -z ${ETP_RT} ]]; then
  echo "No ETP-RT Variable Set";
else
  /usr/bin/crunchy-cli --etp-rt "${ETP_RT}" login;
fi

if [[ -z ${CREDENTIALS} ]]; then
  echo "No CREDENTIALS Variable Set";
else
  /usr/bin/crunchy-cli --credentials "${CREDENTIALS}" login;
fi

exec "$@"
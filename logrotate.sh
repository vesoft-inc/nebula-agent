#!/bin/env bash

ROTATE=5
SIZE=200M

if [ -n "${LOGROTATE_ROTATE}" ]; then
  ROTATE=${LOGROTATE_ROTATE}
fi

if [ -n "${LOGROTATE_SIZE}" ]; then
  SIZE=${LOGROTATE_SIZE}
fi

nebula="
/usr/local/nebula/logs/graphd-stderr.log
/usr/local/nebula/logs/graphd-out.log
/usr/local/nebula/logs/nebula-graphd.INFO
/usr/local/nebula/logs/nebula-graphd.INFO.impl
/usr/local/nebula/logs/nebula-graphd.WARNING
/usr/local/nebula/logs/nebula-graphd.WARNING.impl
/usr/local/nebula/logs/nebula-graphd.ERROR
/usr/local/nebula/logs/nebula-graphd.ERROR.impl
/usr/local/nebula/logs/metad-stderr.log
/usr/local/nebula/logs/metad-out.log
/usr/local/nebula/logs/nebula-metad.INFO
/usr/local/nebula/logs/nebula-metad.INFO.impl
/usr/local/nebula/logs/nebula-metad.WARNING
/usr/local/nebula/logs/nebula-metad.WARNING.impl
/usr/local/nebula/logs/nebula-metad.ERROR
/usr/local/nebula/logs/nebula-metad.ERROR.impl
/usr/local/nebula/logs/storaged-stderr.log
/usr/local/nebula/logs/storaged-out.log
/usr/local/nebula/logs/nebula-storaged.INFO
/usr/local/nebula/logs/nebula-storaged.INFO.impl
/usr/local/nebula/logs/nebula-storaged.WARNING
/usr/local/nebula/logs/nebula-storaged.WARNING.impl
/usr/local/nebula/logs/nebula-storaged.ERROR
/usr/local/nebula/logs/nebula-storaged.ERROR.impl
{
        daily
        rotate ${ROTATE}
        copytruncate
        nocompress
        missingok
        notifempty
        create 644 root root
        size ${SIZE}
}
"

echo "${nebula}" >/etc/logrotate.d/nebula

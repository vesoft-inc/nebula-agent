FROM ubuntu:20.04

ENV LOGROTATE_ROTATE=5 \
    LOGROTATE_SIZE=100M \
    TZ=Asia/Shanghai

RUN mkdir -p /usr/local/nebula/bin \
    && mkdir -p /usr/local/certs
RUN ln -sf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl cron logrotate \
    && apt-get clean all
COPY bin/agent /usr/local/bin/agent
COPY db_playback /usr/local/nebula/bin/db_playback
COPY logrotate.sh /logrotate.sh
RUN echo "0  *    * * *   root    /etc/cron.daily/logrotate" >> /etc/crontab

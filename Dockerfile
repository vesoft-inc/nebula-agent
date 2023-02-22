FROM ubuntu:20.04

RUN mkdir -p /usr/local/nebula/bin
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
                 && echo "Asia/Shanghai" > /etc/timezone
RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl cron logrotate \
    && apt-get clean all

ENV LOGROTATE_ROTATE= \
    LOGROTATE_SIZE=

COPY bin/agent /usr/local/bin/agent
COPY db_playback /usr/local/nebula/bin/db_playback
COPY logrotate.sh /logrotate.sh
RUN echo "0  *    * * *   root    /etc/cron.daily/logrotate" >> /etc/crontab

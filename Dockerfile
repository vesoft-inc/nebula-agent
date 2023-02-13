FROM ubuntu:20.04

RUN mkdir -p /usr/local/nebula/bin
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
                 && echo "Asia/Shanghai" > /etc/timezone
RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends ca-certificates curl \
    && apt-get clean all
ADD bin/agent /usr/local/bin/agent
ADD db_playback /usr/local/nebula/bin/db_playback

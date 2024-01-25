FROM ubuntu:20.04

RUN mkdir -p /usr/local/nebula/bin \
    && mkdir -p /usr/local/certs
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && apt-get clean all
COPY bin/agent /usr/local/bin/agent
COPY db_playback /usr/local/nebula/bin/db_playback
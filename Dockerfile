FROM centos:7

WORKDIR /nebula-agent

COPY bin/agent .

CMD ["./agent", "--agent", "127.0.0.1:8888", "--meta", "127.0.0.1:9559"]
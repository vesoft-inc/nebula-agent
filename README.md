# Overview

Nebula Agent is an daemon service in each machine of [nebula](https://github.com/vesoft-inc/nebula) cluster. It helps to keep track of nebula metad/storaged/graphd, start/stop them, call local rpc of them.
It is only used for [backup and restore](https://github.com/vesoft-inc/nebula-br) tools for now.

# Quick Start

## Download directly

If you are in linux amd64 environment, you could download it directly.

1. Download the agent.

  ```bash
  $ wget https://github.com/vesoft-inc/nebula-agent/releases/download/v0.1.1/agent-0.1.1-linux-amd64
  ```
2. Change the agent name.
  ```bash
  $ mv agent-v0.1.1 agent
  ```
3. Add execute permission to agent.
  ```bash
  $ chmod +x agent
  ```

## Download repo and build

If you are in other environments, you should first install (golang)[https://go.dev/] 1.16+ and git.
Then you could download the repo and build agent binary yourself.
1. Clone repo.
```bash
$ git clone git@github.com:vesoft-inc/nebula-agent.git
```

2. Change directory to nebula-agent.
```bash
$ cd nebula-agent
```

3. Compile with `make`.
```bash
$ make
```

4. Add execute permission to agent.
```bash
$ cd bin
$ chmod +x agent
```

## Usage

If you want to use nebula-agent, you should start one and **only one** nebula-agent service in each nebula cluster machine which contains metad/storaged/metad services.
It should be given an agent daemon address and the metad address. If you have multi-metad in one nebula cluster, any address of them will be OK.

```bash
Usage of agent:
  --agent string
        The agent server address
  --meta string
        The nebula metad service address, any metad address will be ok
  --debug
        Open debug will output more detail info
  --hbs int
        Agent heartbeat interval to nebula meta, in seconds (default 60)
```

An example:

```bash
./agent --agent="127.0.0.1:8888" --meta="127.0.0.1:9559"
```


# Features

Nebula Agent provide two type of services now: file management in nebula machines and agent service. 

## File Management

```C++
// UploadFile upload file from agent machine to external storage
rpc UploadFile(UploadFileRequest) returns (UploadFileResponse);
// DownloadFile download file from external storage to agent machine
rpc DownloadFile(DownloadFileRequest) returns (DownloadFileResponse);

// MoveDir rename dir in agent machine
rpc MoveDir(MoveDirRequest) returns (MoveDirResponse);
// RemoveDir delete dir in agent machine
rpc RemoveDir(RemoveDirRequest) returns (RemoveDirResponse);
// ExistDir check if dir in agent machine exist
rpc ExistDir(ExistDirRequest) returns (ExistDirResponse);
```

## Agent Service

```C++
// start/stop/get status of metad/storaged/graphd service
rpc StartService(StartServiceRequest) returns (StartServiceResponse);
rpc StopService(StopServiceRequest) returns (StopServiceResponse);
rpc ServiceStatus(ServiceStatusRequest) returns (ServiceStatusResponse);

// ban read/write by call graphd's api
rpc BanReadWrite(BanReadWriteRequest) returns (BanReadWriteResponse);
rpc AllowReadWrite(AllowReadWriteRequest) returns (AllowReadWriteResponse);
```


# Overview

Nebula Agent is an daemon service in each machine of [nebula](https://github.com/vesoft-inc/nebula) cluster. It helps to keep track of nebula metad/storaged/graphd, start/stop them, call local rpc of them.

# Features

## File Management

```proto
// UploadFile upload file from agent machine to external storage
rpc UploadFile(UploadFileRequest) returns (UploadFileResponse);
// DownloadFile download file from external storage to agent machine
rpc DownloadFile(DownloadFileRequest) returns (DownloadFileResponse);
// MoveDir rename dir in agent machine
rpc MoveDir(MoveDirRequest) returns (MoveDirResponse);
// RemoveDir delete dir in agent machine
rpc RemoveDir(RemoveDirRequest) returns (RemoveDirResponse);
```

## Agent Service

```proto
// start/stop metad/storaged/graphd service
rpc StartService(StartServiceRequest) returns (StartServiceResponse);
rpc StopService(StopServiceRequest) returns (StopServiceResponse);

// ban read/write by call graphd's api
rpc BanReadWrite(BanReadWriteRequest) returns (BanReadWriteResponse);
rpc AllowReadWrite(AllowReadWriteRequest) returns (AllowReadWriteResponse);
```

# Usage

Agent will be started in each machine automatically, with it's listen address and metad's address given.

```
Usage of bin/agent:
  --agent string
        The agent server address
  --meta string
        The nebula metad service address, any metad address will be ok
  --debug
        Open debug will output more detail info
  --hbs int
        Agent heartbeat interval to nebula meta, in seconds (default 60)
```

An example: `agent --agent="127.0.0.1:8888" --meta="127.0.0.1:9559"`

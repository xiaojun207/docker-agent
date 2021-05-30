docker agent ,which is an agent post docker info、container list、container stats、container logs to server

### Quick start：
```
docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock -e DockerServer="http://192.168.1.200:8080/dockerApi" -e DockerWsServer="ws://192.168.1.200:8080/dockerApi/ws/" -e Token="12345678" xiaojun207/docker-agent:latest
```

or
```
/app/App -DockerServer $DockerServer -DockerWsServer $DockerWsServer -Token $Token
```


### Env:
- DockerServer: The http server accept the agent post docker info;
- DockerWsServer: The websocket server accept the agent post docker info, and push the task to agent, like create and run a new container;
- Token: The http and websocket header authorization for dockerserver auth;

### DockerServer, application/json, must support api：
- /reg  recive agent post docker info data;
- /containers  recive agent post all container list data;

### DockerWsServer, the web socket server：
the message like this:
```
{
    "ch: "docker.container.create", // the channel for docker, like: docker.image.pull, docker.container.start, docker.container.remove, docker.container.create, docker.container.run (create and start), 
    "ts": 1622367529238, // Millisecond timestamps,
    "d": data for channel
}    
```


### futures:
DockerServer and DockerWsServer will open source

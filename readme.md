docker agent ,which is an agent post docker info、container list、container stats、container logs to server

### Quick start：
```
docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock -e DockerServer="http://192.168.1.200:8080/dockerApi" -e Token="12345678" xiaojun207/docker-agent:latest
```

or
```
/app/App -DockerServer $DockerServer -Token $Token
```


### Env:
- DockerServer: The http server accept the agent post docker info;
- Token: The http and websocket header authorization for dockerserver auth;

### DockerServer, application/json, must support api：
- POST {DockerServer}/reg,  recive agent post docker info data;
- POST {DockerServer}/containers,  recive agent post all container list data;
- WS {DockerServer}/ws, websocket path, The websocket server accept the agent post docker info, and push the task to agent, like create and run a new container;

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

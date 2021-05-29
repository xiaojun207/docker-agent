docker agent ,which is an agent post docker info to server

### Quick start：
```
docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock -e DockerServer="http://192.168.3.67:8388/dockerMgrApi" -e DockerWsServer="ws://192.168.3.67:8388/dockerMgrApi/ws/" -e Token="12345678" xiaojun207/docker-agent:latest
```

or
```
/app/App -DockerServer $DockerServer -DockerWsServer $DockerWsServer -Token $Token
```


### Env:
- DockerServer: The http server accept the agent post docker info;
- DockerWsServer: The websocket server accept the agent post docker info, and push the task to agent, like create and run a new container;
- Token: The http header authorization for dockerserver auth;

> 其中DockerServer，需要支持接口：
> 


### futures:
DockerServer and DockerWsServer will open source

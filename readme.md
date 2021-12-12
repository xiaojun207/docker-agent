[中文说明](#zh) | <a name="en">English</a>

docker agent ,which is an agent post docker info、container list、container stats、container logs to server

### Quick start：
```shell
docker pull xiaojun207/docker-agent:latest

docker run -d --name docker-agent -v /var/run/docker.sock:/var/run/docker.sock -e DockerServer="http://192.168.1.200:8068/dockerMgrApi/agent" -e Username="agent" -e Password="12345678" xiaojun207/docker-agent:latest

```

or
```
/app/App -DockerServer $DockerServer -e Username="agent" -e Password="12345678"
```

### DockerServer(docker-manager)
docker manager
```shell
 docker run -d --name docker-manager -p 8068:8068 -v /app/docker-manager/data:/app/data xiaojun207/docker-manager:latest

```


### Parameter Description:

Parameter | required    | default value | description
---|-------------|--------------|--- 
DockerServer | required    | -            | The http server accept the agent post docker info;
Username | no          | agent        | The username for dockerserver auth, default : agent. You can get the token from DockerServer first start console logs;
Password | required    | false        | The password of username for dockerserver auth. You can get the token from DockerServer first start console logs;
Token | Deprecated  | -            | Deprecated, instead by username and password


### Special note
The hostname of each server (the host of docker-agent) must be unique.

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

## Contact email
If you have any ideas or suggestions, please send an email to the following email:

email: xiaojun207@126.com


## 中文说明

<a name="zh">中文说明</a> | [English](#en)


docker-agent，它是一个将docker信息、容器列表、容器统计信息、容器日志发布到服务器(docker-manager)的代理.

### 快速启动：
```shell
docker pull xiaojun207/docker-agent:latest

docker run -d --name docker-agent -v /var/run/docker.sock:/var/run/docker.sock -e DockerServer="http://192.168.1.200:8068/dockerMgrApi/agent" -e Username="agent" -e Password="12345678" xiaojun207/docker-agent:latest

```

或者
```
/app/App -DockerServer $DockerServer -e Username="agent" -e Password="12345678" 
```

### DockerServer(docker-manager)
服务管理端
```shell
 docker run -d --name docker-manager -p 8068:8068 -v /app/docker-manager/data:/app/data xiaojun207/docker-manager:latest

```

### Parameter Description:

参数 | 是否必填 | 默认值   | 描述
---|------|-------|--- 
DockerServer | 必填   | -     | docker-manager的http地址，用于接收docker-agent提交的docker信息;
Username | no   | agent | 登录docker-manager用户名. 在docker-manager中获取，角色类型：AGENT;
Password | 必填   | -     | 登录docker-manager密码
Token | 已弃用  | -     | 已弃用，请使用username和password替换


### 特别说明
每台服务器(docker-agent的宿主机)的hostname，必须唯一

### docker-manager提交数据格式application/json：
- POST {DockerServer}/reg,  接收docker-agent提交的docker基本信息;
- POST {DockerServer}/containers, 接收docker-agent提交的容器列表数据;
- WS {DockerServer}/ws, 接收docker-agent提交的数据, 推送任务到docker-agent, 比如创建并运行一个新容器;

信息格式如下:
```
{
    "ch: "docker.container.create", // 信息通道, 如: docker.image.pull, docker.container.start, docker.container.remove, docker.container.create, docker.container.run (create and start), 
    "ts": 1622367529238, // 毫秒时间戳,
    "d": data for channel // 具体数据
}    
```


## 联系邮箱
如果，你有什么想法或建议，请你发送邮件到下面的邮箱：

email: xiaojun207@126.com



docker run -d --name docker-agent --restart=always \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e DockerServer="http://192.168.3.67:8068/dockerMgrApi/agent" \
  -e Username="agent" -e Password="FdHF8QoUUQxXLq0CFjQRAwmFutZ8MuXw" \
  xiaojun207/docker-agent:latest

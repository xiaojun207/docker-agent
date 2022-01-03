

docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e DockerServer="http://192.168.3.67:8068/dockerMgrApi/agent" \
  -e Username="agent" -e Password="MOoojLftYZO6NI7mlphwJRYuVCc1k2VW" \
  -e HostIp="192.168.3.67" \
  xiaojun207/docker-agent:1.4.1

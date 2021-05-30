version='1.0.0'

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o build/docker-agent-linux App.go
#upx build/docker-agent-linux

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -o build/docker-agent-darwin App.go
#upx build/docker-agent-darwin

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -o build/docker-agent-window.exe App.go
#upx build/docker-agent-window.exe

cd build
tar -czvf ./docker-agent-linux-amd64-${version}.tar.gz ./docker-agent-linux
tar -czvf ./docker-agent-darwin-amd64-${version}.tar.gz ./docker-agent-darwin
tar -czvf ./docker-agent-window-amd64-${version}.tar.gz ./docker-agent-window.exe

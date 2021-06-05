package dto

type ContainerCreateConfig struct {
	TaskId        string            `json:"taskId"`
	ImageName     string            `json:"imageName"`
	ContainerName string            `json:"containerName"`
	Ports         map[string]string `json:"ports"`
	Env           []string          `json:"env"`
	Volumes       []string          `json:"volumes"`
	Memory        int64             `json:"memory"`
	LogType       string            `json:"logType"`
	LogConfigMap  map[string]string `json:"logConfig"`
}

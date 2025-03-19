package models

type Req struct {
	RepoUrl   string `json:"repoUrl"`
	Framework string `json:"framework"`
}

type R2Config struct {
	AccessKeyID     string
	SecretAccessKey string
	EndPoint        string
}

type RedisObject struct {
	ProjectId string `json:"projectId"`
	Framework string `json:"framework"`
}

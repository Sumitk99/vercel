package models

type Req struct {
	RepoUrl string `json:"repoUrl"`
}

type R2Config struct {
	AccessKeyID     string
	SecretAccessKey string
	EndPoint        string
}

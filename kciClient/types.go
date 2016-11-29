package kciClient

import (
	"time"
)

// ------------------------------------------------------
// User represents a user account.
type User struct {
	ID             int64  `json:"id,omitempty"` // No.
	KUserId        uint64 `json:"userId" `      // qcos app id
	RepoType       string `json:"repoType" `
	RepoUserName   string `json:"repoUserName" ` // github user name
	RepoUserEmail  string `json:"repoUserEmail" `
	RepoUserAvatar string `json:"repoUserAvatar" `
	Active         bool   `json:"active"`
}

// ------------------------------------------------------
// Repo represents a version control repository.
type Repo struct {
	RepoType     string `json:"repoType"`
	RepoOwner    string `json:"repoOwner"`
	RepoName     string `json:"repoName"`
	RepoFullName string `json:"repoFullName"`
	RepoAvatar   string `json:"repoAvatar,omitempty"`
}

// ------------------------------------------------------
// Project represents a project in kci.
type Project struct {
	ID            int64     `json:"id"`            // No.
	KUserID       uint64    `json:"userId"`        // qcos app id
	BuildLocation string    `json:"buildLocation"` // buildLocation reserved for v2
	ProjName      string    `json:"name"`          // project name, 该 KuserID 下唯一
	RepoType      string    `json:"repoType"`      // github
	Label         string    `json:"label,omitempt`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
	// 以下为项目repo属性
	RepoOwner    string `json:"repoOwner"`
	RepoName     string `json:"repoName"`
	RepoFullName string `json:"repoFullName"`
	RepoAvatar   string `json:"repoAvatar,omitempty" `
	RepoLink     string `json:"repoLink,omitempty"`
	RepoKind     string `json:"repoScm,omitempty"`
	RepoClone    string `json:"repoClone,omitempty" `
	RepoBranch   string `json:"repoBranch,omitempty"`
	// 以下为项目构建属性
	Timeout      int64 `json:"timeout" `
	PrActive     bool  `json:"prActive"`
	PushActive   bool  `json:"pushActive"`
	DeployActive bool  `json:"deployActive"`
	TagsActive   bool  `json:"tagsActive" `
}

type CreateProjReq struct {
	BuildLocation string `json:"buildLocation"` // buildLocation reserved for v2
	ProjName      string `json:"name"`          // project name, 该 KuserID 下唯一
	RepoType      string `json:"repoType"`      // github
	RepoOwner     string `json:"repoOwner"`
	RepoName      string `json:"repoName"`
}

type CheckProjNameRes struct {
	Avaliable bool `json:"avaliable"`
}

type PatchProj struct {
	Timeout      *int64 `json:"timeout"`
	PrActive     *bool  `json:"prActive"`
	PushActive   *bool  `json:"pushActive"`
	DeployActive *bool  `json:"deployActive"`
	TagsActive   *bool  `json:"tagsActive" `
}

// ------------------------------------------------------
// Build represents the process of compiling and testing work
type Build struct {
	Number   int       `json:"number" `
	Event    string    `json:"event"`
	Status   string    `json:"status"`
	Enqueued time.Time `json:"enqueued"`
	Created  time.Time `json:"created"`
	Started  time.Time `json:"started"`
	Finished time.Time `json:"finished"`
	//Deploy_to `json:""`
	Commit  string `json:"commit"`
	Branch  string `json:"branch"`
	Ref     string `json:"ref"`
	Refspec string `json:"refspec"`
	Remote  string `json:"remote"`
	Title   string `json:"title"`
	Message string `json:"message"`
	//Timestamp 0 `json:""`
	Author       string `json:"author"`
	AuthorAvatar string `json:"authorAvatar"`
	AuthorEmail  string `json:"authorEmail"`
	LinkUrl      string `json:"linkUrl"`
	Signed       bool   `json:"signed"`
	Verified     bool   `json:"verified"`
	//add
	ProjectId int64  `json:"projectId"`
	Jobs      []*Job `json:"jobs,omitempty"`
}

type Job struct {
	Number   int    `json:"number"`
	Error    string `json:"error"`
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode"`
	Enqueued int64  `json:"enqueued"`
	Started  int64  `json:"started"`
	Finished int64  `json:"finished"`

	Environment map[string]string `json:"environment"`
}

// ------------------------------------------------------
// Log represents a line of log during build
type Log struct {
	Proc string
	Time int
	Pod  int
	Out  string
}

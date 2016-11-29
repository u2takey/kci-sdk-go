package kciClient

// Client describes a kci client.
type Client interface {
	// 返回用户信息（绑定的子帐户信息）
	Self() ([]*User, error)

	// 获取仓库列表
	RepoList(repoType string) ([]*Repo, error)

	// 创建项目
	ProjPost(req *CreateProjReq) (*Project, error)

	// 获取项目列表
	ProjList() ([]*Project, error)

	// 获取项目
	Proj(projId int64) (*Project, error)

	// 更新项目设置
	ProjPatch(projId int64, p *PatchProj) (*Project, error)

	// 删除项目
	ProjDel(projId int64) error

	// 手动构建
	BuildPost(projId int64, branch string) (*Build, error)

	// 获取构建历史
	BuildList(projId int64) ([]*Build, error)

	// 获取单次的构建
	BuildById(projId int64, buildNum int) (*Build, error)

	// 获取某次构建的日志
	BuildLogs(projId int64, buildNum, jobNum int) ([]*Log, error)

	// 解除绑定
	AuthDel(repoType string) error

	// 检查项目名是否可用
	CheckProjName(name string) (*CheckProjNameRes, error)

	// 实时日志
	FeedWs(userid uint64) (<-chan []byte, error)
	LogWs(projId int64, buildId, jobNum int) (<-chan []byte, error)
}

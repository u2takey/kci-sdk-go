package kciClient

import (
	"encoding/json"
	"fmt"
	"github.com/franela/goblin"
	"testing"
)

const defaultTimeout = 60
const defaultPrActive = false
const defaultPushActive = true
const defaultDeployActive = true
const defaultTagsActive = false

// !! warning !! test will del/add/mod projects
func TestClient(t *testing.T) {
	g := goblin.Goblin(t)

	client := NewClient(kciUrl, ak, sk)
	g.Describe("TestClient", func() {
		// !!! you should auth with portal first !!!
		var projIdForBuildTest int64
		var buildNum int
		g.Before(func() {
			projs, _ := client.ProjList()
			for _, proj := range projs {
				client.ProjDel(proj.ID)
			}
			req := &CreateProjReq{
				ProjName:  "justtest",
				RepoType:  "github",
				RepoOwner: "u2takey",
				RepoName:  "justtest",
			}

			proj, err := client.ProjPost(req)
			g.Assert(err == nil).IsTrue()

			proj, err = client.Proj(proj.ID)
			g.Assert(err == nil).IsTrue()
			g.Assert(proj.ProjName == req.ProjName).IsTrue()
			g.Assert(proj.RepoType == req.RepoType).IsTrue()
			g.Assert(proj.RepoOwner == req.RepoOwner).IsTrue()
			g.Assert(proj.RepoName == req.RepoName).IsTrue()
			projIdForBuildTest = proj.ID
		})

		// ------------------------------------------------
		g.It("Should get user count info", func() {
			users, err := client.Self()
			g.Assert(err == nil).IsTrue()
			//Detail(users)
			g.Assert(len(users) > 0).IsTrue()
		})

		// ------------------------------------------------
		g.It("Should get user repos", func() {
			repos, err := client.RepoList("github")
			g.Assert(err == nil).IsTrue()
			//Detail(repos)
			g.Assert(len(repos) > 0).IsTrue()
		})

		// ------------------------------------------------
		g.It("Should create and get projs", func() {
			projs, err := client.ProjList()
			g.Assert(err == nil).IsTrue()
			// Detail(projs)
			// should be 1
			g.Assert(len(projs) == 1).IsTrue()

			req := &CreateProjReq{
				ProjName:  "testname123",
				RepoType:  "github",
				RepoOwner: "u2takey",
				RepoName:  "go-cache",
			}
			proj, err := client.ProjPost(req)
			// should success
			g.Assert(err == nil).IsTrue()
			// Detail(proj)
			g.Assert(proj.ProjName == req.ProjName).IsTrue()
			g.Assert(proj.RepoType == req.RepoType).IsTrue()
			g.Assert(proj.RepoOwner == req.RepoOwner).IsTrue()
			g.Assert(proj.RepoName == req.RepoName).IsTrue()

			proj, err = client.ProjPost(req)
			// should fail for repos exsit
			g.Assert(err != nil).IsTrue()

			projs, err = client.ProjList()
			g.Assert(err == nil).IsTrue()
			// Detail(projs)
			// should have 2 projs
			g.Assert(len(projs) == 2).IsTrue()
		})

		// ------------------------------------------------
		g.It("Should not create projs", func() {
			req := &CreateProjReq{
				ProjName:  "justtest",
				RepoType:  "github",
				RepoOwner: "u2takey",
				RepoName:  "boom",
			}
			_, err := client.ProjPost(req)
			// should fail for porj name exsit
			g.Assert(err != nil).IsTrue()

			checkProjNameRes, err := client.CheckProjName(req.ProjName)
			g.Assert(err != nil).IsTrue()
			g.Assert(checkProjNameRes.Avaliable == false).IsTrue()

			req = &CreateProjReq{
				ProjName:  "testname",
				RepoType:  "github",
				RepoOwner: "u2takey",
				RepoName:  "justtest",
			}
			_, err = client.ProjPost(req)
			// should fail for repos exsit
			g.Assert(err != nil).IsTrue()

		})

		// ------------------------------------------------
		g.It("Should update project info ", func() {
			proj, err := client.Proj(projIdForBuildTest)
			g.Assert(err == nil).IsTrue()
			// this is default config
			g.Assert(proj.Timeout == defaultTimeout).IsTrue()
			g.Assert(proj.PrActive == defaultPrActive).IsTrue()
			g.Assert(proj.PushActive == defaultPushActive).IsTrue()
			g.Assert(proj.DeployActive == defaultDeployActive).IsTrue()
			g.Assert(proj.TagsActive == defaultTagsActive).IsTrue()

			patchReq := &PatchProj{}
			var timeout int64 = 30
			tagActive := true
			patchReq.Timeout = &timeout
			// Detail(patchReq)
			patchReq.TagsActive = &tagActive

			proj, err = client.ProjPatch(projIdForBuildTest, patchReq)
			g.Assert(err == nil).IsTrue()

			g.Assert(proj.Timeout == *patchReq.Timeout).IsTrue()
			g.Assert(proj.PrActive == defaultPrActive).IsTrue()
			g.Assert(proj.PushActive == defaultPushActive).IsTrue()
			g.Assert(proj.DeployActive == defaultDeployActive).IsTrue()
			g.Assert(proj.TagsActive == *patchReq.TagsActive).IsTrue()
		})

		// ------------------------------------------------
		g.It("Should post and get build info ", func() {
			proj, err := client.Proj(projIdForBuildTest)
			g.Assert(err == nil).IsTrue()
			build, err := client.BuildPost(projIdForBuildTest, "default")
			//Detail(build)
			g.Assert(err == nil).IsTrue()
			g.Assert(build.ProjectId == proj.ID).IsTrue()

			builds, err := client.BuildList(projIdForBuildTest)
			g.Assert(err == nil).IsTrue()
			//Detail(builds)
			g.Assert(len(builds) == 1).IsTrue()

			buildDetail, err := client.BuildById(projIdForBuildTest, build.Number)
			g.Assert(err == nil).IsTrue()
			//Detail(buildDetail)
			g.Assert(buildDetail.ProjectId == projIdForBuildTest).IsTrue()
			g.Assert(len(buildDetail.Jobs) > 0).IsTrue()

			buildNum = build.Number
		})

		// ------------------------------------------------
		g.After(func() {
			// should get logs after build done
			//logs, err := client.BuildLogs(projIdForBuildTest, buildNum, 0)
			//g.Assert(err == nil).IsTrue()
			//g.Assert(len(logs) > 0).IsTrue()
		})
	})
}

func Detail(v interface{}) {
	buf, _ := json.MarshalIndent(v, "  ", " ")
	fmt.Println(string(buf))
}

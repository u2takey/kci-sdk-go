package kciClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	pathSelf          = "%s/v1/user"
	pathRepo          = "%s/v1/user/%s/repo"
	pathProj          = "%s/v1/project"
	pathCheckProjName = "%s/v1/info/checkname/%s"
	pathProjById      = "%s/v1/project/%d"
	pathBuild         = "%s/v1/build/%d/%s"
	pathBuildList     = "%s/v1/build/%d"
	pathBuildById     = "%s/v1/build/%d/%d"
	pathBuildLogById  = "%s/v1/build/%d/%d/%d/log"
	pathAuth          = "%s/v1/%s/auth"
	pathFeedWs        = "%s/ws/feed/%d"
	pathRealLogs      = "%s/ws/log/%d/%d/%d"

	httpScheme = "https://"
	wsScheme   = "wss://"
)

type client struct {
	client *http.Client
	base   string // base url
	wsbase string
	config *ClientConfig
}

type ClientConfig struct {
	Host      string
	AK        string
	SK        string
	Transport http.RoundTripper
	UserAgent string
}

// NewClient returns a client at the specified url.
func NewClient(host, ak, sk string) Client {
	config := &ClientConfig{
		Host:      host,
		AK:        ak,
		SK:        sk,
		UserAgent: "KCISDK / " + sdkVersion,
	}
	return NewClientWithConfig(config)
}

// NewClient returns a client at the specified url.
func NewClientWithConfig(config *ClientConfig) Client {
	c := &client{base: httpScheme + config.Host, wsbase: wsScheme + config.Host, config: config}
	m := NewMac(config.AK, config.SK)
	c.client = NewMacClient(m, config.Transport)
	return c
}

// 返回用户信息（绑定的子帐户信息）
func (c *client) Self() ([]*User, error) {
	var out []*User
	uri := fmt.Sprintf(pathSelf, c.base)
	err := c.get(uri, &out)
	return out, err
}

// 获取仓库列表
func (c *client) RepoList(repoType string) ([]*Repo, error) {
	var out []*Repo
	uri := fmt.Sprintf(pathRepo, c.base, repoType)
	err := c.get(uri, &out)
	return out, err
}

// 创建项目
func (c *client) ProjPost(req *CreateProjReq) (*Project, error) {
	out := new(Project)
	uri := fmt.Sprintf(pathProj, c.base)
	err := c.post(uri, req, &out)
	return out, err
}

// 获取项目列表
func (c *client) ProjList() ([]*Project, error) {
	var out []*Project
	uri := fmt.Sprintf(pathProj, c.base)
	err := c.get(uri, &out)
	return out, err
}

// 获取项目
func (c *client) Proj(projId int64) (*Project, error) {
	out := new(Project)
	uri := fmt.Sprintf(pathProjById, c.base, projId)
	err := c.get(uri, &out)
	return out, err
}

// 更新项目设置
func (c *client) ProjPatch(projId int64, p *PatchProj) (*Project, error) {
	out := new(Project)
	uri := fmt.Sprintf(pathProjById, c.base, projId)
	err := c.post(uri, p, &out)
	return out, err
}

// 删除项目
func (c *client) ProjDel(projId int64) error {
	uri := fmt.Sprintf(pathProjById, c.base, projId)
	err := c.delete(uri)
	return err
}

// 手动构建
func (c *client) BuildPost(projId int64, branch string) (*Build, error) {
	out := new(Build)
	uri := fmt.Sprintf(pathBuild, c.base, projId, branch)
	err := c.post(uri, nil, &out)
	return out, err
}

// 获取构建历史
func (c *client) BuildList(projId int64) ([]*Build, error) {
	var out []*Build
	uri := fmt.Sprintf(pathBuildList, c.base, projId)
	err := c.get(uri, &out)
	return out, err
}

// 获取单次的构建
func (c *client) BuildById(projId int64, buildId int) (*Build, error) {
	out := new(Build)
	uri := fmt.Sprintf(pathBuildById, c.base, projId, buildId)
	err := c.get(uri, &out)
	return out, err
}

// 获取某次构建的日志
func (c *client) BuildLogs(projId int64, buildId, jobNum int) ([]*Log, error) {
	var out []*Log
	uri := fmt.Sprintf(pathBuildLogById, c.base, projId, buildId, jobNum)
	err := c.get(uri, &out)
	return out, err
}

// 解除绑定
func (c *client) AuthDel(repoType string) error {
	uri := fmt.Sprintf(pathAuth, c.base, repoType)
	err := c.delete(uri)
	return err
}

// 检查项目名是否可用
func (c *client) CheckProjName(name string) (*CheckProjNameRes, error) {
	out := new(CheckProjNameRes)
	uri := fmt.Sprintf(pathCheckProjName, c.base, name)
	err := c.get(uri, &out)
	return out, err
}

//
func (p *client) FeedWs(userid uint64) (<-chan []byte, error) {
	uri := fmt.Sprintf(pathFeedWs, p.wsbase, userid)
	return p.wsMessages(uri)
}

func (p *client) LogWs(projId int64, num, job int) (<-chan []byte, error) {
	uri := fmt.Sprintf(pathRealLogs, p.wsbase, projId, num, job)
	return p.wsMessages(uri)
}

func (p *client) wsMessages(uri string) (<-chan []byte, error) {
	dailer := &websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
	}
	header := make(http.Header)
	if p.config.UserAgent != "" {
		header["User-Agent"] = []string{p.config.UserAgent}
	}
	c, _, err := dailer.Dial(uri, header)
	if err != nil {
		return nil, err
	}
	msg := make(chan []byte, 10)

	go func() {
		defer c.Close()
		defer close(msg)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			msg <- message
		}
	}()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}
			}
		}
	}()
	return msg, nil
}

//
// http request helper functions
//

// helper function for making an http GET request.
func (c *client) get(rawurl string, out interface{}) error {
	return c.do(rawurl, "GET", nil, out)
}

// helper function for making an http POST request.
func (c *client) post(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "POST", in, out)
}

// helper function for making an http PUT request.
func (c *client) put(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "PUT", in, out)
}

// helper function for making an http PATCH request.
func (c *client) patch(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "PATCH", in, out)
}

// helper function for making an http DELETE request.
func (c *client) delete(rawurl string) error {
	return c.do(rawurl, "DELETE", nil, nil)
}

// helper function to make an http request
func (c *client) do(rawurl, method string, in, out interface{}) error {
	// executes the http request and returns the body as
	// and io.ReadCloser
	body, err := c.stream(rawurl, method, in, out)
	if err != nil {
		return err
	}
	defer body.Close()

	// if a json response is expected, parse and return
	// the json response.
	if out != nil {
		return json.NewDecoder(body).Decode(out)
	}
	return nil
}

// helper function to stream an http request
func (c *client) stream(rawurl, method string, in, out interface{}) (io.ReadCloser, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	// if we are posting or putting data, we need to
	// write it to the body of the request.
	var buf io.ReadWriter
	if in == nil {
		// nothing
	} else if rw, ok := in.(io.ReadWriter); ok {
		buf = rw
	} else {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(in)
		if err != nil {
			return nil, err
		}
	}

	// creates a new http request to bitbucket.
	req, err := http.NewRequest(method, uri.String(), buf)
	if c.config.UserAgent != "" {
		req.Header.Set("User-Agent", c.config.UserAgent)
	}
	if err != nil {
		return nil, err
	}
	if in == nil {
		// nothing
	} else if _, ok := in.(io.ReadWriter); ok {
		req.Header.Set("Content-Type", "plain/text")
	} else {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > http.StatusPartialContent {
		defer resp.Body.Close()
		out, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf(string(out))
	}
	return resp.Body, nil
}

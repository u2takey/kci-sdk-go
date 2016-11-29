package kciClient

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"
	"sort"

	"qiniupkg.com/x/bytes.v7/seekable"
)

type Mac struct {
	AccessKey string
	SecretKey []byte
}

func (m *Mac) SignRequest(req *http.Request) (err error) {

	sign, err := signRequest(m.SecretKey, req)
	if err != nil {
		return
	}

	auth := "Qiniu " + m.AccessKey + ":" + base64.URLEncoding.EncodeToString(sign)
	req.Header.Set("Authorization", auth)
	return
}

type Transport struct {
	mac       Mac
	Transport http.RoundTripper
}

func NewMac(accessKey, secretKey string) *Mac {
	return &Mac{accessKey, []byte(secretKey)}
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	err = t.mac.SignRequest(req)
	if err != nil {
		return
	}

	return t.Transport.RoundTrip(req)
}

func NewTransport(mac *Mac, transport http.RoundTripper) *Transport {

	if transport == nil {
		transport = http.DefaultTransport
	}
	t := &Transport{Transport: transport}
	if mac == nil {
		t.mac.AccessKey = ""
		t.mac.SecretKey = []byte("")
	} else {
		t.mac = *mac
	}
	return t
}

func NewMacClient(mac *Mac, transport http.RoundTripper) *http.Client {

	t := NewTransport(mac, transport)
	return &http.Client{Transport: t}
}

const qiniuHeaderPrefix = "X-Qiniu-"

// ---------------------------------------------------------------------------------------

func incBody(req *http.Request, ctType string) bool {

	return req.ContentLength != 0 && req.Body != nil && ctType != "" && ctType != "application/octet-stream"
}

func signQiniuHeaderValues(header http.Header, w io.Writer) {

	var keys []string
	for key := range header {
		if len(key) > len(qiniuHeaderPrefix) && key[:len(qiniuHeaderPrefix)] == qiniuHeaderPrefix {
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return
	}

	if len(keys) > 1 {
		sort.Sort(sortByHeaderKey(keys))
	}
	for _, key := range keys {
		io.WriteString(w, "\n"+key+": "+header.Get(key))
	}
}

func signRequest(sk []byte, req *http.Request) ([]byte, error) {

	h := hmac.New(sha1.New, sk)

	u := req.URL
	data := req.Method + " " + u.Path
	if u.RawQuery != "" {
		data += "?" + u.RawQuery
	}
	io.WriteString(h, data+"\nHost: "+req.Host)

	ctType := req.Header.Get("Content-Type")
	if ctType != "" {
		io.WriteString(h, "\nContent-Type: "+ctType)
	}

	signQiniuHeaderValues(req.Header, h)

	io.WriteString(h, "\n\n")

	if incBody(req, ctType) {
		s2, err2 := seekable.New(req)
		if err2 != nil {
			return nil, err2
		}
		h.Write(s2.Bytes())
	}

	return h.Sum(nil), nil
}

// ---------------------------------------------------------------------------------------

type sortByHeaderKey []string

func (p sortByHeaderKey) Len() int           { return len(p) }
func (p sortByHeaderKey) Less(i, j int) bool { return p[i] < p[j] }
func (p sortByHeaderKey) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// ---------------------------------------------------------------------------------------

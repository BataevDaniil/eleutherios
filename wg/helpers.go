package wg

import (
	"io"
	"net/http"
	"os/exec"
)

func httpGet(url string) (*http.Response, error) { return http.Get(url) }

func httpPost(url, body string) (*http.Response, error) {
	return http.Post(url, "application/json", io.NopCloser(newStrReader(body)))
}

func execCmd(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	return string(out), err
}

type strReader struct{ s string }

func newStrReader(s string) *strReader { return &strReader{s} }
func (r *strReader) Read(p []byte) (int, error) {
	n := copy(p, r.s)
	r.s = r.s[n:]
	if len(r.s) == 0 {
		return n, io.EOF
	}
	return n, nil
}

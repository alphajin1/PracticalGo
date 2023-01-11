package pkgregister

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type pkgRegisterResult struct {
	ID string `json:"id"`
}

func registerPackageData(url string, data pkgData) (pkgRegisterResult, error) {
	p := pkgRegisterResult{}
	b, err := json.Marshal(data) // 1. data to json
	if err != nil {
		return p, err
	}
	reader := bytes.NewReader(b) // 2. json to Bytes & Post
	r, err := http.Post(url, "application/json", reader)
	if err != nil {
		return p, err
	}
	defer r.Body.Close()
	respData, err := io.ReadAll(r.Body) // 3. check response
	if err != nil {
		return p, err
	}
	if r.StatusCode != http.StatusOK { // 4. check status code
		return p, errors.New(string(respData))
	}
	err = json.Unmarshal(respData, &p)
	return p, err
}

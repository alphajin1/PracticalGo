package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
)

func createMultiPartMessage(data pkgData) ([]byte, string, error) {
	var b bytes.Buffer
	var err error
	var fw io.Writer

	mw := multipart.NewWriter(&b)

	fw, err = mw.CreateFormField("name")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Name)

	fw, err = mw.CreateFormField("version")
	if err != nil {
		return nil, "", err
	}
	fmt.Fprintf(fw, data.Version)

	fw, err = mw.CreateFormFile("filedata", data.Filename)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(fw, data.Bytes)
	err = mw.Close()
	if err != nil {
		return nil, "", err
	}

	contentType := mw.FormDataContentType()
	return b.Bytes(), contentType, nil
}

package cmd

import "errors"

var ErrNoServerSpecified = errors.New("you have to specify the remote server")
var UnSupportedHTTPMethod = errors.New("unsupported HTTP method")
var UnSupportedOutputFormat = errors.New("unsupported Output Format")

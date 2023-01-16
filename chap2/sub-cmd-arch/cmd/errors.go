package cmd

import "errors"

// 사용자 지정 에러
var ErrNoServerSpecified = errors.New("You have to specify the remote server.")
var InvalidHttpMethod = errors.New("Invalid HTTP method")

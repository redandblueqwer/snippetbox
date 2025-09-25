package models

import (
	"errors"
)

// 自定义错误类型的变量

var (
	// 找不到对应ID的snippet记录所用错误
	ErrNoRecord = errors.New("models: no matching record found")

	// 无效的认证信息错误
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// 邮箱重复错误
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

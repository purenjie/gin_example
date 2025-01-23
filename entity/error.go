package entity

import "errors"

const (
	ErrCodeRequest = 1
	ErrCodeSystem  = 2
)

var (
	ErrApp  = errors.New("非法的app")
	ErrTS   = errors.New("非法的ts")
	ErrSign = errors.New("非法的sign")

	ErrInvalidUser      = errors.New("非法的用户名")
	ErrInvalidMessage   = errors.New("非法的消息")
	ErrInvalidModelType = errors.New("非法的模型类型")
	ErrSystem           = errors.New("系统内部错误")
)

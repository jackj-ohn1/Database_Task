package errno

import "github.com/pkg/errors"

var (
	ErrWrongPassword = errors.New("账号或密码错误")
	ErrUserNotExist  = errors.New("用户不存在")
	ErrUserExist     = errors.New("该账号已存在")
	ErrWrongParams   = errors.New("参数错误")
)

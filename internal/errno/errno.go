package errno

import "github.com/pkg/errors"

var (
	ErrWrongPassword     = errors.New("账号或密码错误")
	ErrUserNotExist      = errors.New("用户不存在")
	ErrUserExist         = errors.New("该账号已存在")
	ErrWrongParams       = errors.New("参数错误")
	ErrInternalServerErr = errors.New("服务端出错")
	ErrUsedBook          = errors.New("这本书已被借出")
	ErrSpareBook         = errors.New("这本书尚未外借")
	ErrNoPower           = errors.New("你无权进行操作")
	ErrNoLeftResource    = errors.New("已达借书上限")
)

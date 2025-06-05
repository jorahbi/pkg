package errorx

import "errors"

//go:generate stringer -type State -linecomment
const (
	UNPACK_TOKEN_ERROR    State = 1002 //非法请求
	MISS_PARAMS           State = 1003 //缺少参数
	REPETITION_SUBMIT     State = 1004 //重复提交
	UNAME_OR_PWD_IS_ERROR State = 1005 //用户名或密码错误

	//system
	UPLOAD_FILE_MAX //超过文件最大限制
)

type State int

func (e State) Code() int {
	return int(e)
}

var (
	ErrQueueVaildateIsFail = errors.New("queue data vaildate is fail")
	ErrQueuePush           = errors.New("push queue is fail")
	ErrJsonEncode          = errors.New("data json encode error")
)

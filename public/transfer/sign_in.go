package transfer

const (
	CodeSignInSuccess = 1
	CodeSignInFail    = 2
)

// SignIn 设备登录
type SignIn struct {
	DeviceId int64  `json:"device_id"` // 设备id
	UserId   int64  `json:"user_id"`   // 用户id
	Token    string `json:"token"`     // token
}

//  SignInACK 设备登录回执
type SignInACK struct {
	Code    int    `json:"code"`    // 设备id
	Message string `json:"message"` // 用户id
}

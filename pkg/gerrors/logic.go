package gerrors

var (
	ErrDeviceNotFound     = newError(10200, "设备不存在")
	ErrFriendNotFound     = newError(10210, "好友关系不存在")
	ErrAlreadyIsFriend    = newError(10211, "对方已经是好友了")
	ErrGroupNotFound      = newError(10220, "群组不存在")
	ErrNotInGroup         = newError(10201, "用户没有在群组中")
	ErrGroupUserNotFound  = newError(10202, "群组成员不存在")
	ErrUserAlreadyInGroup = newError(10203, "用户已经在群组中")
)

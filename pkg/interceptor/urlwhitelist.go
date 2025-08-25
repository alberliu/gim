package interceptor

var UserWhitelistURL = map[string]struct{}{
	"/user.UserExtService/SignIn": {},
}

var LogicWhitelistURL = map[string]struct{}{
	"/logic.DeviceExtService/RegisterDevice": {},
}

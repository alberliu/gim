package urlwhitelist

var User = map[string]struct{}{
	"/user.UserExtService/SignIn": {},
}

var Logic = map[string]struct{}{
	"/logic.DeviceExtService/RegisterDevice": {},
}

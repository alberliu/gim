package server

import "gim/pkg/protocol/pb/businesspb"

var URLWhitelist = map[string]struct{}{
	businesspb.UserExtService_SignIn_FullMethodName: {},
}

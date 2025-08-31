package interceptor

import "gim/pkg/protocol/pb/businesspb"

var UserWhitelistURL = map[string]struct{}{
	businesspb.UserExtService_SignIn_FullMethodName: {},
}

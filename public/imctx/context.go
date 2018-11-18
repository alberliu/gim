package imctx

import "goim/public/session"

type Context struct {
	Session *session.Session
}

func NewContext(Session *session.Session) *Context {
	return &Context{Session: Session}
}

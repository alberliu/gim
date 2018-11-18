package connect

import (
	"goim/logic/db"
	"goim/public/imctx"
)

func Context() *imctx.Context {
	return imctx.NewContext(db.Factoty.GetSession())
}

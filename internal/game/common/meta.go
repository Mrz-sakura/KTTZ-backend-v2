package common

import "kttz-server/types/protocol"

type Context struct {
	Opts   *protocol.DeskOptions
	RoomNo string
	Uid    int64

	Current bool // 是否该这个玩家操作

	IsSTTH bool
	IsHL   bool
	IsDS   bool
	IsXS   bool
	IsKT   bool
}

func (c *Context) Reset() {
	c.IsSTTH = false
	c.IsHL = false
	c.IsXS = false
	c.IsKT = false
	c.IsDS = false

	c.Current = false

}

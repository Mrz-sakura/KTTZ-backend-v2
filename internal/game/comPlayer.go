package game

import (
	"github.com/lonng/nano/session"
	log "github.com/sirupsen/logrus"
	"kttz-server/internal/game/common"
	"kttz-server/types/vars"
)

type Player struct {
	uid  int64  // 用户ID
	head string // 头像地址
	name string // 玩家名字
	ip   string // ip地址
	sex  int    // 性别
	coin int64  // 房卡数量

	// 玩家数据
	session *session.Session

	ctx *common.Context

	desk  *Desk //当前桌
	score int   //经过n局后,当前玩家余下的分值数,默认为1000

	logger *log.Entry // 日志
}

func NewPlayer(s *session.Session, uid int64, name, head, ip string, sex int) *Player {
	p := &Player{
		uid:   uid,
		name:  name,
		head:  head,
		ctx:   &common.Context{Uid: uid},
		ip:    ip,
		sex:   sex,
		score: 1000,

		logger: log.WithField(vars.PLAYER, uid),
	}

	p.ctx.Reset()
	p.bindSession(s)

	return p
}

func (p *Player) bindSession(s *session.Session) {
	p.session = s
	p.session.Set(vars.PLAYER, p)
}

func (p *Player) removeSession() {
	p.session.Remove(vars.PLAYER)
	p.session = nil
}

package game

import (
	"github.com/lonng/nano/scheduler"
	"kttz-server/pkg/errutil"
	"kttz-server/types/protocol"
	"kttz-server/types/vars"

	"time"

	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
	log "github.com/sirupsen/logrus"
)

const kickResetBacklog = 8

var defaultComonents = NewCom()

type (
	Com struct {
		component.Base
		group      *nano.Group       // 广播channel
		players    map[int64]*Player // 所有的玩家
		chKick     chan int64        // 退出队列
		chReset    chan int64        // 重置队列
		chRecharge chan RechargeInfo // 充值信息
	}

	RechargeInfo struct {
		Uid  int64 // 用户ID
		Gold int64 // 房卡数量
	}
)

func NewCom() *Com {
	return &Com{
		group:      nano.NewGroup("_SYSTEM_MESSAGE_BROADCAST"),
		players:    map[int64]*Player{},
		chKick:     make(chan int64, kickResetBacklog),
		chReset:    make(chan int64, kickResetBacklog),
		chRecharge: make(chan RechargeInfo, 32),
	}
}

func (m *Com) AfterInit() {
	session.Lifetime.OnClosed(func(s *session.Session) {
		_ = m.group.Leave(s)
	})

	// 处理踢出玩家和重置玩家消息(来自http)
	scheduler.NewTimer(time.Second, func() {
	ctrl:
		for {
			select {
			case uid := <-m.chKick:
				p, ok := defaultComonents.player(uid)
				if !ok || p.session == nil {
					logger.Errorf("玩家%d不在线", uid)
				}
				p.session.Close()
				logger.Infof("踢出玩家, UID=%d", uid)

			case uid := <-m.chReset:
				p, ok := defaultComonents.player(uid)
				if !ok {
					return
				}
				if p.session != nil {
					logger.Errorf("玩家正在游戏中，不能重置: %d", uid)
					return
				}
				p.desk = nil
				logger.Infof("重置玩家, UID=%d", uid)

			case ri := <-m.chRecharge:
				player, ok := m.player(ri.Uid)
				// 如果玩家在线
				if s := player.session; ok && s != nil {
					//s.Push("onCoinChange", &protocol.CoinChangeInformation{Coin: ri.Coin})
				}

			default:
				break ctrl
			}
		}
	})
}

func (m *Com) Login(s *session.Session, req *protocol.LoginRequest) error {
	uid := req.Uid
	_ = s.Bind(uid)

	log.Infof("玩家: %d登录: %+v", uid, req)
	if p, ok := m.player(uid); !ok {
		log.Infof("玩家: %d不在线，创建新的玩家", uid)
		p = NewPlayer(s, uid, req.Nickname, req.Avatar, req.IP, req.Gender)
		m.setPlayer(uid, p)
	} else {
		log.Infof("玩家: %d已经在线", uid)
		// 重置之前的session
		if prevSession := p.session; prevSession != nil && prevSession != s {
			// 移除广播频道
			_ = m.group.Leave(prevSession)

			// 如果之前房间存在，则退出来
			if p, err := m.getPlayerBySession(prevSession); err == nil && p != nil && p.desk != nil && p.desk.group != nil {
				_ = p.desk.group.Leave(prevSession)
			}

			prevSession.Clear()
			prevSession.Close()
		}

		// 绑定新session
		p.bindSession(s)
	}

	// 添加到广播频道
	_ = m.group.Add(s)

	res := &protocol.LoginResponse{
		Uid:      s.UID(),
		Nickname: req.Nickname,
		Gold:     req.Gold,
		Avatar:   req.Avatar,
		Gender:   req.Gender,
	}

	return s.Response(res)
}

func (m *Com) player(uid int64) (*Player, bool) {
	p, ok := m.players[uid]

	return p, ok
}

func (m *Com) setPlayer(uid int64, p *Player) {
	if _, ok := m.players[uid]; ok {
		log.Warnf("玩家已经存在，正在覆盖玩家， UID=%d", uid)
	}
	m.players[uid] = p
}

func (m *Com) playerCount() int {
	return len(m.players)
}

func (m *Com) offline(uid int64) {
	delete(m.players, uid)
	log.Infof("玩家: %d从在线列表中删除, 剩余：%d", uid, len(m.players))
}

func (m *Com) getPlayerBySession(s *session.Session) (*Player, error) {
	p, ok := s.Value(vars.PLAYER).(*Player)
	if !ok {
		return nil, errutil.ErrPlayerNotFound
	}
	return p, nil
}

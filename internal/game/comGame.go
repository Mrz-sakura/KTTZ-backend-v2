package game

import (
	"github.com/lonng/nano"
	log "github.com/sirupsen/logrus"
	"kttz-server/types/protocol"
	"kttz-server/types/vars"
)

type Desk struct {
	roomID    int64                 // 房间号
	deskID    int64                 // desk表的pk
	opts      *protocol.DeskOptions // 房间选项
	state     vars.DeskStatus       // 状态
	round     uint32                // 第n局
	creator   int64                 // 创建玩家UID
	createdAt int64                 // 创建时间
	players   []*Player
	room      *Room
	group     *nano.Group // 组播通道
	die       chan struct{}

	//scores    map[*Player]

	logger *log.Entry
}

package game

import (
	"fmt"
	"github.com/lonng/nano/pipeline"
	"math/rand"
	"time"

	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/serialize/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	version     = "" // 游戏版本
	forceUpdate = false
	logger      = log.WithField("component", "game")
)

// Startup 初始化游戏服务器
func Startup() {
	rand.Seed(time.Now().Unix())
	version = viper.GetString("update.version")

	heartbeat := viper.GetInt("core.heartbeat")
	if heartbeat < 5 {
		heartbeat = 5
	}

	forceUpdate = viper.GetBool("update.force")

	logger.Infof("当前游戏服务器版本: %s, 是否强制更新: %t, 当前心跳时间间隔: %d秒", version, forceUpdate, heartbeat)
	logger.Info("game service starup")

	// register game handler
	comps := &component.Components{}
	comps.Register(defaultComonents)
	//comps.Register(defaultDeskManager)
	//comps.Register(new(ClubManager))

	// 加密管道
	c := newCrypto()
	pip := pipeline.New()
	pip.Inbound().PushBack(c.inbound)
	pip.Outbound().PushBack(c.outbound)

	addr := fmt.Sprintf("%s:%d", viper.GetString("gameserver.host"), viper.GetInt("gameserver.port"))
	fmt.Println(addr)
	nano.Listen(addr,
		nano.WithPipeline(pip),
		nano.WithHeartbeatInterval(time.Duration(heartbeat)*time.Second),
		nano.WithLogger(log.WithField("component", "nano")),
		nano.WithSerializer(json.NewSerializer()),
		nano.WithComponents(comps),
	)
}

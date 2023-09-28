package vars

const (
	PLAYER = "player"
)

type DeskStatus int32

const (
	//创建桌子
	DeskStatusCreate DeskStatus = iota
	//游戏
	DeskStatusPlaying
	DeskStatusRoundOver
	//游戏终/中止
	DeskStatusInterruption
	//已销毁
	DeskStatusDestory
	//已经清洗,即为下一轮准备好
	DeskStatusCleaned
)

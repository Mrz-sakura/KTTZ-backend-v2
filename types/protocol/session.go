package protocol

type LoginResponse struct {
	Uid      int64  `json:"uid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
	Gold     int    `json:"gold"`
}

type LoginRequest struct {
	Uid      int64  `json:"uid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"` //[0]未设置 [1]男 [2]女
	Gold     int    `json:"gold"`
	IP       string `json:"ip"`
}

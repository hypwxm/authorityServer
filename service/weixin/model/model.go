package model

// 小程序获取用户信息的原始数据
type MiniProgramUserAuth struct {
	Code          string `json:"code"`
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
	AppId         string `json:"appId"`
	AppSecret     string `json:"appSecret"`
	SessionKey    string
}

type MiniAuthWatermark struct {
	Appid     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}

// 解析小程序敏感信息返回的数据
type MiniProgramUserInfo struct {
	OpenId    string            `json:"openId"`
	NickName  string            `json:"nickName"`
	Gender    int               `json:"gender"`
	City      string            `json:"city"`
	Province  string            `json:"province"`
	Country   string            `json:"country"`
	AvatarUrl string            `json:"avatarUrl"`
	UnionId   string            `json:"unionId"`
	Watermark MiniAuthWatermark `json:"watermark"`
}

// 调用code2session返回的数据
type MiniSessionKey struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"` // 用户在开放平台的唯一标识符，在满足 UnionID 下发条件的情况下会返回，详见 UnionID 机制说明。
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"` // 错误信息
}

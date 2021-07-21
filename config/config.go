package config

type ConfigDefine struct {
	ServerPort            string
	AuthControlMysql      string
	SystemControlMysql    string
	ServiceMysql          string
	ServiceLogisticsMysql string
	UPLOADERHOST          string
	EmailActiveLink       string
	Hostname              string
	GlobalLock            string
	RabbitMqUrl           string
	Pgsql                 string
}

const AppServerTokenKey = "G_SERVER_USER"
const AppLoginUserName = "WB_SERVER_USER_NAME"

var Config = ConfigDefine{}

func InitConfig(cfg ConfigDefine) {
	Config = cfg
}

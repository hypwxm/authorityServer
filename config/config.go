package config

type config struct {
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

const ENV_DEV = "development"
const ENV_PROD = "product"
const ENV_TEST = "test"

const Env = ENV_DEV

var Config = config{}

func init() {
	if Env == ENV_DEV {
		Config = dev
	} else if Env == ENV_TEST {
		Config = test
	} else if Env == ENV_PROD {
		Config = prod
	}
}

const AppUserTokenKey = "WB_APP_USER"
const AppServerTokenKey = "WB_SERVER_USER"

const AppLoginUserName = "WB_SERVER_USER_NAME"

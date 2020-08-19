package config

var dev = config{
	ServerPort:      ":6789",
	UPLOADERHOST:    "",
	EmailActiveLink: "http://localhost:8000/user/index.html#/mine/info",
	Hostname:        "http://localhost:8000",
	GlobalLock:      "memory",
	RabbitMqUrl:     "amqp://root:123456@localhost:5672/",
	Pgsql:           "host=47.96.29.83 port=5432 user=postgres password=123456 dbname=baby_growing sslmode=disable",
}

var test = config{
	ServerPort:         ":3009",
	AuthControlMysql:   "root:123456@tcp(47.96.29.83:3306)/casbin?charset=utf8&parseTime=True&loc=Local",
	SystemControlMysql: "root:123456@tcp(47.96.29.83:3306)/casbin?charset=utf8&parseTime=True&loc=Local",
	ServiceMysql:       "root:123456@tcp(47.96.29.83:3306)/community?charset=utf8&parseTime=True&loc=Local",
	UPLOADERHOST:       "",
	EmailActiveLink:    "http://beta.dropshippinglite.com/user/index.html#/mine/info",
	Hostname:           "http://beta.dropshippinglite.com",
	GlobalLock:         "memory",
}

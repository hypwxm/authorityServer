package config

var prod = config{
	ServerPort:         ":3009",
	AuthControlMysql:   "easyshop:Shopeasyfor!%$@tcp(rm-2ev6yei17di1ae2i6ao.mysql.rds.aliyuncs.com:3306)/casbin?charset=utf8&parseTime=True&loc=Local",
	SystemControlMysql: "easyshop:Shopeasyfor!%$@tcp(rm-2ev6yei17di1ae2i6ao.mysql.rds.aliyuncs.com:3306)/casbin?charset=utf8&parseTime=True&loc=Local",
	ServiceMysql:       "easyshop:Shopeasyfor!%$@tcp(rm-2ev6yei17di1ae2i6ao.mysql.rds.aliyuncs.com:3306)/community?charset=utf8&parseTime=True&loc=Local",
	UPLOADERHOST:       "",
	GlobalLock:         "memory",
	EmailActiveLink:    "http://app.dropshippinglite.com/user/index.html#/mine/info",
	Hostname:           "http://app.dropshippinglite.com",
}

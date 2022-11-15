package common

import "go-micro.dev/v4/config"

type MysqlConfig struct {
	Host     string `json:"host"`     //数据库地址
	User     string `json:"user"`     //用户名
	Password string `json:"password"` //密码
	Database string `json:"database"` //创建的数据库名
	Port     int64  `json:"port"`     //端口
}

func GetMysqlFromConsul(config config.Config, path ...string) (*MysqlConfig, error) {
	mysqlConfig := &MysqlConfig{}
	err := config.Get(path...).Scan(mysqlConfig)
	return mysqlConfig, err
}

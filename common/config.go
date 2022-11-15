package common

import (
	"github.com/go-micro/plugins/v4/config/source/consul"
	"go-micro.dev/v4/config"
	"strconv"
)

// GetConsulConfig 设置配置中心
func GetConsulConfig(host string, port int64, prefix string) (config.Config, error) {
	consulSource := consul.NewSource(
		// 设置配置中心地址
		consul.WithAddress(host+":"+strconv.FormatInt(port, 10)),
		// 设置前缀,不设置是默认/micro/config
		consul.WithPrefix(prefix),
		// 设置是否移除前缀,ture是表示可以不带前缀直接获取对应配置
		consul.StripPrefix(true),
	)
	// 初始化配置
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	// 加载配置
	err = config.Load(consulSource)
	return config, err
}

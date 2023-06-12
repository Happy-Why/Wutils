package setting

import (
	"flag"
	"fmt"
	"testing"
)

var (
	configPaths string // 配置文件路径
	ConfigName  string // private配置文件名
	configType  string // 配置文件类型
)

func setupFlag() {
	// 命令行参数绑定
	flag.StringVar(&ConfigName, "name", "config", "配置文件名")
	flag.StringVar(&configType, "type", "yaml", "配置文件类型")
	flag.StringVar(&configPaths, "path", "./", "指定要使用的配置文件路径,多个路径用逗号隔开")
	flag.Parse()
}

func TestNewSetting(t *testing.T) {
	setupFlag()
	type App struct {
		Port    string
		Version string
	}
	type Mysql struct {
		Host   string
		Port   string
		Dbname string
	}
	type Config struct {
		App
		Mysql
	}
	setting, err := NewSetting(ConfigName, configType, configPaths)
	if err != nil {
		fmt.Println("NewSetting failed, err:", err)
	}
	fmt.Printf("settings:%#v\n", setting)
	var config Config
	if err = setting.BindAll(&config); err != nil {
		fmt.Println("BindAll failed, err:", err)
	}
	fmt.Printf("config:%#v\n", config)

}

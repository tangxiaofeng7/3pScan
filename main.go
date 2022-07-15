package main

import (
	"3pScan/pkg/runner"
	"3pScan/pkg/scan"
	"os"
	"os/signal"

	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gookit/color"
)

func init() {
	config := glog.DefaultConfig()
	config.Path = "log"
	config.File = "{Y-m-d}.log"
	// config.Flags = glog.F_FILE_SHORT

	_ = glog.SetConfig(config)
}

func main() {

	runner.ShowBanner()

	color.Yellow.Println(gtime.Datetime(), "程序初始化...")
	options := runner.ParseOptions()
	color.Yellow.Println("本次扫描模式为:", options.Model)

	// 优雅退出
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		glog.Warningf("CTRL+C 结束任务\n")
		os.Exit(1)
	}()

	t := gtime.Now()

	// 判断扫描模式port
	if gstr.Equal(options.Model, "all") {

		// port扫描
		tempPort, err := scan.NewPortScan(options)
		if err != nil {
			glog.Errorf("无法创建port扫描: %s\n", err)
		}

		// httpx扫描
		_, err = scan.NewHttpxScan(tempPort)
		if err != nil {
			glog.Errorf("无法创httpx扫描: %s\n", err)
		}

	}

	end := gtime.Now().Sub(t)

	color.Magenta.Printf("任务结束,耗时: %s\n", (end))
}

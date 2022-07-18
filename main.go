package main

import (
	"3pScan/pkg/runner"
	"3pScan/pkg/scan"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gookit/color"
)

func main() {

	runner.ShowBanner()

	options := runner.ParseOptions()
	color.White.Println("扫描模式:", options.Model)

	t := gtime.Now()

	// 判断扫描模式port
	if gstr.Equal(options.Model, "all") {

		// port扫描
		tempPort, err := scan.NewPortScan(options)
		if err != nil {
			glog.Errorf("无法创建port扫描: %s\n", err)
		}

		// httpx扫描
		if !gfile.IsEmpty(tempPort) {
			_, err = scan.NewHttpxScan(tempPort)
			if err != nil {
				glog.Errorf("无法创httpx扫描: %s\n", err)
			}
		}
	}

	end := gtime.Now().Sub(t)

	color.Magenta.Printf("%s 任务结束,耗时: %s\n", gtime.Datetime(), (end))
}

package main

import (
	"github.com/tangxiaofeng7/3pScan/pkg/runner"
	"github.com/tangxiaofeng7/3pScan/pkg/scan"

	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gookit/color"
)

func main() {

	runner.ShowBanner()

	options := runner.ParseOptions()

	t := gtime.Now()

	_, err := scan.NewPortScan(options)
	if err != nil {
		glog.Errorf("无法创建port扫描: %s\n", err)
	}

	end := gtime.Now().Sub(t)

	color.Magenta.Printf("%s 任务结束,耗时: %s\n", gtime.Datetime(), (end))
}

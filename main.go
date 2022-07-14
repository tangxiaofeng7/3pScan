package main

import (
	"3pScan/pkg/runner"
	"context"
	"os"
	"os/signal"

	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/guid"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gookit/color"
	httpxRunner "github.com/projectdiscovery/httpx/runner"
)

func init() {
	config := glog.DefaultConfig()
	config.Path = "log"
	config.File = "{Y-m-d}.log"
	config.Flags = glog.F_FILE_SHORT

	_ = glog.SetConfig(config)
}

func main() {
	runner.ShowBanner()
	ctx := context.TODO()

	options := runner.ParseOptions()
	color.Blue.Println(gtime.Datetime(), " 初始化...")
	color.Yellow.Println("本次设置扫描模式为:", options.Model)

	// 优雅退出
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		glog.Warningf(ctx, "CTRL+C 结束任务\n")

		os.Exit(1)
	}()
	// 开始计时
	t := gtime.Now()

	// 判断是否存在扫描模式port
	if gstr.Contains(options.Model, "port") {
		color.Danger.Print(gtime.Datetime(), " 开始端口扫描,运行中...\n")
		tempfile, err := runner.NewRunner(options)
		if err != nil {
			glog.Errorf(ctx, "无法创建扫描: %s\n", err)
		}
		color.Danger.Println(gtime.Datetime(), " 端口扫描完成,结果保存在文件\n")
		color.Yellow.Println(tempfile, "\n")

		// 判断是否存在扫描模式httpx
		if gstr.Contains(options.Model, "httpx") {
			color.Danger.Print(gtime.Datetime(), " 开始httpx扫描,运行中...\n\n")
			temphttpxfile := "./temp/" + guid.S() + ".txt"
			httpxoptions := httpxRunner.Options{
				Methods:      "GET",
				InputFile:    tempfile,
				ExtractTitle: true, //返回title
				StatusCode:   true, //返回状态
				Timeout:      3,    //超时
				Output:       temphttpxfile,
			}

			httpxRunner, err := httpxRunner.New(&httpxoptions)

			if err != nil {
				glog.Error(ctx, err.Error())
			}

			defer httpxRunner.Close()
			httpxRunner.RunEnumeration()

			color.Danger.Println("\n", gtime.Datetime(), "httpx扫描完成,结果保存在文件\n")
			color.Yellow.Println(temphttpxfile, "\n")
		}
	}

	end := gtime.Now().Sub(t)

	color.Magenta.Printf("任务结束,耗时: %s\n", (end))
}

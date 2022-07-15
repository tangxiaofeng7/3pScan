package runner

import (
	"flag"
	"os"

	"github.com/gookit/color"
)

type Options struct {
	Model        string //扫描模式
	Host         string // 主机地址
	HostFile     string // 主机文件
	Ports        string // 端口
	PortsFile    string // 端口文件
	TopPorts     string // 常见端口
	ExcludePorts string // 排除端口
	Timeout      int    // 超时时间
	Tips         int    // 提示信息
	Proxy        string // Socks5代理
	Rate         int    // 扫描速率
	Output       string // 输出文件
	Version      bool   // 显示版本
}

const banner = `
██████╗ ██████╗ ███████╗ ██████╗ █████╗ ███╗   ██╗
╚════██╗██╔══██╗██╔════╝██╔════╝██╔══██╗████╗  ██║
 █████╔╝██████╔╝███████╗██║     ███████║██╔██╗ ██║
 ╚═══██╗██╔═══╝ ╚════██║██║     ██╔══██║██║╚██╗██║
██████╔╝██║     ███████║╚██████╗██║  ██║██║ ╚████║
╚═════╝ ╚═╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═══╝                                
`

const Version = `1.0`

func ShowBanner() {
	color.Magenta.Printf("%s", banner)
	color.Magenta.Print("\t\thttps://github.com/tangxiaofeng7\n\n")
}

func ParseOptions() *Options {
	options := &Options{}
	flag.StringVar(&options.Model, "mode", "all", "扫描模式,可选值:port、all")
	flag.StringVar(&options.Host, "h", "", "主机地址")
	flag.StringVar(&options.HostFile, "hf", "", "主机文件")
	flag.StringVar(&options.Ports, "p", "", "端口")
	flag.StringVar(&options.PortsFile, "pf", "", "端口文件")
	flag.StringVar(&options.TopPorts, "top", "", "常用端口,可选值:full,100,1000")
	flag.StringVar(&options.ExcludePorts, "exclude-ports", "", "排除端口")
	flag.IntVar(&options.Timeout, "t", 1, "超时时间,默认1秒")
	flag.IntVar(&options.Tips, "tips", 5, "端口扫描提示信息间隔,默认5秒")
	flag.StringVar(&options.Proxy, "proxy", "", "Socks5代理")
	flag.StringVar(&options.Output, "o", "", "输出文件")
	flag.BoolVar(&options.Version, "version", false, "显示版本")

	flag.Parse()

	// 输出版本
	if options.Version {
		color.Yellow.Printf("当前版本: %s\n", Version)
		os.Exit(0)
	}

	return options
}

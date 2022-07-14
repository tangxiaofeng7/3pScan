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
	flag.StringVar(&options.Model, "mode", "port", "扫描模式")
	flag.StringVar(&options.Host, "host", "", "主机地址")
	flag.StringVar(&options.HostFile, "hostfile", "", "主机文件")
	flag.StringVar(&options.Ports, "ports", "", "端口")
	flag.StringVar(&options.PortsFile, "ports-file", "", "端口文件")
	flag.StringVar(&options.TopPorts, "top-ports", "", "常用端口")
	flag.StringVar(&options.ExcludePorts, "exclude-ports", "", "排除端口")
	flag.IntVar(&options.Timeout, "timeout", 1, "超时时间,默认1秒")
	flag.StringVar(&options.Proxy, "proxy", "", "Socks5代理")
	flag.StringVar(&options.Output, "output", "", "输出文件")
	flag.BoolVar(&options.Version, "version", false, "显示版本")

	flag.Parse()

	// 输出版本
	if options.Version {
		color.Yellow.Printf("当前版本: %s\n", Version)
		os.Exit(0)
	}

	return options
}

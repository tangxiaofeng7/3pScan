package runner

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gookit/color"
	"golang.org/x/sync/semaphore"
)

type Result struct {
	Host string
	Port int
}

func NewRunner(options *Options) (string, error) {
	var Targets []string
	var res []Result
	var semMaxWeight int64 = 20000
	var semAcquisitionWeight int64 = 100

	sem := semaphore.NewWeighted(semMaxWeight)
	ctx := context.Background()

	// 识别终端提供的Ports参数
	Ports, err := ParsePorts(options)
	if err != nil {
		return "", fmt.Errorf("could not parse ports: %s", err)
	}

	if options.Host != "" {
		Targets = append(Targets, options.Host)
	}

	if options.HostFile != "" {
		TempFile := gfile.GetContents(options.HostFile)
		Targets = gstr.SplitAndTrim(TempFile, "\n")
	}

	glog.Infof(context.TODO(), "设置扫描的主机: %v", Targets)
	// glog.Infof(context.TODO(), "设置扫描端口: %v", Ports)

	for _, target := range Targets {
		for _, port := range Ports {
			// glog.Infof(context.TODO(), "开始主机: %s扫描端口: %d", target, port)

			if err := sem.Acquire(ctx, semAcquisitionWeight); err != nil {
				fmt.Printf("Failed to acquire semaphore (port %d): %v\n", port, err)
				break
			}

			go func(target string, port int) {
				defer sem.Release(semAcquisitionWeight)
				p := scan(target, port)
				if p != 0 {
					res = append(res, Result{target, port})

					// glog.Noticef(context.TODO(), "主机%s扫描到开放端口: %d", target, p)

				}
			}(target, port)
		}
	}
	if err := sem.Acquire(ctx, int64(semMaxWeight)); err != nil {
		fmt.Printf("Failed to acquire semaphore: %v\n", err)
	}

	tempfile := printResults(res)

	return tempfile, nil
}

func scan(host string, port int) int {
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", address, time.Duration(1)*time.Second)
	if err != nil {
		// fmt.Printf("%d CLOSED (%s)\n", port, err)
		return 0
	}
	conn.Close()
	return port
}

// proxy暂时不用
// func scanWithProxy(host, Proxy string, port int) int {
// 	address := fmt.Sprintf("%s:%d", host, port)
// 	dialer, err := proxy.SOCKS5("tcp", Proxy, nil, &net.Dialer{Timeout: time.Duration(1) * time.Second})
// 	if err != nil {
// 		glog.Errorf(context.TODO(), "代理连接异常 (%s)\n", err)
// 	}
// 	conn, err := dialer.Dial("tcp", address)
// 	if err != nil {
// 		fmt.Printf("%d CLOSED (%s)\n", port, err)
// 		return 0
// 	}
// 	conn.Close()
// 	return port
// }

func printResults(res []Result) (tempfile string) {
	color.Green.Println("\nResults\n--------------")
	tempfile = "./temp/" + guid.S() + ".txt"
	for _, b := range res {
		color.Green.Println(b.Host, "\t", b.Port, "\topen")

		port := strconv.Itoa(b.Port)

		gfile.PutContentsAppend(tempfile, b.Host+":"+port+"\n")
	}
	color.Blue.Println("\n")
	return tempfile
}

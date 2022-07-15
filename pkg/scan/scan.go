package scan

import (
	"3pScan/pkg/runner"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"github.com/projectdiscovery/clistats"
	"github.com/projectdiscovery/ipranger"
	"golang.org/x/sync/semaphore"

	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/guid"
	"github.com/gookit/color"
	httpxRunner "github.com/projectdiscovery/httpx/runner"
)

type Result struct {
	Host string
	Port int
}

func NewPortScan(options *runner.Options) (string, error) {
	color.Red.Print(gtime.Datetime(), " --------------------开始端口扫描--------------------\n\n")

	var Targets []string
	var res []Result
	var semMaxWeight int64 = 20000
	var semAcquisitionWeight int64 = 100

	sem := semaphore.NewWeighted(semMaxWeight)
	ctx := context.Background()

	// 识别终端提供的Ports参数
	Ports, err := runner.ParsePorts(options)
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

	// 对用户传入的Targets进行处理
	Targets = HandleIP(Targets)

	Range := len(Targets) * len(Ports)
	color.Yellow.Println("扫描的主机数:", len(Targets), "扫描端口数:", Range)

	// 进度条
	stats, err := clistats.New()

	if err != nil {
		glog.Warningf("Couldn't create progress engine: %s\n", err)
	}

	stats.AddStatic("ports", len(Ports))
	stats.AddStatic("hosts", len(Targets))

	stats.AddStatic("startedAt", time.Now())
	stats.AddCounter("packets", uint64(0))
	stats.AddCounter("total", uint64(Range))

	if err := stats.Start(makePrintCallback(), time.Duration(options.Tips)*time.Second); err != nil {
		glog.Warningf("Couldn't start statistics: %s\n", err)
	}

	for _, target := range Targets {
		for _, port := range Ports {
			stats.IncrementCounter("packets", 1)

			if err := sem.Acquire(ctx, semAcquisitionWeight); err != nil {
				fmt.Printf("Failed to acquire semaphore (port %d): %v\n", port, err)
				break
			}

			go func(target string, port int) {
				defer sem.Release(semAcquisitionWeight)
				p := scan(target, port)
				if p != 0 {
					res = append(res, Result{target, port})
					color.Green.Println(target, p, " open")

				}
			}(target, port)
		}
	}
	if err := sem.Acquire(ctx, int64(semMaxWeight)); err != nil {
		fmt.Printf("Failed to acquire semaphore: %v\n", err)
	}

	tempfile := PrintResults(res)
	stats.Stop()

	if !gfile.IsEmpty(tempfile) {
		color.Red.Println(gtime.Datetime(), " 端口扫描完成,结果保存在文件 ", tempfile, "\n")
	} else {
		color.Red.Println(gtime.Datetime(), " 端口扫描完成,没有任何结果 ", "\n")

	}

	return tempfile, nil
}

func NewHttpxScan(tempfile string) (string, error) {

	color.Red.Print(gtime.Datetime(), " --------------------开始httpx扫描--------------------\n\n")
	temphttpxfile := "./temp/" + guid.S() + ".csv"
	httpxoptions := httpxRunner.Options{
		Methods:            "GET",
		InputFile:          tempfile,
		ExtractTitle:       true, //返回title
		StatusCode:         true, //返回状态
		Timeout:            3,    //超时
		OutputResponseTime: true, //返回响应时间
		OutputServerHeader: true, //返回服务器头
		Probe:              true, //返回探针
		Output:             temphttpxfile,
	}

	httpxRunner, err := httpxRunner.New(&httpxoptions)

	if err != nil {
		glog.Error(context.TODO(), err.Error())
	}

	defer httpxRunner.Close()
	httpxRunner.RunEnumeration()

	color.Red.Println("\n", gtime.Datetime(), "httpx扫描完成,结果保存在文件", temphttpxfile, "\n")

	return temphttpxfile, nil
}

func scan(host string, port int) int {
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", address, time.Duration(1)*time.Second)
	if err != nil {
		return 0
	}
	conn.Close()
	return port
}

func PrintResults(res []Result) (tempfile string) {
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

// 处理IP段
func HandleIP(Target []string) (out []string) {
	color.Red.Println(gtime.Datetime(), "HandleIp Starting......")
	var ipTemp garray.StrArray
	for _, target := range Target {
		if _, _, err := net.ParseCIDR(target); err == nil { // 判断是否为ip段
			a, _ := Hosts(target)
			for _, b := range a {
				ipTemp.Append(b)
				out = append(out, b)
			}
		} else if ipranger.IsIP(target) && !ipTemp.Contains(target) { // 判断是否为ip
			ipTemp.Append(target)
			out = append(out, target)
		} else if _, err := net.LookupIP(target); err == nil { // 判断是否为域名
			out = append(out, target)
		}

	}

	color.Red.Println(gtime.Datetime(), "HandleIp Over......")

	return out
}

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// 进度条
const bufferSize = 128

func makePrintCallback() func(stats clistats.StatisticsClient) {
	builder := &strings.Builder{}
	builder.Grow(bufferSize)

	return func(stats clistats.StatisticsClient) {
		builder.WriteRune('[')
		startedAt, _ := stats.GetStatic("startedAt")
		duration := time.Since(startedAt.(time.Time))
		builder.WriteString(clistats.FmtDuration(duration))
		builder.WriteRune(']')

		hosts, _ := stats.GetStatic("hosts")
		builder.WriteString(" | 扫描主机: ")
		builder.WriteString(clistats.String(hosts))

		ports, _ := stats.GetStatic("ports")
		builder.WriteString(" | 扫描端口: ")
		builder.WriteString(clistats.String(ports))

		packets, _ := stats.GetCounter("packets")
		total, _ := stats.GetCounter("total")

		builder.WriteString(" | 端口扫描进度: ")
		builder.WriteString(clistats.String(packets))
		builder.WriteRune('/')
		builder.WriteString(clistats.String(total))
		builder.WriteRune(' ')
		builder.WriteRune('(')
		//nolint:gomnd // this is not a magic number
		builder.WriteString(clistats.String(uint64(float64(packets) / float64(total) * 100.0)))
		builder.WriteRune('%')
		builder.WriteRune(')')
		builder.WriteRune('\n')

		fmt.Fprintf(os.Stderr, "%s", builder.String())
		builder.Reset()
	}
}

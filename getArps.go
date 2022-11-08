package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ARPResult struct {
	IP              string
	MAC             string
	HostName        string
	HostNameAliases string
	Color           string
	OnOff           bool
	Comment         string
}

func getLocalArp() []*ARPResult {
	fd, err := os.Open("/proc/net/arp")
	if err != nil {
		fmt.Printf("读arp出错: %s\n", err)
		return nil
	}
	defer fd.Close()

	arps := make([]*ARPResult, 0, 10)
	buff := bufio.NewReader(fd)
	for {
		data, _, eof := buff.ReadLine()
		if eof != nil {
			break
		}
		if arp := parseOneLocalArp(data); arp != nil {
			arps = append(arps, arp)
		}
	}
	return mergeArps(arps, 0)
}

func getWisfArp(in []*ARPResult) ([]*ARPResult, []*ARPResult) {
	host := getHttpData("CMD=uploadfw&pnmode=%3B&POSTURL=web 2860 wifi getarptable;")
	if host == nil {
		return nil, nil
	}
	arps1 := parseWifiHost(in, host)

	clients := getHttpData("CMD=uploadfw&pnmode=%3B&POSTURL=/etc_ro/lighttpd/www/cgi-bin/accom.cgi LIST;")
	if clients == nil {
		return arps1, nil
	}
	arps2 := parseWifiClients(in, clients)
	return arps1, arps2
}

func getHttpData(body string) []byte {
	url := "http://192.168.1.30/cgi-bin/accom.cgi"
	DefaultClient := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(1 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*1)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		fmt.Printf("创建http POST请求失败: %s\n", err)
		return nil
	}
	req.Header.Set("content-type", "text/plain")
	req.Header.Set("charset", "UTF-8")
	resp, err := DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("页面获取请求失败: %s\n", err)
		return nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取页面数据失败: %s\n", err)
		return nil
	}
	return data
}

func parseWifiHost(in []*ARPResult, data []byte) []*ARPResult {
	arps := make([]*ARPResult, 0, 4)
	if len(data) < 1 {
		return arps
	}
	for _, each := range strings.Split(string(data), "#") {
		v := strings.Split(string(each), ";")
		if len(v) < 2 {
			continue
		}
		arps = append(arps, &ARPResult{IP: v[0], MAC: v[1]})
	}
	fillinIP(in, arps)
	return mergeArps(arps, 0)
}
func parseWifiClients(in []*ARPResult, data []byte) []*ARPResult {
	arps := make([][]*ARPResult, 0, 5)
	if len(data) < 1 {
		return nil
	}
	for _, each := range strings.Split(string(data), "#") {
		arps_tmp := make([]*ARPResult, 0, 10)
		v := strings.Split(string(each), ";")
		if len(v) < 3 {
			continue
		}
		arps_tmp = append(arps_tmp, &ARPResult{IP: v[1], MAC: v[2]})
		l, _ := strconv.Atoi(v[3])
		for _, mac := range v[4 : 4+l] {
			arps_tmp = append(arps_tmp, &ARPResult{MAC: mac})
		}
		fillinIP(in, arps_tmp)
		arps_tmp = mergeArps(arps_tmp, 1)
		arps = append(arps, arps_tmp)
	}
	return getarps(arps)
}
func getarps(arpss [][]*ARPResult) []*ARPResult {
	sort.SliceStable(arpss, func(i int, j int) bool {
		if len(arpss[i]) == 0 {
			return false
		} else if len(arpss[j]) == 0 {
			return true
		} else {
			return arpss[i][0].IP < arpss[j][0].IP
		}
	})
	res := make([]*ARPResult, 0, 20)
	for _, arps := range arpss {
		arps[0].Color = "\033[4m"
		res = append(res, arps...)
	}
	return res
}
func fillinIP(in, out []*ARPResult) {
	for _, e := range out {
		if e.IP == "" {
			for _, i := range in {
				if i.MAC == e.MAC {
					e.IP = i.IP
				}
			}
		}
	}
}

func parseOneLocalArp(line []byte) *ARPResult {
	var a, b, c, d, e, f, g string
	n, err := fmt.Sscanf(string(bytes.TrimLeft(line, "  	")), "%s %s %s %s %s %s %s", &a, &b, &c, &d, &e, &f, &g)
	if n < 5 || (err != nil && err != io.EOF) || g != "" || d == "00:00:00:00:00:00" {
		return nil
	}
	return &ARPResult{IP: a, MAC: d}
}

func MergeAllResult(arps1, arps21, arps22 []*ARPResult) []*ARPResult {
	arps1 = append(arps1, arps21...)
	arps1 = mergeArps(arps1, 0)
	// 先把wisf的ip尽量补全
	for _, arp2 := range arps22 {
		for _, arp1 := range arps1 {
			if arp1.MAC == arp2.MAC {
				arp1.OnOff = true //arp1中标记不再显示
				if !strings.Contains(arp2.IP, arp1.IP) {
					if arp2.IP == "" {
						arp2.IP = arp1.IP
					} else {
						if arp2.IP < arp1.IP {
							arp2.IP = arp2.IP + "\n" + arp1.IP
						} else {
							arp2.IP = arp1.IP + "\n" + arp2.IP
						}
					}
				}
			}
		}
	}

	// 以arp2为基础，去除重复项
	arps := make([]*ARPResult, 0, len(arps1)+len(arps22))
	for _, arp1 := range arps1 {
		if !arp1.OnOff {
			arps = append(arps, arp1)
		}
	}
	for _, arp2 := range arps22 {
		arps = append(arps, arp2)
	}
	return arps
}

// 排序单个组：各ap组，非wifi组
func mergeArps(arps []*ARPResult, sortstart int) []*ARPResult {
	if len(arps) < sortstart {
		return arps
	}
	sort.SliceStable(arps, func(i int, j int) bool {
		if sortstart > i { //第一个节点不排序
			return "0" < arps[i].MAC
		} else if sortstart > j {
			return arps[i].MAC < "0"
		}
		return arps[i].MAC < arps[j].MAC
	})
	res := make([]*ARPResult, 0, len(arps))
	var old_arp *ARPResult
	for _, arp := range arps {
		if old_arp == nil || old_arp.MAC != arp.MAC {
			res = append(res, arp)
			old_arp = arp
		} else {
			if !strings.Contains(res[len(res)-1].IP, arp.IP) {
				res[len(res)-1].IP += "\n" + arp.IP
			}
		}
	}
	sort.SliceStable(res, func(i int, j int) bool {
		if sortstart > i { //第一个节点不排序
			return "0" < res[j].IP
		} else if sortstart > j {
			return res[i].IP < "0"
		}
		if len(res[i].IP) != len(res[j].IP) {
			return len(res[i].IP) < len(res[j].IP)
		}
		return res[i].IP < res[j].IP
	})
	return res
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var formatFile string = "./devicelists.dat"
var hostsFile string = "/etc/hosts"

type deviceItem struct {
	MAC             string
	HostNameAliases string
}

func formatArps(arps []*ARPResult) {
	devLists := readDeviceListFile(formatFile)
	for _, arp := range arps {
		for _, devitem := range devLists {
			if arp.MAC == devitem.MAC {
				arp.HostNameAliases = devitem.HostNameAliases
				break
			}
		}
	}
	//checkHostFile(hostsFile)
}

var deviceFileReaded bool
var devLists []*deviceItem

func checkHostFile(filename string) {
	readHostsFile(filename)
}

func readDeviceListFile(filename string) []*deviceItem {
	if deviceFileReaded {
		return devLists
	}
	deviceFileReaded = true

	fd, err := os.Open(filename)
	if err != nil {
		fmt.Printf("读devicelists.dat出错: %s\n", err)
		return nil
	}
	defer fd.Close()
	devLists = make([]*deviceItem, 0, 20)
	buff := bufio.NewReader(fd)
	for {
		data, _, eof := buff.ReadLine()
		if eof != nil {
			break
		}

		if dev := parseOneDeviceItem(data); dev != nil {
			devLists = append(devLists, dev)
		}
	}
	return devLists
}
func parseOneDeviceItem(line []byte) *deviceItem {
	l := strings.TrimLeft(string(line), "  	")
	if strings.HasPrefix(l, "#") {
		return nil
	}
	var a, b string
	n, err := fmt.Sscanf(l, "%s %s", &a, &b)
	if n < 2 || (err != nil && err != io.EOF) {
		return nil
	}

	return &deviceItem{MAC: a, HostNameAliases: b}
}

type hostItem struct {
	IP       string
	HostName string
	MAC      string
}

func readHostsFile(filename string) []*hostItem {
	fd, err := os.Open(filename)
	if err != nil {
		fmt.Printf("读%s出错: %s\n", filename, err)
		return nil
	}
	defer fd.Close()
	hostLists := make([]*hostItem, 0, 20)
	buff := bufio.NewReader(fd)
	for {
		data, _, eof := buff.ReadLine()
		if eof != nil {
			break
		}

		if host := parseOneHostItem(data); host != nil {
			hostLists = append(hostLists, host)
		}
	}
	return hostLists
}
func parseOneHostItem(line []byte) *hostItem {
	l := strings.TrimLeft(string(line), "  	")
	if strings.HasPrefix(l, "#") {
		return nil
	}
	var a, b, c, d string
	n, err := fmt.Sscanf(l, "%s %s %s %s", &a, &b, &c, &d)
	if n < 2 || (err != nil && err != io.EOF) {
		return nil
	}

	d = formatMAC(d)
	return &hostItem{IP: a, HostName: b, MAC: d}
}
func formatMAC(mac string) string {
	if len(mac) == 17 {
		return mac
	}
	vs := strings.Split(mac, ":")
	if len(vs) != 6 {
		return ""
	}
	for i, v := range vs {
		for len(v) < 2 {
			v = "0" + v
		}
		vs[i] = v
	}
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s", vs[0], vs[1], vs[2], vs[3], vs[4], vs[5])
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestGetLocalArp(t *testing.T) {
	arps := getLocalArp()
	if len(arps) == 0 {
		t.Errorf("Can't get arps from /proc/net/arp")
	}
	for _, each := range arps {
		fmt.Printf("GetLocalArp: %v\n", *each)
	}
}

func TestparseLocalArps() []*ARPResult {
	data := `
IP address       HW type     Flags       HW address            Mask     Device
192.168.1.31     0x1		 0x2   		 c0:4a:09:27:81:46     *        switch0
192.168.1.50     0x1		 0x2   		 a4:4b:d5:10:74:93     *        switch0
192.168.1.8      0x1		 0x2   		 0c:1d:af:c9:6c:95     *        switch0
192.168.1.204    0x1		 0x2   		 00:0c:e8:c3:6e:6c     *        switch0
192.168.1.52     0x1		 0x2   		 04:fa:83:f3:62:5b     *        switch0
192.168.1.132    0x1		 0x2   		 24:df:a7:9b:73:06     *        switch0
192.168.1.23     0x1		 0x2   		 68:f7:28:04:a6:77     *        switch0
192.168.1.30     0x1		 0x2   		 c0:4a:09:26:46:d8     *        switch0
192.168.1.7      0x1		 0x0   		 68:f7:28:04:a6:77     *        switch0
192.168.1.51     0x1		 0x2   		 04:cf:8c:46:c3:ab     *        switch0
192.168.1.9      0x1		 0x2   		 a8:9c:ed:8c:a8:0a     *        switch0
192.168.1.64     0x1		 0x2   		 44:23:7c:54:15:27     *        switch0
192.168.1.22     0x1		 0x2   		 00:00:00:00:00:00     *        switch0
192.168.1.11     0x1		 0x2   		 00:00:00:00:00:00     *        switch0
`
	arps := make([]*ARPResult, 0, 10)
	buff := bufio.NewReader(bytes.NewReader([]byte(data)))
	for {
		data, _, eof := buff.ReadLine()
		if eof != nil {
			break
		}
		if arp := parseOneLocalArp(data); arp != nil {
			arps = append(arps, arp)
		}
	}
	return mergeArps(arps)
}

func TestParseLocalArps(t *testing.T) {
	arps := TestparseLocalArps()
	if len(arps) == 0 {
		t.Errorf("Can't ParseLocalArps")
	}
	for _, each := range arps {
		fmt.Printf("ParseLocalArps: %v\n", *each)
	}
}

func TestGetHttpData1(t *testing.T) {
	host := getHttpData("CMD=uploadfw&pnmode=%3B&POSTURL=web 2860 wifi getarptable;")
	if host == nil || len(host) < 12 {
		t.Skipf("Can't GetHttpData1")
	}
}
func TestGetHttpData2(t *testing.T) {
	clients := getHttpData("CMD=uploadfw&pnmode=%3B&POSTURL=web 2860 wifi getarptable;")
	if clients == nil || len(clients) < 12 {
		t.Skipf("Can't GetHttpData2")
	}
}

func TestparseWifiHost() []*ARPResult {
	return parseWifiHost([]byte("192.168.1.1;74:ac:b9:a2:eb:91#192.168.1.9;a8:9c:ed:8c:a8:0a#"))
}
func TestParseWifiHost(t *testing.T) {
	arps := TestparseWifiHost()
	for _, each := range arps {
		fmt.Printf("ParseWifiHost: %v\n", *each)
	}
	if len(arps) != 2 {
		t.Errorf("ParseWifiClients 数量应该为2, 得到: %d", len(arps))
	}
}

func TestparseWifiClients() []*ARPResult {
	return parseWifiClients([]byte("AP4_老人房;192.168.1.34;c0:4a:09:29:3d:00;0;0;#AP1_二楼;192.168.1.31;c0:4a:09:27:81:46;4;44:23:7c:54:15:27;24:df:a7:9b:73:06;a4:4b:d5:10:74:93;a8:9c:ed:8c:a8:0a;0;#AP2_餐厅;192.168.1.32;c0:4a:09:27:83:46;2;04:fa:83:f3:62:5b;d0:ff:98:59:f0:2d;0;#AP3_客厅;192.168.1.33;c0:4a:09:27:7c:f2;4;62:18:c4:12:c8:a4;0c:1d:af:c9:6c:95;00:0c:e8:c3:6e:6c;04:cf:8c:46:c3:ab;0;#"))
}
func TestParseWifiClients(t *testing.T) {
	arps := TestparseWifiClients()
	for _, each := range arps {
		fmt.Printf("ParseWifiClients: %v\n", *each)
	}
	if len(arps) != 14 {
		t.Errorf("ParseWifiClients 数量应该为14, 得到: %d", len(arps))
	}
}

func TestmergeAllResult() []*ARPResult {
	arps1 := TestparseLocalArps()
	arps21 := TestparseWifiHost()
	arps22 := TestparseWifiClients()
	return MergeAllResult(arps1, arps21, arps22)
}
func TestMergeAllResult(t *testing.T) {
	arps := TestmergeAllResult()
	for _, each := range arps {
		fmt.Printf("MergeAllResult: %v\n", *each)
	}
	if len(arps) != 17 {
		t.Errorf("MergeAllResult 数量应该为17, 得到: %d", len(arps))
	}
}

package main

import (
	"fmt"
	"testing"
)

func TestReadDeviceListFile(t *testing.T) {
	lists := readDeviceListFile("testdata/devicelists.dat")
	for _, each := range lists {
		fmt.Printf("ReadDeviceListFile: %v\n", *each)
	}
	if len(lists) != 16 {
		t.Errorf("ReadDeviceListFile 数量应该为16, 得到: %d", len(lists))
	}
}

func TestReadHostsFile(t *testing.T) {
	lists := readHostsFile("testdata/hosts")
	for _, each := range lists {
		fmt.Printf("ReadHostsFile: %v\n", *each)
	}
	if len(lists) != 25 {
		t.Errorf("ReadHostsFile 数量应该为25, 得到: %d", len(lists))
	}
}

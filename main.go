package main

import (
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	initPWD()
	initInputOutput()
	runtime.LockOSThread()
	for {
		start := time.Now()
		MyPrintf("%s\n", start.Format("2006-01-02 15:04:05.999"))

		arp1 := getLocalArp()
		arp21, arp22 := getWisfArp(arp1)
		arps := MergeAllResult(arp1, arp21, arp22)
		formatArps(arps)
		printARP(arps)
		printLink()

		printConsole()
		elapsed := time.Since(start)
		if elapsed > time.Second*3 {
			elapsed = time.Second*3 - time.Millisecond*100
		}
		select {
		case v := <-inputChain:
			for {
				switch v {
				case 'q', 'e', 'x':
					return
				}
				if len(inputChain) > 0 {
					v = <-inputChain
				} else {
					break
				}
			}
		case <-time.After(time.Second*3 - elapsed):
		}
	}
}

func initPWD() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	os.Chdir(dir)
}

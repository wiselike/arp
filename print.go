package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

var cmds []byte = make([]byte, 0, 12)

func printARP(arps []*ARPResult) {
	MyPrintf("IP\t\tMAC\t\t\t设备名\n")
	for _, arp := range arps {
		if arp.HostNameAliases == "" {
			MyPrintf("\033[0;31m%v\t%v\033[0m\n", arp.IP, arp.MAC)
		} else if arp.IP == "" {
			MyPrintf("\033[0;35m?\t\t%v\t%v\033[0m\n", arp.MAC, arp.HostNameAliases)
		} else if arp.Color != "" {
			MyPrintf(arp.Color+"%v    %v       %v\033[0m\n", arp.IP, arp.MAC, arp.HostNameAliases)
		} else {
			MyPrintf("%v\t%v\t%v\n", arp.IP, arp.MAC, arp.HostNameAliases)
		}
	}
}

func printLink() {
	for i := 0; i < 5; i++ {
		if v, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/net/eth%d/carrier", i)); err == nil {
			if v[0] == '0' {
				MyPrintf("Eth%d\t\t%s\n", i, "DOWN")
				continue
			}
			if v, err = ioutil.ReadFile(fmt.Sprintf("/sys/class/net/eth%d/speed", i)); err == nil {
				MyPrintf("Eth%d\t\t%s\n", i, bytes.TrimSpace(v))
				continue
			}
		}
		MyPrintf("Eth%d\t\t%s\n", i, "Unknown")
	}
}

func printConsole() {
	fmt.Printf("\033[H\033[2J\033[0;J")
	fmt.Printf("\033[H\033[2J\033[3;J")
	Output.WriteTo(os.Stdout)
	Output.Reset()
	//fmt.Fprintf(Stdout, "\033[?25l")
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWidth() *winsize {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		fmt.Printf("无法获取终端大小: %v\n", errno)
		return nil
	}
	return ws
}

var Output *bytes.Buffer
var Stdout *os.File
var inputChain chan byte = make(chan byte, 10)

func initInputOutput() {
	Output = new(bytes.Buffer)

	go func() {
		var ch byte
		var err error
		runtime.LockOSThread()
		readinit()
		for {
			ch, err = readch()
			if err != nil {
				return
			}
			inputChain <- ch
			if ch == '\n' {
				cmds = cmds[:0]
				continue
			}

			if len(cmds) > 10 {
				cmds = format(cmds)
			}
			for i, c := range cmds {
				if c == ch {
					cmds[i] = 0
					break
				}
			}
			cmds = append(cmds, ch)
		}
	}()
}

func format(in []byte) []byte {
	newCmds := make([]byte, 0, len(in))
	for _, c := range in {
		if c != 0 {
			newCmds = append(newCmds, c)
		}
	}
	return newCmds
}

func MyPrintf(format string, a ...interface{}) {
	fmt.Fprintf(Output, format, a...)
}

var termiosOld syscall.Termios

func readinit() {
	syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termiosOld)))
	termiosOld.Lflag &= ^uint32(syscall.ICANON)
	//termiosOld.Lflag &= ^uint32(syscall.ECHO)
	termiosOld.Cc[syscall.VMIN] = 1
	termiosOld.Cc[syscall.VTIME] = 0
}
func readch() (byte, error) {
	syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termiosOld)))
	tmp := make([]byte, 1)
	_, err := os.Stdin.Read(tmp)
	return tmp[0], err
}

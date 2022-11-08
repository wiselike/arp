package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var termiosOld syscall.Termios

func readinit() {
	syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termiosOld)))
	termiosOld.Lflag &= ^uint32(syscall.ICANON)
	termiosOld.Lflag &= ^uint32(syscall.ECHO)
	termiosOld.Cc[6] = 1
	termiosOld.Cc[5] = 0
}
func readch() (byte, error) {
	syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termiosOld)))
	tmp := make([]byte, 1)
	_, err := os.Stdin.Read(tmp)
	return tmp[0], err
}

func main() {
	readinit()
	for {
		ch, err := readch()
		if err != nil {
			fmt.Printf("err!\n")
		}
		fmt.Printf("%v", ch)
	}
}

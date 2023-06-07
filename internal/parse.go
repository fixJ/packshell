package internal

import (
	"bytes"
	"compress/gzip"
	"github.com/winlabs/gowin32"
	"golang.org/x/sys/windows"
	"io/ioutil"
	"syscall"
)

func FindResource(module windows.Handle, name uintptr, resType uintptr) (resInfo windows.Handle, err error) {
	r0, _, e1 := syscall.SyscallN(windows.NewLazySystemDLL("kernel32.dll").NewProc("FindResourceW").Addr(), uintptr(module), uintptr(name), uintptr(resType))
	resInfo = windows.Handle(r0)
	if resInfo == 0 {
		err = e1
	}
	return
}

func LoadResource(module windows.Handle, resInfo windows.Handle) (resData windows.Handle, err error) {
	r0, _, e1 := syscall.SyscallN(windows.NewLazySystemDLL("kernel32.dll").NewProc("LoadResource").Addr(), uintptr(module), uintptr(resInfo), 0)
	resData = windows.Handle(r0)
	if resData == 0 {
		err = e1
	}
	return
}

func LockResource(resData windows.Handle) (addr uintptr, err error) {
	r0, _, e1 := syscall.SyscallN(windows.NewLazyDLL("kernel32.dll").NewProc("LockResource").Addr(), uintptr(resData), 0, 0)
	addr = uintptr(r0)
	if addr == 0 {
		err = e1
	}
	return
}

func GetResource(module windows.Handle, name int, typ gowin32.ResourceType) (uintptr, error) {
	FileDataLenRes, err := FindResource(0, uintptr(gowin32.IntResourceId(uint(name))), uintptr(typ))
	if err != nil {
		return 0, err
	}
	FileDataLenResData, err := LoadResource(0, FileDataLenRes)
	if err != nil {
		return 0, err
	}
	pFileDataLen, err := LockResource(FileDataLenResData)
	if err != nil {
		return 0, err
	}
	return pFileDataLen, nil
}

func DecompressFile(data []byte) ([]byte, error) {
	rdata := bytes.NewReader(data)
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return nil, err
	}
	s, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func DecodeFile(data []byte, xorval byte) []byte {
	res := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		res[i] = data[i] ^ xorval
	}
	return res
}

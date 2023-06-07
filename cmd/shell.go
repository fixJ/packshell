package main

import (
	"encoding/binary"
	"fmt"
	"github.com/winlabs/gowin32"
	"packshell/internal"
	"unsafe"
)

var (
	rsrcFileDataLenID = 10
	rsrcFileDataID    = 20
	rsrcFileKeyID     = 30
)

func main() {
	// 1.get data from rsrc section
	fmt.Printf("[*] Starting get source file from rsrc section\n")
	pFileDataLen, err := internal.GetResource(0, rsrcFileDataLenID, gowin32.ResourceTypeRCData)
	if err != nil {
		internal.CheckErr(err)
		return
	}
	pFileData, err := internal.GetResource(0, rsrcFileDataID, gowin32.ResourceTypeRCData)
	if err != nil {
		internal.CheckErr(err)
		return
	}
	pFileKey, err := internal.GetResource(0, rsrcFileKeyID, gowin32.ResourceTypeRCData)
	if err != nil {
		internal.CheckErr(err)
		return
	}
	fmt.Printf("[*] Get rsrc section successful\n")
	// 2.get data
	fileDataLenBytes := make([]byte, 4)
	var i uintptr
	var j uint32
	for i = 0; i < 4; i++ {
		fileDataLenBytes[i] = *(*byte)(unsafe.Pointer(pFileDataLen + i))
	}
	fileDataLen := binary.LittleEndian.Uint32(fileDataLenBytes)

	fileDataBytes := make([]byte, fileDataLen)
	for j = 0; j < fileDataLen; j++ {
		fileDataBytes[j] = *(*byte)(unsafe.Pointer(pFileData + uintptr(j)))
	}
	var fileKey byte
	fileKey = *(*byte)(unsafe.Pointer(pFileKey))

	// 3.decompress and decode file
	fmt.Printf("[*] Starting decompress and decode source file\n")

	decompressFile, err := internal.DecompressFile(fileDataBytes)
	if err != nil {
		internal.CheckErr(err)
		return
	}
	originFile := internal.DecodeFile(decompressFile, fileKey)
	fmt.Printf("[*] Decompress and decode source file successful\n")

	// 4. run in memory
	internal.Run(originFile)
}

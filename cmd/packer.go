package main

import (
	"encoding/binary"
	"fmt"
	"github.com/winlabs/gowin32"
	"os"
	"packshell/internal"
)

var (
	fileDataLenID = 10
	fileDataID    = 20
	fileKeyID     = 30
)

func main() {

	args := os.Args
	if len(args) != 3 {
		fmt.Printf("[X] Need two args, first is input file second is output file\n")
		return
	}
	in := args[1]
	out := args[2]
	shellFile := "shell.exe"
	fmt.Printf("[*] Source file is [%s], destnation file is [%s]\n", in, out)
	// 1.get original exe file
	fmt.Printf("[*] Starting get original exe file\n")
	data, err := internal.GetPEFileData(in)
	if err != nil {
		internal.CheckErr(err)
		return
	}
	fmt.Printf("[*] Get original exe file successful\n")
	// 2.encode file using xor
	fmt.Printf("[*] Starting encode original file using xor\n")
	encodeData, xorVal := internal.EncodeFile(data)
	fmt.Printf("[*] Encode original file successful\n")
	// 3.compress file using gzip
	fmt.Printf("[*] Starting compress file\n")
	fixData, err := internal.CompressFile(encodeData)
	if err != nil {
		internal.CheckErr(err)
		return
	}
	fmt.Printf("[*] Compress file successful\n")

	fixDataLen := len(fixData)
	fixDataLenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(fixDataLenBytes, uint32(fixDataLen))

	// 4.copy the shell file
	fmt.Printf("[*] Creating shell as a new file\n")
	err = internal.CopyFile(shellFile, out)
	if err != nil {
		internal.CheckErr(err)
		return
	}
	fmt.Printf("[*] Create new file successful\n")
	// 5.update resource to new file
	fmt.Printf("[*] Add source file as resource to new file\n")
	update, err := gowin32.NewResourceUpdate(out, false)
	if err != nil {
		internal.CheckErr(err)
		return
	}
	err = update.Update(gowin32.ResourceTypeRCData, gowin32.IntResourceId(uint(fileDataLenID)), gowin32.LanguageSystemDefault, fixDataLenBytes)
	if err != nil {
		internal.CheckErr(err)
		return
	}
	err = update.Update(gowin32.ResourceTypeRCData, gowin32.IntResourceId(uint(fileDataID)), gowin32.LanguageSystemDefault, fixData)
	if err != nil {
		internal.CheckErr(err)
		return
	}
	err = update.Update(gowin32.ResourceTypeRCData, gowin32.IntResourceId(uint(fileKeyID)), gowin32.LanguageSystemDefault, []byte{xorVal})
	if err != nil {
		internal.CheckErr(err)
		return
	}
	err = update.Save()
	if err != nil {
		internal.CheckErr(err)
		return
	}
	fmt.Printf("[*] Add source successful\n")
	fmt.Printf("[*] Packer is all done!\n")
}

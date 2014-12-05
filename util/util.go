package util

import (
	"os"
	"path"
	"io"
)

func checkError (err error) {
	if err != nil {
		panic(err)
	}
}
func MkdirAll(dstPath string, mode os.FileMode) {
	err := os.MkdirAll(dstPath, mode)
	checkError(err)
}

func CopyFile(srcPath, dstPath string, keepMode bool) {
	if !keepMode {
		//not keep mode
		MkdirAll(path.Dir(dstPath), (os.FileMode)(0777))
	}
	srcStat, err:= os.Lstat(srcPath)
	checkError(err)
	if srcStat.IsDir() {
		panic("srcPath is filepath")
	}
	srcFile, err := os.Open(srcPath)
	checkError(err)
	dstFile, err := os.Create(dstPath)
	checkError(err)
	io.Copy(dstFile, srcFile)
}

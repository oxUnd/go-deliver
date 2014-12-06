package util

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

func checkError(err error) {
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
	srcStat, err := os.Lstat(srcPath)
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

func hit(p string, include, exclude *regexp.Regexp) bool {
	_hit := true

	if include != nil {
		_hit = include.MatchString(p)
	}

	if _hit && exclude != nil {
		_hit = _hit && !exclude.MatchString(p)
	}

	return _hit
}

func Find(parameters ...interface{}) []string {
	if len(parameters) < 1 {
		panic("The first parameter must be a directory path.")
	}

	dir, err := filepath.Abs(parameters[0].(string))
	checkError(err)

	var include, exclude *regexp.Regexp
	if len(parameters) == 3 {
		include = parameters[1].(*regexp.Regexp)
		exclude = parameters[2].(*regexp.Regexp)
	} else if len(parameters) == 2 {
		include = parameters[3].(*regexp.Regexp)
	}
	files := []string{}
	infos, err := ioutil.ReadDir(dir)
	checkError(err)

	for _, info := range infos {
		path_ := path.Join(dir, info.Name())

		if !hit(path_, include, exclude) {
			continue
		}
		if info.IsDir() {
			files = append(Find(path_, include, exclude))
		} else {
			files = append(files, path_)
		}
	}
	return files
}

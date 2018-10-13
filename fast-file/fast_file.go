package fastFile

import (
	"os"
	"fmt"
	"golib/fast-common"
)

type Setting struct {
	RootDir string
	FileMode os.FileMode
	Flag int
}

//os.O_WRONLY|os.O_CREATE os.ModePerm
func (ctx *Setting) PushFileData(fileName, data string) error {
	var (
		file *os.File
		err error
		filePath string
	)

	filePath = fastCommon.JoinUrl(ctx.RootDir, fileName)

	if file, err = os.OpenFile(filePath, ctx.Flag, ctx.FileMode); err != nil{
		return fmt.Errorf("open file failure, path : %s", filePath)
	}

	defer file.Close()

	if _, err = file.Write([]byte(data)); err != nil{
		return fmt.Errorf("write data to file failure, path : %s", filePath)
	}

	return nil
}

func NewDir(rootDir string, flag int, fileMode os.FileMode) (*Setting, error) {
	if _, err := os.Stat(rootDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
				return nil, fmt.Errorf("create directory failure, path : %s", rootDir)
			}
			if err := os.Chmod(rootDir, os.ModePerm); err != nil{
				return nil, fmt.Errorf("chmod directory failure, path : %s", rootDir)
			}
		} else if os.IsPermission(err) {
			return nil, fmt.Errorf("permission denied, path : %s", rootDir)
		}else{
			return nil, err
		}
	}

	return &Setting{rootDir,fileMode,flag}, nil
}

package logging

import "os"

// fileExist ファイルが存在するかチェックします.
func fileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ファイルを作成します.
func fileCreate(path string) *os.File {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return file
}

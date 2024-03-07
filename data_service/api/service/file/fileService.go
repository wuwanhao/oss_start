package file

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// 文件服务
type FileService struct {
	StorageRoot  string
	StorageIndex string
}

func NewFileService(storageRoot, storageIndex string) *FileService {
	return &FileService{
		StorageRoot:  storageRoot,
		StorageIndex: storageIndex,
	}
}

// 获取文件流
func (fs *FileService) GetFile(filename string) (io.Reader, error) {
	filePath := fs.StorageRoot + fs.StorageIndex + filename
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer file.Close()
	return file, nil
}

// 上传文件
func put(w http.ResponseWriter, r *http.Request) {
	// 创建目录
	err := os.MkdirAll(os.Getenv("STORAGE_ROOT")+"/file/", 0755)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 创建文件
	file, err := os.Create(os.Getenv("STORAGE_ROOT") + "/file/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	io.Copy(file, r.Body)
}

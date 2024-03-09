package file

import (
	"api_service/common/config"
	"io"
	"log"
	"os"
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

// 存储文件流
func (fs *FileService) PutFile(fileStream io.ReadCloser, filename string) error {
	// 1.创建目录
	err := os.MkdirAll(config.Config.Oss.StorageRoot+config.Config.Oss.StorageIndex, 0755)
	log.Println(err)
	if err != nil {
		log.Println(err)
		return err
	}

	// 2.创建文件
	file, err := os.Create(config.Config.Oss.StorageRoot + config.Config.Oss.StorageIndex + filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	io.Copy(file, fileStream)
	return nil
}

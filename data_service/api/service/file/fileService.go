package file

import (
	"data_service/common/config"
	"data_service/utils"
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
	log.Println("[data_service] get file from path: " + filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("[data_service] get file failed from path: %s, err: %s", filePath, err)
		return nil, err
	}

	// 确保在函数返回前将文件指针重置到文件开头
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Printf("[data_service] seek file failed from path: %s, err: %s", filePath, err)
		return nil, err
	}
	return file, nil
}

// 存储文件流
func (fs *FileService) PutFile(fileStream io.ReadCloser, filename string) error {
	// 1. 检查并创建目录路径
	dirPath := config.Config.Oss.StorageRoot + config.Config.Oss.StorageIndex
	err := utils.CreateDir(dirPath)
	if err != nil {
		log.Println("[data_service] Error creating directory:", err)
		return err
	}

	// 2. 创建文件
	filePath := dirPath + filename
	log.Println("[data_service] File Path:", filePath)
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("[data_service] Error creating file:", err)
		return err
	}
	defer file.Close()

	// 3. 将文件流拷贝到新文件中
	_, err = io.Copy(file, fileStream)
	if err != nil {
		log.Println("[data_service] Error copying file:", err)
		return err
	}

	log.Println("[data_service] File successfully saved:", filePath)
	return nil
}

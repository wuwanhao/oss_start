package es

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func PutMetadata(name string, version int, size int64, hash string) error {
	//向ES中添加API: PUT /metadata/objects/<object_name>_<version>?op_type=create
	doc := fmt.Sprintf(`{"name":"%s", "version":"%d", "size":"%d", "hash":"%s"}`, name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d?op_type=create", os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	// 如果同时有多个客户端上传同一个元数据，结果会发生冲突，只有第一个文档被成功创建，即ES中已经有当前版本的了，升一个版本号，重新添加
	if response.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}

	if response.StatusCode != http.StatusCreated {
		result, _ := io.ReadAll(response.Body)
		return fmt.Errorf("failed to put metadata: %d %s", response.StatusCode, string(result))
	}

	return nil
}

// 获取对象元数据的最新版本，升版本后存入ES
func AddVersion(name, hash string, size int64) error {
	version, e := SearchLatestVersion(name)
	if e != nil {
		return e
	}
	return PutMetadata(name, version.Version+1, size, hash)
}

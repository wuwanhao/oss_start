package es

import (
	es "connector/es/entity"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func PutMetadata(name string, version int, size int64, hash string) error {
	// 向ES中查询
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

// 获取指定名称对象的所有版本，若不指定，则获取全部对象的全部版本  from, size 代表分页
func GetAllVersions(name string, from, size int) ([]es.Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?sort=name,version&from=%d&size=%d", os.Getenv("ES_SERVER"), from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	metadatas := make([]es.Metadata, 0)
	result, _ := io.ReadAll(response.Body)
	var sr es.SearchResult
	json.Unmarshal(result, &sr)
	for i := range sr.HIts.Hits {
		metadatas = append(metadatas, sr.HIts.Hits[i].Source)
	}
	return metadatas, nil
}

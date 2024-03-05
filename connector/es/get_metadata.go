package es

import (
	es "connector/es/entity"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// 根据对象的name和versionId从ES中获取对象的元数据
func GetMetadata(name string, version int) (es.Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

func getMetadata(name string, versionId int) (meta es.Metadata, e error) {
	// 精确定位，直接返回元数据的内容
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d/_source", os.Getenv("ES_SERVER"), name, versionId)
	r, e := http.Get(url)
	if e != nil {
		return
	}

	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("failed to get %s_%d: %d", name, versionId, r.StatusCode)
		return
	}

	result, _ := io.ReadAll(r.Body)
	json.Unmarshal(result, &meta)
	return
}

// 根据名称查找最新版本的元数据，版本号以降序排列只返回第一个结果
func SearchLatestVersion(name string) (meta es.Metadata, e error) {
	// 使用搜索，返回hit list
	url := fmt.Sprintf("http://%s/metadata/_search?q=name: %s&size=1&sort=version:desc", os.Getenv("ES_SERVER"), url.PathEscape(name))
	r, e := http.Get(url)
	if e != nil {
		return
	}

	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("failed to search latest metadata: %d", r.StatusCode)
		return
	}

	result, _ := io.ReadAll(r.Body)
	var sr es.SearchResult
	json.Unmarshal(result, &sr)
	if len(sr.HIts.Hits) != 0 {
		meta = sr.HIts.Hits[0].Source
	}
	return
}

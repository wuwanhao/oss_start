package es

import (
	"api_service/common/config"
	"common_service/logs"
	es "connector_service/es/entity"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// 根据对象的name和versionId从ES中获取对象的元数据
func GetMetadata(name string, version int) (es.Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

func getMetadata(name string, versionId int) (meta es.Metadata, e error) {
	// 精确定位，直接返回元数据的内容 API: GET /metadata/file/<object_name>_<version_id>/_source
	url := fmt.Sprintf("http://%s:%d/metadata/_doc/%s_%d/_source", config.Config.Es.Host, config.Config.Es.Port, name, versionId)
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

	// 使用搜索，返回hit list， ES搜索API: GET /metadata/_search?@=name: <object_name>&size=l&sort=version:desc
	esUrl := fmt.Sprintf("http://%s:%d/metadata/_search?q=name:%s&size=1&sort=version:desc", config.Config.Es.Host, config.Config.Es.Port, url.PathEscape(name))
	logs.Info("ES get metadata url: %v", esUrl)
	r, e := http.Get(esUrl)
	if e != nil {
		logs.Warn("ES get metadata error: %v", e)
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

// 获取指定名称对象的所有版本，若不指定，则获取全部对象的全部版本  from, size 代表分页
func GetAllVersions(name string, from, size int) ([]es.Metadata, error) {
	// GET全体对象版本列表API: GET /metadata/_search?sort=name, version&from=<from›&size=<size>
	// GET指定对象版本列表API: GET /metadata/_search?sort=name, version&from=<from>&size=<size>&q=name:‹object_name>
	es_url := fmt.Sprintf("http://%s:%d/metadata/_search?sort=name,version&from=%d&size=%d", config.Config.Es.Host, config.Config.Es.Port, from, size)
	if name != "" {
		es_url += "&q=name:" + name
	}
	response, err := http.Get(es_url)
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

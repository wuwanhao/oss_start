package es

// object的元数据
type Metadata struct {
	Name    string // 对象名
	Version int    // 对象版本
	Size    int64  // 对象大小
	Hash    string // 散列
}

// 与ES的搜索API返回的结构保持一致
/*
{
  "took":2,
  "timed_out":false,
  "_shards":{"total":5,"successful":5,"failed":0},
  "hits":{
    "total":2,
    "max_score":1.0,
    "hits":[
      {
        "_index":"accounts",
        "_type":"person",
        "_id":"AV3qGfrC6jMbsbXb6k1p",
        "_score":1.0,
        "_source": {
          "user": "李四",
          "title": "工程师",
          "desc": "系统管理"
        }
      },
      {
        "_index":"accounts",
        "_type":"person",
        "_id":"1",
        "_score":1.0,
        "_source": {
          "user" : "张三",
          "title" : "工程师",
          "desc" : "数据库管理，软件开发"
        }
      }
    ]
  }
}
*/
// ES搜索API返回的hits->hits->obj
type hit struct {
	Source Metadata `json:"_source"`
}

type SearchResult struct {
	HIts struct {
		total int
		Hits  []hit
	}
}

package versions

import (
	"connector_service/es"
)

func Handler(name string, from, size int) () {
	var sr es.SearchResult
	// 阻塞
	//for {
	//	metas, err := es.GetAllVersions(name, from, size)
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//
	//	for i := range metas {
	//		// 得到字节流
	//		b, _ := json.Marshal(metas[i])
	//		// 写入响应正文
	//		w.Write(b)
	//		w.Write([]byte("\n"))
	//	}
	//
	//	// 当前批是最后一批数据，metas的大小是小于size的，终止循环，返回最终的数据
	//	if len(metas) != size {
	//		return
	//	}
	//
	//	from += size
	//}

}

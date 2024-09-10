package versions

import (
	"connector_service/es"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	from := 0
	size := 1000
	name := strings.Split(r.URL.EscapedPath(), "/")[2] // 拿到key
	// 阻塞
	for {
		metas, err := es.GetAllVersions(name, from, size)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i := range metas {
			// 得到字节流
			b, _ := json.Marshal(metas[i])
			// 写入响应正文
			w.Write(b)
			w.Write([]byte("\n"))
		}

		// 当前批是最后一批数据，metas的大小是小于size的，终止循环，返回最终的数据
		if len(metas) != size {
			return
		}

		from += size
	}

}

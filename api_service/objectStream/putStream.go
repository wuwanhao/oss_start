package objectStream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct {
	writer *io.PipeWriter
	c      chan error
}

func NewPutStream(server, object string) *PutStream {
	// 创建一个无缓冲的pipe，写入writer的内容可以从reader中读出来
	reader, writer := io.Pipe()
	c := make(chan error)
	go func() {
		request, _ := http.NewRequest("PUT", "http://"+server+"/objects/"+object, reader)
		client := http.Client{}
		response, err := client.Do(request)
		// 重新构造错误
		if err == nil && response.StatusCode != http.StatusOK {
			err = fmt.Errorf("DataServer return http code %d", response.StatusCode)
		}
		c <- err
	}()

	return &PutStream{
		writer: writer,
		c:      c,
	}
}

func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *PutStream) Close() error {
	w.writer.Close()
	return <-w.c
}

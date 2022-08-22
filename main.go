package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// 1. Добавить в пример с файловым сервером возможность получить список всех файлов
// на сервере (имя, расширение, размер в байтах)

type sizeInserter struct {
	basePath string
	buf      *bytes.Buffer
}

func (si *sizeInserter) Process() {

	node, err := html.Parse(si.buf)
	if err != nil {
		fmt.Println("failed ot parse html")
	}

	var preEl *html.Node
	var bodyEl *html.Node
	var getPreFunc func(*html.Node)
	var getBodyFunc func(*html.Node)
	var addSize func(*html.Node)

	getPreFunc = func(node *html.Node) {
		if node == nil {
			return
		}
		getBodyFunc(node)
		if bodyEl != nil {
			preEl = bodyEl.FirstChild
			return
		}
	}
	getBodyFunc = func(node *html.Node) {
		if node == nil {
			return
		}
		if node.Type == html.DocumentNode || node.Type == html.ElementNode && node.Data == "html" {
			getBodyFunc(node.FirstChild)
		}
		if node.Data == "head" {
			getBodyFunc(node.NextSibling)
		}
		if node.Data == "body" {
			bodyEl = node
			return
		}
	}

	getPreFunc(node)

	addSize = func(node *html.Node) {
		if node.FirstChild != nil { //text
			if strings.ContainsRune(node.FirstChild.Data, os.PathSeparator) {
				return
			}
			curPath, _ := os.Getwd()
			filePath := filepath.Join(curPath, si.basePath, node.FirstChild.Data)
			fileInfo, err := os.Stat(filePath)
			var fileSizeStr string
			if err != nil {
				fmt.Println("fail to get fileInfo: ", err)
			} else {
				fileSizeStr = fmt.Sprintf(" \t %v bytes", fileInfo.Size())
			}

			node.FirstChild.Data += fileSizeStr
		}
	}

	if preEl != nil {
		for chNode := preEl.FirstChild; chNode != nil; chNode = chNode.NextSibling {
			if chNode.FirstChild != nil { //text
				addSize(chNode)
			}
		}
		html.Render(si.buf, preEl)
		return
	}

	tempBuf := &bytes.Buffer{}
	html.Render(tempBuf, node)

	if bodyEl.FirstChild != nil {
		si.buf.WriteString(bodyEl.FirstChild.Data)
	}
}

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (mrw *MyResponseWriter) Write(p []byte) (int, error) {
	return mrw.buf.Write(p)
}

type fsWrapper struct {
	dirToServe http.Dir
	handler    http.Handler
}

func createFSWrapper(dir2Server string) fsWrapper {
	wr := fsWrapper{dirToServe: http.Dir(dir2Server)}
	wr.handler = http.FileServer(wr.dirToServe)
	return wr
}

func (h *fsWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	// Create a response wrapper to fill buf with the responce:
	mrw := &MyResponseWriter{
		ResponseWriter: w,
		buf:            &bytes.Buffer{},
	}

	h.handler.ServeHTTP(mrw, r)

	mSizeInserter := sizeInserter{
		basePath: filepath.Join(string(h.dirToServe), r.RequestURI),
		buf:      mrw.buf,
	}

	mSizeInserter.Process()

	if _, err := io.Copy(w, mrw.buf); err != nil {
		fmt.Printf("Failed to copy response: %v", err)
	}
}

func main() {
	myHandler := createFSWrapper("files")

	fs := &http.Server{
		Addr:         ":8080",
		Handler:      &myHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fs.ListenAndServe()
}

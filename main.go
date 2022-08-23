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

// С помощью query-параметра, реализовать фильтрацию выводимого списка по
// расширению (то есть, выводить только .png файлы, или только .jpeg)

type intermProcessor struct {
	basePath string
	buf      *bytes.Buffer
	preEl    *html.Node
	bodyEl   *html.Node
}

func (si *intermProcessor) getBody(node *html.Node) {
	if node == nil {
		return
	}
	if node.Type == html.DocumentNode || node.Type == html.ElementNode && node.Data == "html" {
		si.getBody(node.FirstChild)
	}
	if node.Data == "head" {
		si.getBody(node.NextSibling)
	}
	if node.Data == "body" {
		si.bodyEl = node
		return
	}
}

func (si *intermProcessor) getPre(node *html.Node) {
	if node == nil {
		return
	}
	si.getBody(node)
	if si.bodyEl != nil {
		si.preEl = si.bodyEl.FirstChild
		return
	}
}

func (si *intermProcessor) Process(ext string) {
	si.FilterByExtension(ext)
	si.insertFileSize()
}

func (si *intermProcessor) insertFileSize() {

	node, err := html.Parse(si.buf)
	if err != nil {
		fmt.Println("failed ot parse html")
	}

	var addSize func(*html.Node)

	si.getPre(node)

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

	if si.preEl != nil {
		for chNode := si.preEl.FirstChild; chNode != nil; chNode = chNode.NextSibling {
			if chNode.FirstChild != nil { //text
				addSize(chNode)
			}
		}
		html.Render(si.buf, si.preEl)
		return
	}

	tempBuf := &bytes.Buffer{}
	html.Render(tempBuf, node)

	if si.bodyEl.FirstChild != nil {
		si.buf.WriteString(si.bodyEl.FirstChild.Data)
	}
}

func (si *intermProcessor) FilterByExtension(ext string) {

	if ext == "" {
		return
	}

	node, err := html.Parse(si.buf)
	if err != nil {
		fmt.Println("failed ot parse html")
	}

	var filterANode func(*html.Node)

	si.getPre(node)

	filterANode = func(node *html.Node) {
		if node.FirstChild != nil { //text node
			if strings.ContainsRune(node.FirstChild.Data, os.PathSeparator) {
				return
			}
			if _, afterStr, res := strings.Cut(filepath.Ext(node.FirstChild.Data), "."); res {
				if afterStr != "" && afterStr != ext {
					attr := html.Attribute{Key: "style", Val: "display:none;"}
					node.Attr = append(node.Attr, attr)
					node.Parent.RemoveChild(node.NextSibling)
				}
			}
		}
	}

	if si.preEl != nil {
		for chNode := si.preEl.FirstChild; chNode != nil; chNode = chNode.NextSibling {
			if chNode.FirstChild != nil { //text
				filterANode(chNode)
			}
		}
		html.Render(si.buf, si.preEl)
		return
	}

	if si.bodyEl.FirstChild != nil {
		si.buf.WriteString(si.bodyEl.FirstChild.Data)
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

	defer func() {
		if r.Body != nil {
			r.Body.Close()
		}
	}()

	var ext string
	if r.Method == http.MethodGet {
		ext = r.FormValue("ext")
	}

	// Create a response wrapper to fill buf with the responce:
	mrw := &MyResponseWriter{
		ResponseWriter: w,
		buf:            &bytes.Buffer{},
	}

	h.handler.ServeHTTP(mrw, r)

	mIntermProcessor := intermProcessor{
		basePath: filepath.Join(string(h.dirToServe), r.URL.Path),
		buf:      mrw.buf,
	}

	mIntermProcessor.Process(ext)

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

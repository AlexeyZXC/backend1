package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTxtExtension(t *testing.T) {

	req, err := http.NewRequest("GET", "/?ext=txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	myHandler := createFSWrapper("files")

	myHandler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("Returned status: %v", status)
	}

	expected := "<pre><a href=\"file1.txt\">file1.txt \t 8 bytes</a>\n<a href=\"file2.avi\" style=\"display:none;\">file2.avi \t 0 bytes</a><a href=\"sub2/\">sub2/</a>\n<a href=\"subfolder1/\">subfolder1/</a>\n</pre>"

	if resp.Body.String() != expected {
		t.Errorf("Body is wrong, expected:\n%v\nActual:\n%v", expected, resp.Body.String())
	}
}

func TestAviExtension(t *testing.T) {

	req, err := http.NewRequest("GET", "/?ext=avi", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()

	myHandler := createFSWrapper("files")

	myHandler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("Returned status: %v", status)
	}

	expected := "<pre><a href=\"file1.txt\" style=\"display:none;\">file1.txt 	 8 bytes</a><a href=\"file2.avi\">file2.avi 	 0 bytes</a>\n<a href=\"sub2/\">sub2/</a>\n<a href=\"subfolder1/\">subfolder1/</a>\n</pre>"

	if resp.Body.String() != expected {
		t.Errorf("Body is wrong, expected:\n%v\nActual:\n%v", expected, resp.Body.String())
	}
}

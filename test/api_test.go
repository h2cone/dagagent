package test

import (
	"bytes"
	"encoding/json"
	"github.com/h2cone/dagagent/api/server"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func before(wg *sync.WaitGroup) {
	defer wg.Done()
	server.Start()
}

func TestUpload(t *testing.T) {
	var wg sync.WaitGroup
	go before(&wg)
	wg.Wait()

	url := "http://" + server.Address + "/upload"
	method := "POST"
	src := filepath.Join("dags", "example", "example_complex.py")

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.WriteField("subDir", "example")
	if err != nil {
		t.Error(err)
	}
	srcFile, err := os.Open(src)
	if err != nil {
		t.Fatal(err)
	}
	defer srcFile.Close()
	srcStat, err := os.Stat(src)
	if err != nil {
		t.Error(err)
	}

	part, err := writer.CreateFormFile("file", filepath.Base(src))
	if err != nil {
		t.Error(err)
	}
	_, err = io.Copy(part, srcFile)
	if err != nil {
		t.Error(err)
	}
	err = writer.Close()
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Basic YWlyZmxvdzphaXJmbG93")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		t.Fatalf("StatusCode: %d, Status: %s", res.StatusCode, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	dst := string(body)
	t.Logf("body: %s", dst)

	var result map[string]string
	json.Unmarshal(body, &result)
	dstStat, err := os.Stat(result["fileloc"])
	if err != nil {
		t.Error(err)
	}
	srcModtime := srcStat.ModTime()
	dstModTime := dstStat.ModTime()
	t.Logf("srcModtime: %v, dstModTime: %v\n", srcModtime, dstModTime)
	if srcStat.ModTime().Equal(dstStat.ModTime()) {
		t.Error("modTime unchanged")
	}
}

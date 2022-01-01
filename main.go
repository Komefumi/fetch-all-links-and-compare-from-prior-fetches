package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

const fetchContentDirName string = "fetches"
const timeFormatString string = "2006-01-02---15_04_05"

var urlExtractRegex *regexp.Regexp = regexp.MustCompile("^https?://(.*)/?$")

func main() {
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	for _, url := range os.Args[1:] {
		urlFetchDir, err := getUrlFetchDir(url)
		if err != nil {
			fmt.Sprintf("While examining for %s: %v", url, err)
			continue
		}

		outputDirRead, err := os.Open(urlFetchDir)
		if err != nil {
			fmt.Sprintf("While examining for %s: %v", url, err)
			continue
		}

		outputDirFiles, err := outputDirRead.Readdir(0)
		for outputIndex := range outputDirFiles {
			outputFileHere := outputDirFiles[outputIndex]
			outputNameHere := outputFileHere.Name()
			fileSize := outputFileHere.Size()
			fmt.Printf("%s: %7d bytes\n", outputNameHere, fileSize)
		}
	}
}

func getUrlFetchDir(url string) (string, error) {
	extractedFromUrl := urlExtractRegex.FindStringSubmatch(url)
	urlDirName := fmt.Sprintf("%s/%s", fetchContentDirName, extractedFromUrl[1])
	fmt.Printf("urlDirName: %v\n", urlDirName)
	_, err := os.Stat(fetchContentDirName)
	if !os.IsNotExist(err) && err != nil {
		return "", err
	}

	err = os.Mkdir(fetchContentDirName, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return "", err
	}
	err = os.Mkdir(urlDirName, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return "", err
	}

	return urlDirName, nil
}

func fetch(url string, ch chan string) {
	now := time.Now()
	fileName := fmt.Sprintf("%s.html", now.Format(timeFormatString))
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	urlDir, err := getUrlFetchDir(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	filePath := fmt.Sprintf("%s/%s", urlDir, fileName)
	fmt.Println(urlDir)
	fmt.Println(filePath)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	fmt.Println(f)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	_, err = io.Copy(f, resp.Body)
	resp.Body.Close()
	ch <- fmt.Sprintf("%s fetched to %s", url, filePath)
	return
}

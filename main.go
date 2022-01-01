package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const fetchContentDirName string = "fetches"
const timeFormatString string = "2006-01-02---15_04_05"

var urlExtractRegex *regexp.Regexp = regexp.MustCompile("^https?://(.*)/?$")
var htmlExtensionRegex *regexp.Regexp = regexp.MustCompile("\\.html")
var isOkayRegex *regexp.Regexp = regexp.MustCompile("fetched")

func main() {
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}
	for range os.Args[1:] {
		returnedString := <-ch
		if !isOkayRegex.MatchString(returnedString) {
			fmt.Println(returnedString)
		}
	}
	urlCount := len(os.Args[1:])
	for urlIndex, url := range os.Args[1:] {
		urlFetchDir, err := getUrlFetchDir(url)
		if err != nil {
			fmt.Sprintf("While examining for %s: %v", url, err)
			continue
		}

		fmt.Println("-------")
		fmt.Printf("For %s:\n", urlFetchDir)
		filepath.Walk(urlFetchDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatalf(err.Error())
			}
			if info.IsDir() {
				return nil
			}
			outputNameHere := info.Name()
			justFileName := htmlExtensionRegex.ReplaceAllString(outputNameHere, "")
			dateExtracted, err := time.Parse(timeFormatString, justFileName)
			if err != nil {
				panic(err)
			}
			fileSize := info.Size()
			fmt.Printf("%s: %7d bytes\n", dateExtracted.Format(time.RFC822), fileSize)
			return nil
		})
		fmt.Println("-------\n")
		if urlIndex < urlCount-1 {
			fmt.Println()
		}
	}
}

func getUrlFetchDir(url string) (string, error) {
	extractedFromUrl := urlExtractRegex.FindStringSubmatch(url)
	urlDirName := fmt.Sprintf("%s/%s", fetchContentDirName, extractedFromUrl[1])
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
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	_, err = io.Copy(f, resp.Body)
	resp.Body.Close()
	ch <- fmt.Sprintf("%s fetched to %s", url, filePath)
	return
}

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func StoreImages(email string) (err error) {
	// INIT - vars
	var dirName string

	// SPLIT - parts
	m := strings.Split(email, "\n------=_")

	for i, v := range m {
		// FILTER OUT - non image parts
		// CASE 1 - non-image parts
		if !strings.Contains(v, "image/jpeg") {
			continue
		}

		// CASE 2 - image part
		fmt.Println(i, len(v))

		// SPLIT - header, body
		mm := strings.Split(v, "\r\n\r\n")
		// CASE 2-1 - wrong image part
		if len(mm) < 2 {
			fmt.Println("ERR0: wrong image part")
			continue
		}

		// CASE 2-2 - valid image part
		var (
			header   = mm[0]
			filename = ""
			body     = strings.TrimSpace(strings.ReplaceAll(mm[1], "\n", ""))
		)

		// EXTRACT & SET - filename
		mmm := strings.Split(header, "filename=\"")
		// CASE 3-1 - wrong image header
		if len(mmm) < 2 {
			fmt.Println("ERR0: wrong image header")
			continue
		}

		// CASE 3-2 - valid image header
		// SET - filename
		filename = strings.ReplaceAll(mmm[1], "\"", "")
		if dirName == "" {
			timestampStr := strings.ReplaceAll(filename, ".jpg", "")
			timestamp, err := strconv.ParseInt(strings.Split(timestampStr, "_")[0], 10, 64)
			if err != nil {
				fmt.Println("ERR0: ", err)
				continue
			}
			eventTime := time.Unix(timestamp, 0)
			dirName = eventTime.String()
			os.MkdirAll(dirName, os.ModePerm)
		}

		// DECODE - base64 image
		imageBytes, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			fmt.Println("ERR1: ", err)
			continue
		}

		// WRITE - file
		filepath := fmt.Sprintf("%s/%s", dirName, filename)
		err = ioutil.WriteFile(filepath, imageBytes, 0644)
		if err != nil {
			fmt.Println("ERR2: ", err)
			continue
		}
	}
	return
}

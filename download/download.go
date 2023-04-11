package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"
)

var cli client.Client

func main() {
	if len(os.Args) < 3 {
		exit("it needs 2 params, rpm url and token")
	}

	rpmUrl := os.Args[1]
	token := os.Args[2]
	u, err := url.Parse(rpmUrl)
	checkErr(err)

	v := strings.Split(u.Path, "/")
	org := v[1]
	repo := v[2]
	fileName := v[len(v)-1]

	if v[3] != "raw" {
		exit("source file must be raw format")
	}

	branches := getBranches(token, org, repo)

	for _, branch := range branches {
		if strings.Contains(rpmUrl, branch.Name) {
			filePath := strings.Split(u.Path, branch.Name)[1]

			content, err := cli.GetPathContent(org, repo, filePath, branch.Name)
			if err != nil {
				continue
			}

			decodeContent, err := base64.StdEncoding.DecodeString(content.Content)
			checkErr(err)

			err = os.WriteFile(fileName, decodeContent, 0644)
			checkErr(err)

			break
		}
	}

	_, err = os.Stat(fileName)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("download %s failed", fileName)
		os.Exit(1)
	}
}

func getBranches(token, org, repo string) []gitee.Branch {
	cli = client.NewClient(func() []byte {
		return []byte(token)
	})

	branches, err := cli.GetRepoAllBranch(org, repo)
	checkErr(err)

	return branches
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/lspaccatrosi16/go-cli-tools/input"
)

type ScaffoldData struct {
	TargetPath string
	Template   string
}

func main() {
	fmt.Println("Hello world")
}

//scaffold does this:
// get template from gh
// copy it to target dir
// runs and then deletes postinstall.sh in target

func getTemplates() ([]string, []*[]byte) {
	resp, err := http.Get("https://github.com/lspaccatrosi16/scaffold/archive/refs/heads/master.zip")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		panic(err)
	}

	for _, zipFile := range zipReader.File {
		unzipped, err := readZippedFile(zipFile)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Reading %s\n", zipFile.Name)

		_ = unzipped
	}
	return []string{}, []*[]byte{}
}

func readZippedFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return []byte{}, err
	}

	defer f.Close()
	return ioutil.ReadAll(f)
}

func getData(availTemplates []string) ScaffoldData {
folder_input:
	targetPath := input.GetInput("Target path")
	stats, err := os.Stat(targetPath)

	if err != nil {
		if os.IsNotExist(err) {
			createErr := os.Mkdir(targetPath, 0o755)
			if createErr != nil {
				panic(fmt.Sprintf("Error creating target path: %s", createErr.Error()))
			}
		} else {
			panic(err)
		}
	}

	if !stats.IsDir() {
		fmt.Printf("Path %s is not a directory", targetPath)
		goto folder_input
	}
	sort.Strings(availTemplates)
	options := []input.SelectOption{}
	for _, name := range availTemplates {
		options = append(options, input.SelectOption{
			Name:  name,
			Value: name,
		})
	}

	selectedTemplate, err := input.GetSelection("Pick a template", options)
	if err != nil {
		panic(err)
	}

	return ScaffoldData{
		TargetPath: targetPath,
		Template:   selectedTemplate,
	}
}

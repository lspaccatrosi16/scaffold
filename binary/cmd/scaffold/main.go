package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/lspaccatrosi16/go-cli-tools/input"
	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/mandelsoft/vfs/pkg/memoryfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
)

type ScaffoldData struct {
	TargetPath string
	Template   string
}

func main() {
	logger := logging.GetLogger()
	logger.LogDivider()
	logger.Log("Getting Template List")
	templateNames, templates := getTemplates()
	logger.Log("Got Template List")
	logger.LogDivider()
	selected := getData(templateNames)
	chosenTemplate, ok := templates[selected.Template]
	if !ok {
		panic("selected template was selected but does not exist in map")
	}
	createScaffold(selected.TargetPath, chosenTemplate)
	executePostInstall(selected.TargetPath)
	logger.LogDivider()
	logger.Log("Done")
}

//scaffold does this:
// get template from gh
// copy it to target dir
// runs and then deletes postinstall.sh in target

func executePostInstall(path string) {
	expectedPath := filepath.Join(path, "postinstall.sh")
	_, err := os.Stat(expectedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			panic(err)
		}
	}
	cmd := exec.Command("/bin/sh", expectedPath)
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
	err = os.Remove(expectedPath)
	if err != nil {
		panic(err)
	}
}

func createScaffold(path string, templateZip *[]byte) {
	mfs := memoryfs.New()
	_, err := unzipFolder(templateZip, "", mfs)
	if err != nil {
		panic(err)
	}
	err = vfsToDisk(path, "", mfs)
	if err != nil {
		panic(err)
	}
}

func getTemplates() ([]string, map[string]*[]byte) {
	mfs := memoryfs.New()
	resp, err := http.Get("https://github.com/lspaccatrosi16/scaffold/archive/refs/heads/master.zip")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	baseFolderName, err := unzipFolder(&data, "", mfs)
	if err != nil {
		panic(err)
	}
	templates := map[string]*[]byte{}
	templateNames := []string{}
	templatesPath := filepath.Join(baseFolderName, "templates")
	templateFolders, err := vfs.ReadDir(mfs, templatesPath)
	if err != nil {
		panic(err)
	}
	for _, folder := range templateFolders {
		if !folder.IsDir() {
			continue
		}
		templateNames = append(templateNames, folder.Name())
		startPath := filepath.Join(templatesPath, folder.Name())
		zipped, err := zipFolder(startPath, mfs)
		if err != nil {
			panic(err)
		}
		templates[folder.Name()] = zipped
	}
	return templateNames, templates
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
	} else {
		if !stats.IsDir() {
			fmt.Printf("Path %s is not a directory", targetPath)
			goto folder_input
		}
	}
	sort.Strings(availTemplates)
	options := []input.SelectOption{}
	for _, name := range availTemplates {
		options = append(options, input.SelectOption{
			Name:  name,
			Value: name,
		})
	}

	selectedTemplate, err := input.GetSearchableSelection("Pick a template", options)
	if err != nil {
		panic(err)
	}

	return ScaffoldData{
		TargetPath: targetPath,
		Template:   selectedTemplate,
	}
}

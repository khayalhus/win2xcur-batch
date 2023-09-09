package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

// Mappings are the top-level object holding
type Mappings struct {
	Pairs []Pair `json:"mappings"`
}

// Pair holds the corresponding Linux cursors for the Windows cursor
type Pair struct {
	Windows string   `json:"windows"`
	Linux   []string `json:"linux"`
}

func main() {
	pwd, _ := os.Getwd()
	windowsDir := "Unzipped" // Replace with the path to your directory
	fullPathToWindowsDir := pwd + "/" + windowsDir

	// Read windows to linux mappings from file
	file, _ := ioutil.ReadFile("map.json")
	mapped := Mappings{}
	_ = json.Unmarshal([]byte(file), &mapped)

	// Read directory containing original files
	directories, lsErr := os.ReadDir(fullPathToWindowsDir)
	if lsErr != nil {
		fmt.Println(lsErr)
		return
	}

	// Convert from Windows to Linux
	convertDir := "Converted"
	fullPathToConvertedDir := pwd + "/" + convertDir
	convertMkdirErr := os.MkdirAll(fullPathToConvertedDir, 0700)
	if convertMkdirErr != nil {
		fmt.Println(convertMkdirErr)
		return
	}
	converterFolderRequirement := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	for i := range directories {
		if !directories[i].IsDir() {
			continue
		}
		dirName := directories[i].Name()
		fullPathToDir := fullPathToWindowsDir + "/" + dirName
		relativePathToDir := windowsDir + "/" + dirName
		matched := converterFolderRequirement.Match([]byte(dirName))
		if matched {
			dirName = converterFolderRequirement.ReplaceAllString(dirName, "")
			relativePathToDir = windowsDir + "/" + dirName
			renameErr := os.Rename(fullPathToDir, pwd+"/"+relativePathToDir)
			if renameErr != nil {
				fmt.Println(renameErr)
				return
			}
		}

		convertMkdirErr := os.MkdirAll(fullPathToConvertedDir+"/"+dirName, 0700)
		if convertMkdirErr != nil {
			fmt.Println(convertMkdirErr)
			return
		}
		convertCmd := fmt.Sprintf("win2xcur %s/*.ani -o %s", relativePathToDir, convertDir+"/"+dirName)
		fmt.Println(convertCmd)
		cmd := exec.Command("bash", "-c", convertCmd)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Sort converted files
	sortDir := pwd + "/" + "Sorted"
	sortedMkdirErr := os.MkdirAll(sortDir, 0700)
	if sortedMkdirErr != nil {
		fmt.Println(sortedMkdirErr)
		return
	}
	directories, lsErr = os.ReadDir(fullPathToConvertedDir)
	if lsErr != nil {
		fmt.Println(lsErr)
		return
	}
	for i := range directories {
		fmt.Println("Sorting " + directories[i].Name())
		linuxDir := sortDir + "/" + directories[i].Name()
		orgDir := fullPathToConvertedDir + "/" + directories[i].Name()
		linuxCursorDir := linuxDir + "/" + "cursors"
		os.Mkdir(linuxDir, 0744)
		os.Mkdir(linuxCursorDir, 0744)
		winFiles, lsErr := os.ReadDir(orgDir)
		if lsErr != nil {
			fmt.Println(lsErr)
			return
		}

		for _, winFile := range winFiles {
			for _, pair := range mapped.Pairs {
				if pair.Windows == winFile.Name() {
					for _, linuxName := range pair.Linux {
						fmt.Println("  " + winFile.Name() + " -> " + linuxName)
						input, err := ioutil.ReadFile(orgDir + "/" + winFile.Name())
						if err != nil {
							fmt.Println(err)
							return
						}

						err = ioutil.WriteFile(linuxCursorDir+"/"+linuxName, input, 0644)
						if err != nil {
							fmt.Println("Error creating", linuxCursorDir+"/"+linuxName)
							fmt.Println(err)
							return
						}
					}
				}
			}
		}

		// Add cursor.theme
		cursorConfig := fmt.Sprintf("[Cursor Theme]\nName=%s\nInherits=Bibata-Modern-Classic\n", directories[i].Name())
		byteCursorConfig := []byte(cursorConfig)
		err := ioutil.WriteFile(linuxDir+"/"+"cursor.theme", byteCursorConfig, 0644)
		if err != nil {
			fmt.Println("Error creating", linuxDir+"/"+"cusror.theme")
			fmt.Println(err)
			return
		}
	}
}

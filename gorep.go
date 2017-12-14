package main

import (
	"path/filepath"
	"flag"
	"fmt"

	"strings"
	"io/ioutil"
	"os"
	"github.com/Pallinder/go-randomdata"
)

func main() {
	var flagPath, flagBanned string
	flag.StringVar(&flagPath, "path", "", "path files to replace")
	flag.StringVar(&flagBanned, "banned", "", "comma separated words that are banned")
	flag.Parse()
	banned := strings.Split(flagBanned, ",")
	if flagBanned == "" || len(banned) == 0 {
		fmt.Println("ERROR", `argument -banned="" is required. Remember comma separated`)
		return
	}

	bans := map[string]string{}
	unique := map[string]bool{}
	for _, ban := range banned {
		//Keep trying new replacements until we find a unique one again.
		var replace string
		for replace == "" || unique[replace] {
			replace = randomdata.SillyName()
		}
		bans[ban] = replace
	}

	if flagPath == "." || flagPath == "" {
		flagPath, _ = os.Getwd()
	}

	pathsFound := []string{}

	err := filepath.Walk(flagPath, func(path string, info os.FileInfo, err error) error {
		for ban, _ := range bans {
			if strings.Contains(path, ban) {
				if info.IsDir() {
					pathsFound = append(pathsFound, path)
				} else {
					pathsFound = append([]string{path}, pathsFound...)
				}
			}
		}
		if info.IsDir() {
			return nil
		}

		bts, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(bts)

		for ban, replace := range bans {
			content = strings.Replace(content, ban, replace, -1)
			content = strings.Replace(content, ban, replace, -1)
		}

		err = ioutil.WriteFile(path, []byte(content), info.Mode())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println("ERROR", err.Error())
	}

	for _, path := range pathsFound {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		newPath := path
		if !info.IsDir() {
			newPath = info.Name()
		}
		for ban, replace := range bans {
			newPath = strings.Replace(newPath, ban, replace, -1)
		}
		if !info.IsDir() {
			newPath = strings.Replace(path, info.Name(), newPath, -1)
		}
		os.Rename(path, newPath)
		println(fmt.Sprintf("Replacing %s with %s", path, newPath))
	}

}

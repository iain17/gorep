package main

import (
	"path/filepath"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
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

	basedir := "/tmp/"+time.Now().String()+"/"
	err := filepath.Walk(flagPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		newPath := basedir+path
		for ban, replace := range bans {
			newPath = strings.Replace(newPath, ban, replace, -1)
		}
		
		bts, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		content := string(bts)
		for ban, replace := range bans {
			content = strings.Replace(content, ban, replace, -1)
		}
		os.MkdirAll(filepath.Dir(newPath), 0777)

		err = ioutil.WriteFile(newPath, []byte(content), info.Mode())
		if err != nil {
			return err
		}

		return nil
	})

	os.RemoveAll(flagPath)
	os.Rename(basedir+flagPath, flagPath+"/")

	if err != nil {
		panic(err)
	}

}

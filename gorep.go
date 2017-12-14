package main

import "path/filepath"
import "flag"
import "fmt"

import "strings"
import "io/ioutil"
import "os"
import "github.com/Pallinder/go-randomdata"

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

	found := []string{}

	err := filepath.Walk(flagPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			for ban, replace := range bans {
				if strings.Contains(info.Name(), ban) {
					os.Rename(path, strings.Replace(path, ban, replace, -1))
				}
			}
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
		found = append(found, path)

		err = ioutil.WriteFile(path, []byte(content), info.Mode())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println("ERROR", err.Error())
	}

	for _, path := range found {
		fmt.Printf("found in %s\n", path)
	}

	if len(found) == 0 {
		fmt.Println("Nothing replaced")
	} else {
		fmt.Printf("Total %d file replaced\n", len(found))
	}

}

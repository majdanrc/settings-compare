package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
)

type SettingsFile struct {
	File        string
	Dict        Dict
	RawSettings []Setting `xml:"variable,omitempty"`
}

type Setting struct {
	Key   string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Dict map[string]string

func main() {
	files := os.Args[1:]

	var dicts []SettingsFile

	for _, f := range files {
		dicts = append(dicts, readFile(f))
	}

	keys := make(map[string]bool)

	for _, d := range dicts {
		for k := range d.Dict {
			keys[k] = true
		}
	}

	for k := range keys {
		var values []string

		for _, d := range dicts {
			values = append(values, d.Dict[k])
		}

		if !allSameStrings(values) {
			fmt.Printf("%s -> ", k)
			for c, v := range values {
				fmt.Printf("[")
				var cn color.Attribute
				cn = 91 + color.Attribute(c)

				color.Set(cn)
				fmt.Printf("%s", v)
				color.Unset()
				fmt.Printf("]")
			}
			fmt.Println()
		}
	}

	fmt.Println("------------------")
	fmt.Println(dicts[0])
}

func readFile(filePath string) SettingsFile {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var settings SettingsFile
	settings.File = filePath

	xml.Unmarshal(byteValue, &settings)

	settings.Dict = make(Dict)

	for _, r := range settings.RawSettings {
		settings.Dict[r.Key] = r.Value
	}

	return settings
}

func allSameStrings(a []string) bool {
	for i := 1; i < len(a); i++ {
		if a[i] != a[0] {
			return false
		}
	}
	return true
}

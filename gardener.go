package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"github.com/gomarkdown/markdown"
	"path"
	"strings"
	"path/filepath"
)

func main() {
	files, err := ioutil.ReadDir("./mdfiles")
	if (err != nil) {
		log.Fatal(err)
	}

	for _, f := range files {
		fileName := f.Name()
		input_filepath := path.Join("./mdfiles", fileName)
		file, err := os.Open(input_filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		b, err := ioutil.ReadAll(file)
		output := markdown.ToHTML(b, nil, nil)

		output_filename := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		ouput_filepath := path.Join("./ouput", output_filename + ".html")
		
		err = ioutil.WriteFile(ouput_filepath, output, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Completed!")
}

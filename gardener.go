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
	"flag"
)

func main() {

	var inputDir string;
	flag.StringVar(&inputDir, "input", "", "Directory of markdown files")

	var outputDir string;
	flag.StringVar(&outputDir, "output", "", "Directory to generate html files")

	var checkString string;
	flag.StringVar(&checkString, "check-string", "", "Include only files containing check string")

	flag.Parse()

	files, err := ioutil.ReadDir(inputDir)
	if (err != nil) {
		log.Fatal(err)
	}

	processedFileCount := 0

	for _, f := range files {
		fileName := f.Name()
		input_file_extension := filepath.Ext(fileName)
		if (input_file_extension != ".md") {
			fmt.Println("Skipping", fileName)
			continue
		}

		input_filepath := path.Join(inputDir, fileName)

		file, err := os.Open(input_filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		b, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		if (len(b) == 0) {
			continue
		}

		content := string(b)
		if (!strings.Contains(content, checkString)) {
			continue
		}

		output := markdown.ToHTML(b, nil, nil)

		output_filename := strings.TrimSuffix(fileName, input_file_extension)
		output_filename = strings.ReplaceAll(output_filename, " ", "-")
		ouput_filepath := path.Join(outputDir, output_filename + ".html")

		err = ioutil.WriteFile(ouput_filepath, output, 0644)
		if err != nil {
			log.Fatal(err)
		}
		processedFileCount += 1
	}
	fmt.Println(processedFileCount,"/", len(files),  "files processed")
}

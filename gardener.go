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
	"regexp"
)

func getMDFileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, ".md")
}

func getHTMLFilePathFromFileName(fileName string) string {
	var mdFileNameWithoutExtension = getMDFileNameWithoutExtension(fileName)
	return strings.ToLower(strings.ReplaceAll(mdFileNameWithoutExtension, " ", "-")) + ".html"
}

func getLinksFromContent(content string) []string {
	re := regexp.MustCompile(`\[([^\[\]]*)\]`)
	matches := re.FindAllStringSubmatch(content, -1)
	var links []string
	for _, match := range matches {
		links = append(links, match[1])
	}
	return links
}

func main() {

	var inputDir string;
	flag.StringVar(&inputDir, "input", "", "Directory of markdown files")

	var outputDir string;
	flag.StringVar(&outputDir, "output", "", "Directory to generate html files")

	var checkString string;
	flag.StringVar(&checkString, "check-string", "", "Include only files containing check string")

	var baseTemplatePath string;
	flag.StringVar(&baseTemplatePath, "base-template", "", "Base template file")

	flag.Parse()

	if (len(outputDir) == 0) {
		fmt.Println("You need to speceify an output directory!")
		os.Exit(1)
	}

	os.RemoveAll(outputDir)
	os.MkdirAll(outputDir, 0755)

	files, err := ioutil.ReadDir(inputDir)
	if (err != nil) {
		log.Fatal(err)
	}

	processedFileCount := 0
	fileNameToContentMap := make(map[string][]byte)
	backLinksMap := make(map[string][]string)

	for _, f := range files {
		fileName := f.Name()
		input_file_extension := filepath.Ext(fileName)
		if (input_file_extension != ".md") {
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
		if len(checkString) != 0 && !strings.Contains(content, checkString) {
			continue
		}

		links := getLinksFromContent(content)
		for _, link := range links {
			backLinksMap[link] = append(backLinksMap[link], fileName)
		}
		mdFileNameWithoutExtension := getMDFileNameWithoutExtension(fileName)
		fileNameToContentMap[mdFileNameWithoutExtension] = b
	}
	
	baseTemplateBytes, err := ioutil.ReadFile(baseTemplatePath)
	if (err != nil) {
		log.Fatal(err)
	}

	for mdFileNameWithoutExtension, byteContent := range fileNameToContentMap {
		byteOutput := markdown.ToHTML(byteContent, nil, nil)
		postContentOutput := string(byteOutput)

		baseTemplate := string(baseTemplateBytes)
		postPageOutput := strings.ReplaceAll(baseTemplate, "{{ content }}", postContentOutput)
		postPageOutput = strings.ReplaceAll(postPageOutput, "{{ title }}", mdFileNameWithoutExtension)
		links := getLinksFromContent(postPageOutput)
		for _, link := range links {
			markdownLink := fmt.Sprintf("[[%s]]", link)
			htmlLink := link
			if wow, ok := fileNameToContentMap[link]; ok {
				htmlLink = fmt.Sprintf("<a href='%s'>%s</a>", getHTMLFilePathFromFileName(link), link)
			}
			postPageOutput = strings.ReplaceAll(postPageOutput, markdownLink, htmlLink)
		}

		htmlFileName := strings.ToLower(strings.ReplaceAll(mdFileNameWithoutExtension, " ", "-")) + ".html"
		ouput_filepath := path.Join(outputDir, htmlFileName)

		err = ioutil.WriteFile(ouput_filepath, []byte(postPageOutput), 0644)
		if err != nil {
			log.Fatal(err)
		}
		processedFileCount += 1
	}

	listOfPages := "<ul>\n"
	for fileName := range fileNameToContentMap {
		htmlFilePath := getHTMLFilePathFromFileName(fileName)
		mdFileNameWithoutExtension := getMDFileNameWithoutExtension(fileName)
		row := fmt.Sprintf("<li><a href=\"%s\">%s</a></li>\n", htmlFilePath, mdFileNameWithoutExtension)
		listOfPages += row
	}
	listOfPages += "\n</ul>"

	baseTemplate := string(baseTemplateBytes)
	output := strings.ReplaceAll(baseTemplate, "{{ content }}", listOfPages)
	output = strings.ReplaceAll(output, "{{ title }}", "Home")
	ouput_filepath := path.Join(outputDir, "index.html")
	err = ioutil.WriteFile(ouput_filepath, []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(processedFileCount)
}

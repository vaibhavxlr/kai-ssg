package ssg

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"gopkg.in/yaml.v3"
)

func segregateFMandBody(fileData string) (string, string) {
	fileSli := strings.Split(fileData, "---")
	return fileSli[1], fileSli[2]
}

func deserializeContent(front, body, contType string) any {
	var siteMd SiteMeta
	err := yaml.Unmarshal([]byte(front), &siteMd)
	if err != nil {
		println("Failed to parse front matter", err)
		os.Exit(1)
	}

	//
	// switch type {
	// case "about":
	// case "blog":
	// case "post":
	// }

	htmlBody := markdown.ToHTML([]byte(body), nil, nil)
	var pageVals AboutPage
	pageVals.Content = template.HTML(htmlBody)
	pageVals.Author = "Vaibhav Upadhyay"
	pageVals.SiteTitle = "vaibhavxlr's blog"
	pageVals.Year = 2025
	return pageVals
}

func Generate(content, htmlTemplatesPath, output string) {
	contentList, err := os.ReadDir(content)
	if err != nil {
		println("[ERROR] unable to read the markdown content:", err)
		os.Exit(1)
	}

	tmpl := template.Must(template.ParseGlob(htmlTemplatesPath + "*.html"))

	for _, val := range contentList {
		if val.IsDir() {
		} else {
			if val.Name() == "about.md" {
				println("[INFO] Found about.md")
				fileBytes, err := os.ReadFile(filepath.Join(content, val.Name()))
				if err != nil {
					println("[ERROR] unable to read about.md:", err)
					os.Exit(0)
				}

				fileStr := string(fileBytes)

				front, body := segregateFMandBody(fileStr)
				println(front, body)
				pageCont := deserializeContent(front, body, "")

				outputFile, _ := os.Create(filepath.Join(output, "index.html"))
				tmpl.ExecuteTemplate(outputFile, "about", pageCont)
			}
		}
	}

}

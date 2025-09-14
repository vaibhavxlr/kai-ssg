package ssg

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"gopkg.in/yaml.v3"
)

func segFrontMatandBody(path string) (string, string) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("[ERROR] unable to read about.md: %v", err)
	}

	fileData := string(fileBytes)
	fileSli := strings.Split(fileData, "---")
	return fileSli[1], fileSli[2]
}

func deserializeContent(front, body, contType string) any {

	switch contType {
	case "about":

		var pageVals AboutPage
		err := yaml.Unmarshal([]byte(front), &pageVals)
		if err != nil {
			log.Fatal("Failed to parse front matter - 'about' level", err)
		}
		htmlBody := markdown.ToHTML([]byte(body), nil, nil)
		pageVals.Content = template.HTML(htmlBody)
		pageVals.Author = "Vaibhav Upadhyay"
		pageVals.SiteTitle = "vaibhavxlr's blog"
		pageVals.Year = time.Now().Year()
		fmt.Printf("struct - %v", pageVals)
		return pageVals

	case "blog":

		var pageVals BlogPage
		pageVals.Author = "Vaibhav Upadhyay"
		pageVals.SiteTitle = "vaibhavxlr's blog"
		pageVals.Year = time.Now().Year()
		pageVals.PostsByYear = nil
		fmt.Printf("struct - %v", pageVals)
		return pageVals

	case "post":

		var pageVals Post
		err := yaml.Unmarshal([]byte(front), &pageVals)
		if err != nil {
			log.Fatalf("Failed to parse front matter - 'post' level: %v", err)
		}
		htmlBody := markdown.ToHTML([]byte(body), nil, nil)
		pageVals.Year, _ = strconv.Atoi(strings.Split(pageVals.Date, "-")[0])
		pageVals.Content = template.HTML(htmlBody)
		fmt.Printf("struct - %v", pageVals)
		return pageVals

	}

	return nil
}

func Generate(content, htmlTemplatesPath, output string) {

	tmpl := template.Must(template.ParseGlob(htmlTemplatesPath + "*.html"))

	postbyYear := make(map[int][]Post)
	filepath.WalkDir(content, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err

		}

		if !d.IsDir() && d.Name() == "about.md" {
			log.Println("[INFO] Found about.md")
			front, body := segFrontMatandBody(path)
			pageCont := deserializeContent(front, body, "about")
			os.MkdirAll(output, 0755)
			outputFile, _ := os.Create(filepath.Join(output, "index.html"))
			err = tmpl.ExecuteTemplate(outputFile, "about", pageCont)
			if err != nil {
				log.Fatalf("[ERROR] unable to execute 'about' template: %v", err)
			}
		} else if !d.IsDir() && filepath.Ext(path) == ".md" {
			log.Println("[INFO] Found", d)
			front, body := segFrontMatandBody(path)
			pageCont := deserializeContent(front, body, "post").(Post)
			dir := filepath.Join(output, "posts", strconv.Itoa(pageCont.Year), pageCont.Slug)
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				log.Fatalf("[ERROR] unable to create 'post' dir: %v", err)
			}
			outputFile, _ := os.Create(filepath.Join(dir, "index.html"))
			tmpl.ExecuteTemplate(outputFile, "post", pageCont)

			postbyYear[pageCont.Year] = append(postbyYear[pageCont.Year], pageCont)
			srcDir := filepath.Dir(path)
			filesInPostDir, _ := os.ReadDir(srcDir)
			for _, val := range filesInPostDir {
				if !val.IsDir() && filepath.Ext(val.Name()) != ".md" {
					out, _ := os.Create(filepath.Join(dir, val.Name()))
					in, _ := os.Open(filepath.Join(srcDir, val.Name()))
					io.Copy(out, in)
				}
			}

		}

		return nil
	})
	dir := filepath.Join(output, "blog")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalf("[ERROR] unable to create 'blog' dir: %v", err)
	}
	outputFile, _ := os.Create(filepath.Join(output, "blog", "index.html"))

	pageCont := deserializeContent("", "", "blog").(BlogPage)

	fmt.Printf("\npageCont: %v", pageCont)
	pageCont.PostsByYear = postbyYear

	fmt.Printf("\npostbytyear2: %v", postbyYear)
	fmt.Printf("\npageCont2: %v", pageCont)
	err = tmpl.ExecuteTemplate(outputFile, "blog", pageCont)
	if err != nil {
		log.Fatalf("[ERROR] error executing 'blog' template %v", err)
	}

}

package main

import (
	"flag"
	"kai-ssg/ssg"
)

//
// type PageData struct {
// 	SiteTitle string
// 	Title     string
// 	Date      string
// 	Author    string
// 	Slug      string
// 	Year      int
// 	Content   template.HTML
// }

func main() {
	println("***--kai-ssg--***")

	content := flag.String("content", "./content/", "The place where markdowns live")

	htmlTemplates := flag.String("templates", "./templates/", "The place where templates live")

	output := flag.String("output", "./output/", "The place where output goes")
	flag.Parse()
	println("[INFO] flags: ", *content, *htmlTemplates, *output)

	ssg.Generate(*content, *htmlTemplates, *output)
}

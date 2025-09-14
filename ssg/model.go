package ssg

import "html/template"

type AboutPage struct {
	SiteTitle string
	Author    string
	Year      int
	Content   template.HTML
}

type Post struct {
	Title   string
	Date    string
	Slug    string
	Year    int
	Content template.HTML
}

type BlogPage struct {
	SiteTitle   string
	BlogIntro   string
	Author      string
	Year        int
	PostsByYear map[int][]Post
}

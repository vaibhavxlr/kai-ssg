package ssg

import "html/template"

type SiteMeta struct {
	SiteTitle string
	Author    string
	Year      int
	BlogIntro string
}

type Post struct {
	Title   string
	Date    string
	Year    int
	Slug    string
	Content template.HTML
}

type AboutPage struct {
	SiteMeta
	Content template.HTML
}

type BlogIndex struct {
	SiteMeta
	PostsByYear map[int][]Post
}

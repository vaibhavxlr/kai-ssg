package ssg

import "html/template"

type AboutPage struct {
	SiteTitle     string
	Author        string
	CopyrightYear int
	Content       template.HTML
}

type Post struct {
	SiteTitle     string
	Title         string
	Date          string
	Slug          string
	Year          int
	Content       template.HTML
	Author        string
	CopyrightYear int
}

type BlogPage struct {
	SiteTitle     string
	BlogIntro     string
	Author        string
	CopyrightYear int
	Years         []int
	PostsByYear   map[int][]Post
}

package main

import (
	"path/filepath"
)


type TemplateSet struct {
	ext string
	base string
	files []string
}

func (t TemplateSet) normalize() {
	if filepath.Ext(t.base) == "" {
		t.base += tplExt
	}

	for i, f := range t.files {
		if filepath.Ext(f) == "" {
			t.files[i] += tplExt
		}
	}
}

func (t TemplateSet) Layout() string {
	return t.base
}

func (t TemplateSet) Others() []string {
	return t.files
}

func newSet(base string, files ...string) TemplateSet {
	return TemplateSet{base: base, files: files}
}

// Sources isnt used. I've left it here as a comparison
type Sources struct{
	Base string
	Page string 
}
func (t Sources) Layout() string {
	return t.Base
}

func (t Sources) Others() []string {
	return []string{t.Page}
}
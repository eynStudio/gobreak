package io

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	. "github.com/eynstudio/gobreak"
)

type UrlStatus struct {
	Status
	Url string
}

type SaveYyyyMmFile struct {
	Error
	ext        []string
	rootFolder string
	file       multipart.File
	header     *multipart.FileHeader
	newName    string
	saveName   string
}

func NewSaveYyyyMmFile(root, newName string, ext []string) *SaveYyyyMmFile {
	return &SaveYyyyMmFile{rootFolder: root, newName: newName, ext: ext}
}

func (p *SaveYyyyMmFile) Save(r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		p.setUploadErr()
		return
	}
	if p.file, p.header, p.Err = r.FormFile("file"); p.IsErr() {
		p.setUploadErr()
		return
	}
	defer p.file.Close()
	p.CheckExt()
	p.NoErrExec(p.checkSaveFileName)

	var f *os.File
	if f, p.Err = os.OpenFile(p.saveName, os.O_WRONLY|os.O_CREATE, 0666); p.IsErr() {
		p.setUploadErr()
		return
	}
	defer f.Close()
	if _, p.Err = io.Copy(f, p.file); p.IsErr() {
		p.setUploadErr()
		return
	}
}

func (p *SaveYyyyMmFile) CheckExt() {
	if len(p.ext) == 0 {
		return
	}
	ext := filepath.Ext(p.header.Filename)
	ok := false
	for _, it := range p.ext {
		if it == ext {
			ok = true
			break
		}
	}
	p.SetErrfIf(!ok, "只允许上传%v格式文件", p.ext)
}

func (p *SaveYyyyMmFile) checkSaveFileName() {
	p.saveName = path.Join(p.rootFolder, time.Now().Format("2006/01/"))
	if p.Err = os.MkdirAll(p.saveName, 0666); p.IsErr() {
		p.setUploadErr()
		return
	}

	if p.newName != "" {
		p.saveName = path.Join(p.saveName, p.newName+filepath.Ext(p.header.Filename))
	} else {
		p.saveName = path.Join(p.saveName, p.header.Filename)
	}
}
func (p *SaveYyyyMmFile) GetFileName() string    { return p.header.Filename }
func (p *SaveYyyyMmFile) setUploadErr()          { p.SetErr("文件上传失败") }
func (p SaveYyyyMmFile) GetSaveFilePath() string { return "/" + p.saveName }

func (p SaveYyyyMmFile) GetUrlStatus() (m UrlStatus) {
	m.Status = p.GetStatus().GetStatus()
	if p.NotErr() {
		m.Url = p.GetSaveFilePath()
	}
	return
}

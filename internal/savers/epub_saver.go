package savers

import (
	"archive/zip"
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	textTemplate "text/template"

	"github.com/crispyarty/novelparser/internal"
	"github.com/crispyarty/novelparser/internal/config"
)

const templateDirPath = "templates/epub"
const uploadsDirPath = "uploads"

var funcMap = template.FuncMap{
	"inc": func(i int) int { return i + 1 },
}

var check = func(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func SaveNovel(name string, novels []*internal.NovelData) string {
	data := newContent(name, novels)
	zipFile, err := os.Create(archiveName(data))
	check(err)
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	copyFileToZip(zipWriter, "mimetype", zip.Store)
	copyFileToZip(zipWriter, "META-INF/container.xml", zip.Deflate)
	addFileToZip(zipWriter, "content.opf", genFromTextTemplate("content.opf", data), zip.Deflate)
	addFileToZip(zipWriter, "content.html", genFromHtmlTemplate("content.html", data), zip.Deflate)
	addFileToZip(zipWriter, "toc.ncx", genFromTextTemplate("toc.ncx", data), zip.Deflate)

	return zipFile.Name()
}

func templatePath(finename string) string {
	// return fmt.Sprintf("%v/%v", config.AssetPath(templateDirPath), finename)
	return config.AssetPath(templateDirPath, finename)
}

func archiveName(data *Content) string {
	// dir := fmt.Sprintf("uploads/%v", data.NovelName)
	dir := config.AssetPath(uploadsDirPath, data.NovelName)
	os.MkdirAll(dir, os.ModePerm)

	// return fmt.Sprintf("%v/%v.epub", dir, data.Title())
	return filepath.Join(dir, data.Title()+".epub")
}

func genFromHtmlTemplate(filename string, data *Content) []byte {
	buf := &bytes.Buffer{}
	tmpl, err := template.New(filename).Funcs(funcMap).ParseFiles(templatePath(filename))
	check(err)
	err = tmpl.Execute(buf, data)
	check(err)

	return buf.Bytes()
}

func genFromTextTemplate(filename string, data *Content) []byte {
	buf := &bytes.Buffer{}
	tmpl, err := textTemplate.New(filename).Funcs(funcMap).ParseFiles(templatePath(filename))
	check(err)
	err = tmpl.Execute(buf, data)
	check(err)

	return buf.Bytes()
}

func copyFileToZip(zipWriter *zip.Writer, filename string, compression uint16) {
	bytes, err := os.ReadFile(templatePath(filename))
	check(err)

	addFileToZip(zipWriter, filename, bytes, compression)
}

func addFileToZip(zipWriter *zip.Writer, filename string, content []byte, compression uint16) {
	header := &zip.FileHeader{
		Name:   filename,
		Method: compression,
	}

	fileWriter, err := zipWriter.CreateHeader(header)
	check(err)

	_, err = fileWriter.Write(content)
	check(err)
}

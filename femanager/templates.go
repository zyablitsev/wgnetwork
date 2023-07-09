package femanager

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed spa/dist/assets/index.js
//go:embed spa/dist/assets/index.css
var assets embed.FS

func loadAssets() ([]byte, []byte, error) {
	var (
		css, js []byte
		err     error
	)

	css, err = assets.ReadFile("spa/dist/assets/index.css")
	if err != nil {
		return nil, nil, err
	}

	js, err = assets.ReadFile("spa/dist/assets/index.js")
	if err != nil {
		return nil, nil, err
	}

	return css, js, nil
}

//go:embed spa/index.gohtml
var indexTmpl string

func loadIndex(apiUrl string) ([]byte, error) {
	tmpl, err := template.New("index.template").Parse(indexTmpl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	data := struct {
		ApiUrl string
	}{
		ApiUrl: apiUrl,
	}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

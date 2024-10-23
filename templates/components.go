package templates

import (
	"bytes"
	"html/template"
	text "text/template"
)

func Execute(name string, t string, data interface{}) string {
	tmpl, err := template.New(name).Parse(t)
	if err != nil {
		return err.Error()
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		return err.Error()
	}

	return b.String()
}

// The reason both functions are needed is because html/template sanitizes
// the html input, which is something we want, unless we already
// sanitized the input
func ExecuteText(name string, t string, data interface{}) string {
	tmpl, err := text.New(name).Parse(t)
	if err != nil {
		return err.Error()
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		return err.Error()
	}

	return b.String()
}

func Banner() string {
	return `
        <div class="header">
            <nav style="margin: 0px 10px">
                    <ul>
                            <li><a href="/credleak"><kbd>Cred Leak</kbd></a></li>
                            <li><a href="/openport"><kbd>Open Port</kbd></a></li>
                            <li><a href="/actor"><kbd>Actor</kbd></a></li>
                            <li><a href="/event"><kbd>Event</kbd></a></li>
                    </ul>
            </nav>
        </div>
    `
}

func header() string {
	return `
        <head>
            <title>JCTF Form Generator</title>
            <script src="https://unpkg.com/htmx.org@1.9.12" integrity="sha384-ujb1lZYygJmzgSwoxRggbCHcjc0rB2XoQrxeTUQyRjrOnlCoYta87iKBWq3EsdM2" crossorigin="anonymous"></script>
	    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.blue.min.css">
	    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.colors.min.css">
            <link rel="stylesheet" type="text/css" href="/style.css">
        </head>
        `
}

func BuildPage(body string) string {
	data := struct {
		Header string
		Body   string
		Banner string
	}{
		Header: header(),
		Body:   body,
		Banner: Banner(),
	}

	const page = `
        <!DOCTYPE html>
        <html lang="en">
        {{.Header}}
        <body hx-boost="true">
	    {{.Banner}}
            <div class="grid">
                <div></div>
                <div>{{.Body}}</div>
                <div></div>
            </div>
        </body>
        </html>
        `

	return ExecuteText("page", page, data)
}
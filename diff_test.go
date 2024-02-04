package diff_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
)

func red(s string) string {
	return "\x1b[101m\x1b[30m" + s + "\x1b[0m"
}

func green(s string) string {
	return "\x1b[102m\x1b[30m" + s + "\x1b[0m"
}

func redent(s string) string {
	return strings.TrimSpace(dedent.Dedent(s))
}

func colorize(text string) string {
	rep := strings.NewReplacer(
		"{red}", "\x1b[101m\x1b[30m",
		"{green}", "\x1b[102m\x1b[30m",
		"{reset}", "\x1b[0m",
	)
	return rep.Replace(text)
}

func TestDiffOk(t *testing.T) {
	is := is.New(t)
	err := diff.Diff("hi", "hi")
	is.NoErr(err)
}
func TestDeepOk(t *testing.T) {
	is := is.New(t)
	type C struct{ D string }
	type A struct{ *C }
	type B struct{}
	type Web struct {
		*A
		B
	}
	web1 := &Web{&A{&C{"D"}}, B{}}
	web2 := &Web{&A{&C{"D"}}, B{}}
	err := diff.Diff(web1, web2)
	is.NoErr(err)
}

func TestDiffNotOk(t *testing.T) {
	is := is.New(t)
	err := diff.Diff(3, "hi")
	is.True(err != nil)
	de, ok := err.(*diff.Error)
	is.True(ok)
	is.Equal(red("3")+green("hi"), de.Diff)
}
func TestDeepNotOk(t *testing.T) {
	is := is.New(t)
	type C struct{ D string }
	type A struct{ *C }
	type B struct{}
	type Web struct {
		*A
		B
	}
	web1 := &Web{&A{&C{"D"}}, B{}}
	web2 := &Web{&A{&C{"F"}}, B{}}
	err := diff.Diff(web1, web2)
	is.True(err != nil)
	de, ok := err.(*diff.Error)
	is.True(ok)
	is.Equal(de.Diff, redent(`
		&Web{A: &A{
			C: &C{D: "`+red("D")+green("F")+`"},
		}}
	`))
}

func TestStringOk(t *testing.T) {
	is := is.New(t)
	err := diff.String("hi", "hi")
	is.NoErr(err)
}

func TestStringNotOk(t *testing.T) {
	is := is.New(t)
	err := diff.String("hi", "cool")
	is.True(err != nil)
	de, ok := err.(*diff.Error)
	is.True(ok)
	is.Equal(red("hi")+green("cool"), de.Diff)
}

func TestHTTPOk(t *testing.T) {
	is := is.New(t)
	a := `
		HTTP/1.1 200 OK
		Connection: close
	`
	err := diff.HTTP(a, a)
	is.NoErr(err)
}

func TestHTTPRequest(t *testing.T) {
	is := is.New(t)
	req := httptest.NewRequest("POST", "http://example.com", bytes.NewBufferString(`{"hello": "world"}`))
	out, err := httputil.DumpRequestOut(req, true)
	is.NoErr(err)
	diff.TestHTTP(t, string(out), string(out))
	diff.TestHTTP(t, string(out), `
		POST / HTTP/1.1
		Host: example.com
		User-Agent: Go-http-client/1.1
		Content-Length: 18
		Accept-Encoding: gzip

		{"hello": "world"}
	`)
}

func TestHTTPResponse(t *testing.T) {
	is := is.New(t)
	req := httptest.NewRequest("POST", "http://example.com", bytes.NewBufferString(`{"hello": "world"}`))
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"hello": "world"}`))
	})
	handler.ServeHTTP(rec, req)
	out, err := httputil.DumpResponse(rec.Result(), true)
	is.NoErr(err)
	diff.TestHTTP(t, string(out), string(out))
	diff.TestHTTP(t, string(out), `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"hello": "world"}
	`)
}

func TestHTTPNotOk(t *testing.T) {
	is := is.New(t)
	a := `
		HTTP/1.1 200 OK
		Connection: close
	`
	b := `
		HTTP/1.1 404 Not Found
		Connection: close
	`
	err := diff.HTTP(a, b)
	is.True(err != nil)
	de, ok := err.(*diff.Error)
	is.True(ok)
	is.Equal(de.Actual, "HTTP/1.1 200 OK\nConnection: close")
	is.Equal(de.Expect, "HTTP/1.1 404 Not Found\nConnection: close")
	is.Equal(de.Diff, colorize("HTTP/1.1 {red}200·OK{reset}{green}404·Not·Found{reset}\nConnection: close"))
}

func TestHTMLOK(t *testing.T) {
	is := is.New(t)
	a := `
		HTTP/1.1 200 OK
		Content-Length: 96
		Content-Type: text/html

		<html>
			<head>
				<title>Duo</title>
			</head>
			<body>
				<h1>Hello Berlin!</h1>
			</body>
		</html>
	`
	b := `
		HTTP/1.1 200 OK
		Content-Length: 96
		Content-Type: text/html

		<html>
			<head>
				<title>Duo</title>
			</head>
			<body>
				<h1>Hello Berlin!</h1>
			</body>
		</html>
	`
	err := diff.HTTP(a, b)
	is.NoErr(err)
}

func TestHTMLNotOK(t *testing.T) {
	is := is.New(t)
	a := `
		HTTP/1.1 200 OK
		Content-Length: 96
		Content-Type: text/html

		<html>
			<head>
				<title>Duo</title>
			</head>
			<body>
				<h1>Hello Berlin!</h1>
			</body>
		</html>
	`
	b := `
		HTTP/1.1 200 OK
		Content-Length: 96
		Content-Type: text/html


		<html>
			<head>
				<title>Duo</title>
			</head>
			<body>
				<h1>Hello Berlin!</h1>
			</body>
		</html>
	`
	err := diff.HTTP(a, b)
	is.True(err != nil)
	de, ok := err.(*diff.Error)
	is.True(ok)
	// Just check for the extra newline that's green
	is.True(strings.Contains(de.Diff, "\x1b[102m\x1b[30m\\n\x1b[0m"))
}

func TestContentOK(t *testing.T) {
	is := is.New(t)
	const a = `
		a
		b
		c
	`
	const b = `
	a
	b
	c
	`
	err := diff.Content(a, b)
	is.NoErr(err)
}

func TestContentNotOK(t *testing.T) {
	is := is.New(t)
	const a = `
		a
		b
		c
		d
	`
	const b = `
	a
	b
	c
	`
	err := diff.Content(a, b)
	is.True(err != nil)
	de, ok := err.(*diff.Error)
	is.True(ok)
	is.Equal(de.Diff, colorize("a\nb\nc{red}\\nd{reset}"))
}

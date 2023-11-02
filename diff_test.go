package diff_test

import (
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

func TestDiffOk(t *testing.T) {
	is := is.New(t)
	result := diff.Diff("hi", "hi")
	// No color means no difference
	is.Equal(`hi`, result)
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
	result := diff.Diff(web1, web2)
	// No color means no difference
	is.Equal(result, redent(`
		&Web{A: &A{
			C: &C{D: "D"},
		}}
	`))
}

func TestDiffNotOk(t *testing.T) {
	is := is.New(t)
	result := diff.Diff(3, "hi")
	// Color means difference
	is.Equal(red("3")+green("hi"), result)
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
	result := diff.Diff(web1, web2)
	// Color means difference
	is.Equal(result, redent(`
		&Web{A: &A{
			C: &C{D: "`+red("D")+green("F")+`"},
		}}
	`))
}

func TestStringOk(t *testing.T) {
	is := is.New(t)
	result := diff.String("hi", "hi")
	is.Equal("hi", result)
}

func TestStringNotOk(t *testing.T) {
	is := is.New(t)
	result := diff.String("hi", "cool")
	is.Equal("\x1b[101m\x1b[30mhi\x1b[0m\x1b[102m\x1b[30mcool\x1b[0m", result)
}

func TestHTTPOk(t *testing.T) {
	is := is.New(t)
	a := `
		HTTP/1.1 200 OK
		Connection: close
	`
	result := diff.HTTP(a, a)
	is.Equal(result, diff.HTTP(a, result))
}

func TestHTTPNotOk(t *testing.T) {
	is := is.New(t)
	a := `
		HTTP/1.1 200 OK
		Connection: close
	`
	b := `
		HTTP/1.1 404 Not FOund
		Connection: close
	`
	result := diff.HTTP(a, b)
	is.Equal("HTTP/1.1 \x1b[101m\x1b[30m2\x1b[0m\x1b[102m\x1b[30m4\x1b[0m0\x1b[101m\x1b[30m0\x1b[0m\x1b[102m\x1b[30m4 Not\x1b[0m \x1b[102m\x1b[30mF\x1b[0mO\x1b[101m\x1b[30mK\x1b[0m\x1b[102m\x1b[30mund\x1b[0m\nConnection: close", result)
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
	result := diff.HTTP(a, b)
	is.Equal(result, diff.HTTP(a, result))
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
	result := diff.HTTP(a, b)
	is.Equal("HTTP/1.1 200 OK\nContent-Length: 96\nContent-Type: text/html\n\n\x1b[102m\x1b[30m\n\x1b[0m<html>\n\t<head>\n\t\t<title>Duo</title>\n\t</head>\n\t<body>\n\t\t<h1>Hello Berlin!</h1>\n\t</body>\n</html>", result)
}

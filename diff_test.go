package diff_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matthewmueller/diff"
)

func TestDiffOk(t *testing.T) {
	result := diff.Diff("hi", "hi")
	assert.Equal(t, "hi", result)
}

func TestDiffNotOk(t *testing.T) {
	result := diff.Diff(3, "hi")
	assert.Equal(t, "\x1b[102m\x1b[30mh\x1b[0mi\x1b[101m\x1b[30mnt(3)\x1b[0m", result)
}

func TestStringOk(t *testing.T) {
	result := diff.String("hi", "hi")
	assert.Equal(t, "hi", result)
}

func TestStringNotOk(t *testing.T) {
	result := diff.String("hi", "cool")
	assert.Equal(t, "\x1b[101m\x1b[30mhi\x1b[0m\x1b[102m\x1b[30mcool\x1b[0m", result)
}

func TestHTTPOk(t *testing.T) {
	a := `
		HTTP/1.1 200 OK
		Connection: close
	`
	result := diff.HTTP(a, a)
	// Note, this isn't perfect, but visually it's not a problem.
	assert.Equal(t, "HTTP/1.1 200 OK\x1b[101m\x1b[30m\r\x1b[0m\nConnection: close", result)
}

func TestHTTPNotOk(t *testing.T) {
	a := `
		HTTP/1.1 200 OK
		Connection: close
	`
	b := `
		HTTP/1.1 404 Not FOund
		Connection: close
	`
	result := diff.HTTP(a, b)
	assert.Equal(t, "HTTP/1.1 \x1b[101m\x1b[30m20\x1b[0m\x1b[102m\x1b[30m4\x1b[0m0\x1b[102m\x1b[30m4\x1b[0m \x1b[102m\x1b[30mNot F\x1b[0mO\x1b[101m\x1b[30mK\r\x1b[0m\x1b[102m\x1b[30mund\x1b[0m\nConnection: close", result)
}

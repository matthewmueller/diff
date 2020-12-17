package diff

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/kr/pretty"
	"github.com/lithammer/dedent"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Diff two interfaces
func Diff(a, b interface{}) string {
	as := pretty.Sprint(a)
	bs := pretty.Sprint(b)
	return String(as, bs)
}

// String diffs two strings
func String(a, b string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(a, b, false)
	var buf bytes.Buffer
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			buf.WriteString("\x1b[102m\x1b[30m")
			buf.WriteString(diff.Text)
			buf.WriteString("\x1b[0m")
		case diffmatchpatch.DiffDelete:
			buf.WriteString("\x1b[101m\x1b[30m")
			buf.WriteString(diff.Text)
			buf.WriteString("\x1b[0m")
		case diffmatchpatch.DiffEqual:
			buf.WriteString(diff.Text)
		}
	}
	result := buf.String()
	result = strings.Replace(result, "\\n", "\n", -1)
	result = strings.Replace(result, "\\t", "\t", -1)
	return result
}

// HTTP diffs two response dumps via httputil.DumpResponse
func HTTP(a, b string) string {
	a = strings.ReplaceAll(strings.TrimSpace(dedent.Dedent(a)), "\r\n", "\n")
	b = strings.ReplaceAll(strings.TrimSpace(dedent.Dedent(b)), "\r\n", "\n")
	return String(a, b)
}

// Test tests two diffs
func Test(t testing.TB, expected interface{}, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		return
	}
	exp := pretty.Sprint(expected)
	act := pretty.Sprint(actual)
	var b bytes.Buffer
	b.WriteString("\n\x1b[4mExpected\x1b[0m:\n")
	b.WriteString(exp)
	b.WriteString("\n\n")
	b.WriteString("\x1b[4mActual\x1b[0m: \n")
	b.WriteString(act)
	b.WriteString("\n\n")
	b.WriteString("\x1b[4mDifference\x1b[0m: \n")
	b.WriteString(String(exp, act))
	b.WriteString("\n")
	t.Fatal(b.String())
}

// TestHTTP diffs two HTTP dumps from httputil.DumpResponse
func TestHTTP(t testing.TB, expected, actual string) {
	expected = strings.Replace(strings.TrimSpace(dedent.Dedent(expected)), "\n", "\r\n", -1)
	actual = strings.TrimSpace(dedent.Dedent(actual))
	Test(t, expected, actual)
}

// TestString ignores
func TestString(t testing.TB, expected string, actual string) {
	if expected == actual {
		return
	}
	var b bytes.Buffer
	b.WriteString("\n\x1b[4mExpected\x1b[0m:\n")
	b.WriteString(expected)
	b.WriteString("\n\n")
	b.WriteString("\x1b[4mActual\x1b[0m: \n")
	b.WriteString(actual)
	b.WriteString("\n\n")
	b.WriteString("\x1b[4mDifference\x1b[0m: \n")
	b.WriteString(String(expected, actual))
	b.WriteString("\n")
	t.Fatal(b.String())
}

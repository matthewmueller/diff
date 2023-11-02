package diff

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hexops/valast"
	"github.com/lithammer/dedent"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Diff two interfaces
func Diff(actual, expect interface{}) string {
	return String(format(actual), format(expect))
}

// String diffs two strings
func String(actual, expect string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(actual, expect, false)
	var buf bytes.Buffer
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			buf.WriteString(green(diff.Text))
		case diffmatchpatch.DiffDelete:
			buf.WriteString(red(diff.Text))
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

// Test tests actual with expected
func Test(t testing.TB, actual, expect interface{}) {
	t.Helper()
	TestString(t, format(actual), format(expect))
}

// TestHTTP diffs two HTTP dumps from httputil.DumpResponse
func TestHTTP(t testing.TB, actual, expect string) {
	t.Helper()
	actual = strings.ReplaceAll(strings.TrimSpace(dedent.Dedent(actual)), "\r\n", "\n")
	expect = strings.ReplaceAll(strings.TrimSpace(dedent.Dedent(expect)), "\r\n", "\n")
	TestString(t, expect, actual)
}

// Report the differences between actual and expect
func Report(actual, expect string) string {
	if actual == expect {
		return ""
	}
	s := new(strings.Builder)
	s.WriteString("\n\x1b[4mExpect\x1b[0m:\n")
	s.WriteString(expect)
	s.WriteString("\n\n")
	s.WriteString("\x1b[4mActual\x1b[0m: \n")
	s.WriteString(actual)
	s.WriteString("\n\n")
	s.WriteString("\x1b[4mDifference\x1b[0m: \n")
	s.WriteString(String(actual, expect))
	s.WriteString("\n")
	return s.String()
}

// TestString diffs two strings
func TestString(t testing.TB, actual string, expect string) {
	t.Helper()
	report := Report(actual, expect)
	if report == "" {
		return
	}
	t.Fatal(report)
}

func format(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return valast.StringWithOptions(v, &valast.Options{
		Unqualify: true,
	})
}

func red(s string) string {
	return "\x1b[101m\x1b[30m" + s + "\x1b[0m"
}

func green(s string) string {
	return "\x1b[102m\x1b[30m" + s + "\x1b[0m"
}

func bold(s string) string {
	return "\n\x1b[4m" + s + "\x1b[0m"
}

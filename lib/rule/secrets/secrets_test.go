package secrets_test

import (
	"testing"
	"testing/fstest"

	"github.com/duythinht/tg/lib/rule/secrets"
)

func TestGoSourceCheck(t *testing.T) {
	src := `
	package main

	func main() {
		s := "private_key xyz"
	}
	`

	fs := fstest.MapFS{
		"cmd/main.go": {
			Data: []byte(src),
		},
	}

	f, err := fs.Open("cmd/main.go")

	if err != nil {
		t.Fatalf("could not open mock file: %s", err)
	}

	r := &secrets.Rule{}

	reports, err := r.Check(f)

	if err != nil {
		t.Fatalf("error on rule check: %s", err)
	}

	if len(reports) < 1 {
		t.Logf("should have reports exists")
		t.Fail()
	}

	switch {
	case reports[0].Location.Path != "main.go":
		t.Logf("path should be cmd/main.go")
		t.Fail()
	case reports[0].Type != secrets.Type:
		t.Logf("report type should be %s, got: %s", secrets.Type, reports[0].Type)
		t.Fail()

	case reports[0].Location.Positions.Begin.Line != 5:
		t.Logf("report begin line should be %d, got: %d", 5, reports[0].Location.Positions.Begin.Line)
		t.Fail()
	}

}

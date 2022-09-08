package secrets

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"io/fs"
	"strings"

	"github.com/duythinht/tg/lib/rule"
)

const (
	ID   = "G143"
	Type = "sast"
)

type Rule struct{}

// Name of secrets rule
func (r *Rule) Name() string {
	return "secrets"
}

// Check scan secrets for a file
// this version currently support only check secret on go source code
func (r *Rule) Check(f fs.File) ([]rule.Report, error) {

	stat, err := f.Stat()

	if err != nil {
		return nil, fmt.Errorf("error reading file stats %w", err)
	}

	src, err := io.ReadAll(f)

	if err != nil {
		return nil, fmt.Errorf("error reading file content: %s - %w", stat.Name(), err)
	}

	switch {
	case strings.HasSuffix(stat.Name(), ".go"):
		return r.checkGoSource(stat.Name(), src)
	}

	return []rule.Report{}, nil
}

func (r *Rule) checkGoSource(name string, src []byte) ([]rule.Report, error) {
	// Initialize the scanner.
	var s scanner.Scanner
	fset := token.NewFileSet()                        // positions are relative to fset
	file := fset.AddFile(name, fset.Base(), len(src)) // register input "file"
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

	reports := make([]rule.Report, 0)

	for pos, tok, lit := s.Scan(); tok != token.EOF; pos, tok, lit = s.Scan() {
		if tok == token.STRING {

			// trim string/raw string from lit token
			// eg:
			//		"private_key xyz" => private_key xyz
			//      `public_key xyz`  => public_key xyz
			stringValue := lit[1 : len(lit)-1]

			if !strings.HasPrefix(stringValue, "private_key") && !strings.HasPrefix(stringValue, "public_key") {
				continue
			}

			reports = append(reports, rule.Report{
				Type:   Type,
				RuleID: ID,
				Location: rule.Location{
					Positions: rule.Positions{
						Begin: rule.Begin{
							Line: fset.Position(pos).Line,
						},
					},
				},
				Metadata: rule.Metadata{
					Description: "secrets found the the source code",
					Severity:    "HIGH",
				},
			})
		}
	}
	return reports, nil
}

var _ = rule.Rule(&Rule{})

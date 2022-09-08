package rule

import "io/fs"

// Rule signature interface, every rule should be implement to easily extensible
type Rule interface {
	Check(f fs.File) ([]Report, error)
	Name() string
}

// Report structure
type Report struct {
	Type     string   `json:"type"`
	RuleID   string   `json:"ruleId"`
	Location Location `json:"location"`
	Metadata Metadata `json:"metadata"`
}

// Location of a report
type Location struct {
	Path      string    `json:"path"`
	Positions Positions `json:"positions"`
}

// Position on a file
type Positions struct {
	Begin Begin `json:"begin"`
}

// Begin of issue
type Begin struct {
	Line int `json:"line"`
}

// Metadata of a report
type Metadata struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

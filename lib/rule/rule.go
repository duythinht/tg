package rule

import "io/fs"

type Report struct {
	Type     string   `json:"type"`
	RuleID   string   `json:"ruleId"`
	Location Location `json:"location"`
	Metadata metadata `json:"metadata"`
}

type Location struct {
	Path      string    `json:"path"`
	Positions Positions `json:"positions"`
}

type Positions struct {
	Begin Begin `json:"begin"`
}

type Begin struct {
	Line int `json:"line"`
}

type metadata struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type Rule interface {
	Check(f fs.File) ([]Report, error)
}

package fieldparser

import (
	"errors"

	"github.com/caelifer/gotests/json2/ranker"
)

// Errors
var (
	NoMatchingParserError = errors.New("No matching parser found")
)

type Meta interface {
	ID() string
	Type() string
	Items() string
	Custom() string
}

// MetaTag a struct used to de-serialize JSON field schema, implements Meta interface
type Schema struct {
	ID     string // Converted field
	Type   string `json:"type"`
	Items  string `json:"items"`
	Custom string `json:"custom"`
}

// FieldParser - an interface for custom Jira field parsers
type FieldParser interface {
	Score(interface{}) ranker.MatchScore
	Parse(jsonChunk []byte) (string, error)
}

// ParseJiraField
func ParseJiraField(meta Meta, jsn []byte) (string, error) {
	fp := ranker.New(meta).BestMatch(toScorer(rfp)...)
	if fp == nil {
		return "", NoMatchingParserError
	}

	return fp.(FieldParser).Parse(jsn) // Assert correct type
}

// RegisterParser
func RegisterParser(fp FieldParser) {
	rfp = append(rfp, fp) // Register FieldParser
}

// rfp - a package global to store all registered field parsers
var rfp []FieldParser

// Helper converter
func toScorer(p []FieldParser) []ranker.Scorer {
	s := make([]ranker.Scorer, len(p))
	for i, v := range p {
		s[i] = v
	}
	return s
}

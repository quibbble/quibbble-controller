package gamenotation

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"text/scanner"
)

func Parse(raw string) (*Snapshot, error) {
	snapshot := &Snapshot{Tags: Tags{}, Actions: []Action{}}
	r := strings.NewReader(raw)
	sc := &scanner.Scanner{}
	sc.Init(r)
	// this uses the GoTokens scanner defaults other than ScanFloats
	sc.Mode = scanner.ScanIdents | scanner.ScanChars | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments | scanner.SkipComments
	err := parseTags(sc, snapshot)
	if err != nil {
		return nil, err
	}
	for _, required := range RequiredTags {
		if snapshot.Tags[required] == "" {
			return nil, fmt.Errorf("missing required %s tag", required)
		}
	}
	err = parseActions(sc, snapshot)
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

func parseTags(s *scanner.Scanner, snapshot *Snapshot) error {
	run := s.Peek()
	inside := false
	for run != scanner.EOF {
		switch run {
		case '[':
			if inside {
				return fmt.Errorf("cannot nest right bracket for tags")
			}
			inside = true
			_ = s.Next()
		case ']':
			if !inside {
				return fmt.Errorf("missing starting right bracket for tags")
			}
			inside = false
			_ = s.Next()
		case '\n', '\r', '\t', ' ':
			_ = s.Next()
		default:
			if !inside {
				return nil
			}
			s.Scan()
			tag := s.TokenText()
			s.Scan()
			val := s.TokenText()
			snapshot.Tags[tag] = strings.Trim(val, "\"")
		}
		run = s.Peek()
	}
	return nil
}

func parseActions(s *scanner.Scanner, snapshot *Snapshot) error {
	run := s.Peek()
	var action *Action
	for run != scanner.EOF {
		switch run {
		case ' ':
			_ = s.Next()
		case '{':
			for run != '}' {
				run = s.Next()
				if run == scanner.EOF {
					return fmt.Errorf("missing comment closure")
				}
			}
		default:
			s.Scan()
			if s.TokenText() == "&" {
				if action == nil {
					return fmt.Errorf("cannot have action details before base")
				}
				s.Scan()
				details := s.TokenText()
				ignore := []rune{' ', '\n', '\t', scanner.EOF}
				for !slices.Contains(ignore, s.Peek()) {
					s.Scan()
					if s.TokenText() == "&" {
						return fmt.Errorf("multiple ampersands found in action")
					}
					details += s.TokenText()
				}
				split := strings.Split(details, ".")
				action.Details = split
				snapshot.Actions = append(snapshot.Actions, *action)
				action = nil
			} else {
				if action != nil {
					snapshot.Actions = append(snapshot.Actions, *action)
				}
				base := s.TokenText()
				ignore := []rune{' ', '\n', '\t', '&', scanner.EOF}
				for !slices.Contains(ignore, s.Peek()) {
					s.Scan()
					base += s.TokenText()
				}
				if base == "" && s.Peek() == scanner.EOF {
					return nil
				}
				if len(base) < 2 {
					return fmt.Errorf("invalid action base")
				}
				team, err := strconv.Atoi(base[:len(base)-1])
				if err != nil {
					return err
				}
				action = &Action{
					Index: team,
					Key:   string(base[len(base)-1]),
				}
			}
		}
		run = s.Peek()
		if run == scanner.EOF && action != nil {
			snapshot.Actions = append(snapshot.Actions, *action)
		}
	}
	return nil
}

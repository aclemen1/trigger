// Package trigger parses the keys that say when something wakes up, and
// matches a fired key against them.
//
//	at:2026-10-15T09:00          a date
//	rrule:FREQ=WEEKLY;BYDAY=MO   a recurrence
//	run:<name>                   a recurrence computed by a command
//	on:<tool>:<event>[:<id>]     an event of a tool
//	with:<what>[:<id>]           a situation: with:contact:JMR, with:home
//	after:<tool>:<id>            the end of something else
//
// A key is a list of segments separated by ':'. A trigger matches a fired key
// when its segments start the fired key's: on:mail:reply matches
// on:mail:reply:thread/abc, on:mail:rep does not.
package trigger

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	At    = "at"
	RRule = "rrule"
	Run   = "run"
	On    = "on"
	With  = "with"
	After = "after"
)

// Kinds are every kind of key, in the order of the documentation.
var Kinds = []string{At, RRule, Run, On, With, After}

// Key is a parsed trigger key.
type Key struct {
	kind     string
	segments []string
}

var (
	dateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(T\d{2}:\d{2})?$`)
	nameRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
)

// Parse checks a key and returns it. The error says what form was expected.
func Parse(s string) (Key, error) {
	s = strings.TrimSpace(s)
	kind, rest, ok := strings.Cut(s, ":")
	if !ok || rest == "" {
		return Key{}, fmt.Errorf("trigger %q: expected <kind>:<value>, kinds %s. Example: on:mail:reply:thread/abc", s, strings.Join(Kinds, ", "))
	}
	switch kind {
	case At:
		if !dateRe.MatchString(rest) {
			return Key{}, fmt.Errorf("trigger %q: expected at:2026-10-15 or at:2026-10-15T09:00", s)
		}
		return Key{kind, []string{rest}}, nil
	case RRule:
		if !strings.Contains(strings.ToUpper(rest), "FREQ=") {
			return Key{}, fmt.Errorf("trigger %q: expected rrule:FREQ=…, e.g. rrule:FREQ=WEEKLY;BYDAY=MO", s)
		}
		return Key{kind, []string{rest}}, nil
	case Run:
		if !nameRe.MatchString(rest) {
			return Key{}, fmt.Errorf("trigger %q: expected run:<name>, e.g. run:school-homework", s)
		}
		return Key{kind, []string{rest}}, nil
	}
	segs := strings.Split(rest, ":")
	for _, g := range segs {
		if strings.TrimSpace(g) == "" || g != strings.TrimSpace(g) {
			return Key{}, fmt.Errorf("trigger %q: a segment between ':' is empty or padded with spaces", s)
		}
	}
	switch kind {
	case On:
		if len(segs) < 2 {
			return Key{}, fmt.Errorf("trigger %q: expected on:<tool>:<event>[:<id>], e.g. on:mail:reply:thread/abc", s)
		}
	case With:
	case After:
		if len(segs) < 2 {
			return Key{}, fmt.Errorf("trigger %q: expected after:<tool>:<id>, e.g. after:ticket:T-12", s)
		}
	default:
		return Key{}, fmt.Errorf("trigger %q: unknown kind %q; kinds %s", s, kind, strings.Join(Kinds, ", "))
	}
	return Key{kind, segs}, nil
}

// MustParse is Parse for keys known to be valid; it panics otherwise.
func MustParse(s string) Key {
	k, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return k
}

func (k Key) Kind() string { return k.kind }

// Segments are the parts after the kind.
func (k Key) Segments() []string { return append([]string(nil), k.segments...) }

func (k Key) String() string {
	if k.kind == "" {
		return ""
	}
	return k.kind + ":" + strings.Join(k.segments, ":")
}

func (k Key) IsZero() bool { return k.kind == "" }

// Starts says whether the segments of k start those of fired, with the same kind.
func (k Key) Starts(fired Key) bool {
	if k.kind == "" || k.kind != fired.kind || len(k.segments) > len(fired.segments) {
		return false
	}
	for i, s := range k.segments {
		if fired.segments[i] != s {
			return false
		}
	}
	return true
}

// Matches says whether a trigger matches a fired key; either one invalid never matches.
func Matches(trigger, fired string) bool {
	t, err := Parse(trigger)
	if err != nil {
		return false
	}
	f, err := Parse(fired)
	if err != nil {
		return false
	}
	return t.Starts(f)
}

// HasPrefix says whether a key starts with a prefix given as text, by whole
// segments: "with:" and "with:contact" match with:contact:JMR, "with:cont" does not.
func HasPrefix(key, prefix string) bool {
	prefix = strings.TrimSuffix(strings.TrimSpace(prefix), ":")
	if prefix == "" {
		return true
	}
	return key == prefix || strings.HasPrefix(key, prefix+":")
}

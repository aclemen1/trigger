package trigger

import "testing"

func TestParse(t *testing.T) {
	valid := map[string]string{
		"at:2026-10-15":               At,
		"at:2026-10-15T09:00":         At,
		"rrule:FREQ=DAILY;INTERVAL=2": RRule,
		"run:school-homework":         Run,
		"on:mail:reply":               On,
		"on:mail:reply:thread/abc":    On,
		"on:task:done:PT-0012":        On,
		"on:oj:RDIR":                  On,
		"with:contact:JMR":            With,
		"with:home":                   With,
		"with:place:La Tour-de-Peilz": With,
		"after:task:PT-0012":          After,
		"  on:mail:reply  ":           On,
		"on:cal:event:abc:def":        On,
	}
	for s, kind := range valid {
		k, err := Parse(s)
		if err != nil {
			t.Errorf("%q: %v", s, err)
			continue
		}
		if k.Kind() != kind {
			t.Errorf("%q: kind %s, want %s", s, k.Kind(), kind)
		}
	}
	invalid := []string{
		"", "on", "on:", "on:mail", "on:mail:", "on::reply", "on:mail: reply",
		"at:2026-10", "at:15.10.2026", "rrule:BYDAY=MO", "run:", "run:a b",
		"after:task", "with:", "when:contact:JMR", "ON:mail:reply",
	}
	for _, s := range invalid {
		if _, err := Parse(s); err == nil {
			t.Errorf("%q: accepted", s)
		}
	}
}

func TestMatches(t *testing.T) {
	cases := []struct {
		trigger, fired string
		want           bool
	}{
		{"on:mail:reply", "on:mail:reply:thread/abc", true},
		{"on:mail:reply:thread/abc", "on:mail:reply:thread/abc", true},
		{"on:mail:reply:thread/abc", "on:mail:reply:thread/xyz", false},
		{"on:mail:reply:thread/abc", "on:mail:reply", false},
		{"on:mail:rep", "on:mail:reply", false},
		{"on:task:done", "on:task:done:PT-0012", true},
		{"with:contact:JMR", "with:contact:JMR", true},
		{"with:contact:JMR", "with:contact:jmr", false},
		{"with:contact", "with:contact:JMR", true},
		{"with:contact:JMR", "on:contact:JMR", false},
		{"after:task:PT-0012", "after:task:PT-0012", true},
		{"bad", "on:mail:reply", false},
		{"on:mail:reply", "bad", false},
	}
	for _, c := range cases {
		if got := Matches(c.trigger, c.fired); got != c.want {
			t.Errorf("Matches(%q, %q) = %v, want %v", c.trigger, c.fired, got, c.want)
		}
	}
}

func TestHasPrefix(t *testing.T) {
	cases := []struct {
		key, prefix string
		want        bool
	}{
		{"with:contact:JMR", "with:", true},
		{"with:contact:JMR", "with", true},
		{"with:contact:JMR", "with:contact", true},
		{"with:contact:JMR", "with:cont", false},
		{"on:mail:reply", "with:", false},
		{"on:mail:reply", "", true},
	}
	for _, c := range cases {
		if got := HasPrefix(c.key, c.prefix); got != c.want {
			t.Errorf("HasPrefix(%q, %q) = %v, want %v", c.key, c.prefix, got, c.want)
		}
	}
}

func TestString(t *testing.T) {
	if s := MustParse(" on:mail:reply:thread/abc ").String(); s != "on:mail:reply:thread/abc" {
		t.Errorf("String: %q", s)
	}
}

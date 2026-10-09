# trigger

Keys that say when something wakes up, shared by the tools that keep tasks and
rules, and the rule that matches a fired key against them.

```go
import "github.com/aclemen1/trigger"

k, err := trigger.Parse("on:mail:reply:thread/abc") // k.Kind() == trigger.On
trigger.Matches("on:mail:reply", "on:mail:reply:thread/abc") // true
trigger.HasPrefix("with:contact:JMR", "with:")               // true
```

| Kind | Form | Example |
|---|---|---|
| `at` | `at:<date>[T<hh:mm>]` | `at:2026-10-15T09:00` |
| `rrule` | `rrule:<RRULE>` | `rrule:FREQ=WEEKLY;BYDAY=MO` |
| `run` | `run:<name>` | `run:school-homework` |
| `on` | `on:<tool>:<event>[:<id>]` | `on:mail:reply:thread/abc` |
| `with` | `with:<what>[:<id>]` | `with:contact:JMR` |
| `after` | `after:<tool>:<id>` | `after:ticket:T-12` |

- Segments are separated by `:`; none is empty or padded with spaces.
- A trigger matches a fired key when it has the same kind and its segments
  start the key's: `on:mail:reply` matches `on:mail:reply:thread/abc`,
  `on:mail:rep` does not. Matching is case-sensitive.
- `HasPrefix` filters keys by whole segments: `with:` and `with:contact` match
  `with:contact:JMR`, `with:cont` does not.

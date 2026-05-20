# logslice

Fast log file parser with time-range filtering and structured output for common formats.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

## Usage

Parse a log file and filter entries within a specific time range:

```bash
logslice --file /var/log/app.log --from "2024-01-15 08:00:00" --to "2024-01-15 09:00:00"
```

Output results as JSON:

```bash
logslice --file /var/log/nginx/access.log --from "2024-01-15 08:00:00" --format json
```

Use as a library:

```go
package main

import "github.com/yourusername/logslice"

func main() {
    entries, err := logslice.Parse("app.log", logslice.Options{
        From:   "2024-01-15 08:00:00",
        To:     "2024-01-15 09:00:00",
        Format: logslice.JSON,
    })
    if err != nil {
        panic(err)
    }
    for _, e := range entries {
        fmt.Println(e)
    }
}
```

## Supported Formats

- Nginx / Apache access logs
- Syslog
- JSON structured logs
- Common timestamp patterns (RFC3339, ISO8601, Unix)

## License

MIT — see [LICENSE](LICENSE) for details.
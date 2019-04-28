This is a reporter for the [go-metrics](https://github.com/rcrowley/go-metrics)
library which will post the metrics to [Librato](https://www.librato.com/). It
was originally part of the `go-metrics` library itself, but has been split off
to make maintenance of both the core library and the client easier.

### Usage

```go
import "github.com/mihasya/go-metrics-librato"
import "github.com/rcrowley/go-metrics"

go librato.Librato(metrics.DefaultRegistry,
    10e9,                  // interval
    "example@example.com", // account owner email address
    "token",               // Librato API token
    "hostname",            // source
    []float64{0.95},       // percentiles to send
    time.Millisecond,      // time unit
)

func Increment(name string) {
	m := metrics.GetOrRegisterMeter(name, nil)
	m.Mark(1)
}

func Time(name string, value time.Duration) {
	t := metrics.GetOrRegisterTimer(name, nil)
	t.Update(value)
}
```

### Migrating from `rcrowley/go-metrics` implementation

Simply modify the import from `"github.com/rcrowley/go-metrics/librato"` to
`"github.com/mihasya/go-metrics-librato"` and it should Just Work.

This is a reporter for the [go-metrics](https://github.com/rcrowley/go-metrics)
library which will post the metrics to [Librato](https://www.librato.com/). It
was originally part of the `go-metrics` library itself, but has been split off
to make maintenance of both the core library and the client easier.

## ...Fork specific changes and instructions...

### GOMAXPROCS Setting (you might not need to do this)
If the requests to post metrics are failing, set GOMAXPROCS to something > 1
See: https://github.com/golang/go/issues/4677
```
// Do this at the start of your program
runtime.GOMAXPROCS(2)
```

### Librato EOF

Librato servers can't understand the Counter + Gauge json data together for some reason.

Error on librato servers is on the lines of: http://stackoverflow.com/questions/6591388/configure-jackson-to-deserialize-single-quoted-invalid-json

This error is not actually returned by librato on the POST call that contains both Counter and Gauge data. However, the http call
just causes an EOF. To avoid this problem, this fork sends the Counter and Gauge data as two separate http calls until we can figure out a better solution.

## .../Fork specific changes and instructions...

### Usage

```go
import "github.com/mihasya/go-metrics-librato"

go librato.Librato(metrics.DefaultRegistry,
    10e9,                  // interval
    "example@example.com", // account owner email address
    "token",               // Librato API token
    "hostname",            // source
    []float64{0.95},       // percentiles to send
    time.Millisecond,      // time unit
)
```

### Migrating from `rcrowley/go-metrics` implementation

Simply modify the import from `"github.com/rcrowley/go-metrics/librato"` to
`"github.com/mihasya/go-metrics-librato"` and it should Just Work.

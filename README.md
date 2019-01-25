# About 
Opencensus enables you to easily collect metrics in your go application. This repository declares an exporter for your favorite metrics that directly pumps all those data points directly into InfluxDB, in order to be analyzed and graphed. Isn't that awesome?

# Usage
Adding this exporter to your project should be as simple as:

```golang
import (
    "github.com/influxdata/influxdb/client/v2"

    "go.opencensus.io/plugin/ochttp"
    "go.opencensus.io/stats/view"
    "go.opencensus.io/trace"
    "go.opencensus.io/zpages"

    "github.com/egymgmbh/opencensus-go-exporter-influxdb"
)

....

func main() {
    influxUrl  := "<url-of-the-influxdb>"
    influxUser := "<username>"
    influxPass := "<password>"
    influxDB   := "<db-name>"

    // Connect to the influxDB
    influxClient, err := client.NewHTTPClient(client.HTTPConfig{
        Addr:     influxUrl,
        Username: influxUser,
        Password: influxPass,
    })
    if err != nil {
        log.Fatalf("InfluxDB client connection failure: %v", err)
    }
    defer influxClient.Close()

    // Register our custom exporter to opencensus
    exporter := influxdb.NewExporter(influxClient, influxDB, func(err error) {
        log.Fatalf("error while registering Opencensus exporter: %v", err)
    }, nil)
    view.RegisterExporter(exporter)

    // Useful metrics to be exported from opencensus using the exporter
    view.Register(ochttp.ServerRequestCountByMethod)
    view.Register(ochttp.ServerResponseCountByStatusCode)
    view.Register(ochttp.ServerLatencyView)
    view.Register(ochttp.ServerRequestBytesView)
    view.Register(ochttp.ServerResponseBytesView)
}
```

## Custom Tags
The exporter can be configured to export additional tags with each metric in order to filter/refine your graphs. Just pass any wanted tags to the exporter as an additional parameter:

```golang
customTags := make(map[string]string)
customTags["hostname"] = getHostname()
view.RegisterExporter(influxdb.NewExporter(influxClient, influxDB, func(err error) {
    log.Fatal(err)
}, customTags))
```

# Contribute
Feel free to contribute and extend this exporter as you need it, we are always happy to accept pull requests.

# License
```
Copyright 2018-2018 eGym GmbH <support@egym.de>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
package influxdb

import (
	"fmt"

	"github.com/influxdata/influxdb/client/v2"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// NewExporter create a new InfluxDb exporter.
// It requires an influxDb client, a name for the database, one error handler and additionally accepts custom tags
// for the export.
func NewExporter(influxCli client.Client, database string, errorHandler func(error), customTags map[string]string) view.Exporter {
	return &exporter{influxCli: influxCli, database: database, errorHandler: errorHandler, customTags: customTags}
}

type exporter struct {
	influxCli    client.Client
	database     string
	errorHandler func(error)
	customTags   map[string]string
}

func (e *exporter) ExportView(viewData *view.Data) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  e.database,
		Precision: "s",
	})
	if err != nil {
		e.errorHandler(err)
		return
	}

	for _, row := range viewData.Rows {
		fields := make(map[string]interface{})

		switch d := row.Data.(type) {
		case *view.CountData:
			fields["value"] = float64(d.Value)
		case *view.DistributionData:
			fields["min"] = d.Min
			fields["max"] = d.Max
			fields["mean"] = d.Mean
			fields["count"] = d.Count
		case *view.LastValueData:
			fields["value"] = float64(d.Value)
		case *view.SumData:
			fields["value"] = float64(d.Value)
		default:
			e.errorHandler(fmt.Errorf("unknown AggregationData type: %T", row.Data))
			return
		}

		tagsMap := make(map[string]string)
		appendAndReplace(tagsMap, e.customTags)
		appendAndReplace(tagsMap, convertTags(row.Tags))

		pt, err := client.NewPoint(viewData.View.Name, tagsMap, fields, viewData.End)
		if err != nil {
			e.errorHandler(err)
		}
		bp.AddPoint(pt)
	}

	err = e.influxCli.Write(bp)
	if err != nil {
		e.errorHandler(err)
	}
}

// appendAndReplace appends all the data from the 'elementsMap' to the
// 'appendable' map. If both have the same key, the one from 'elementsMap'
// is taken.
func appendAndReplace(appendable, elementsMap map[string]string) {
	if appendable == nil {
		return
	}

	for k, v := range elementsMap {
		appendable[k] = v
	}
}

func convertTags(tags []tag.Tag) map[string]string {
	res := make(map[string]string)
	for _, tag := range tags {
		res[tag.Key.Name()] = tag.Value
	}
	return res
}

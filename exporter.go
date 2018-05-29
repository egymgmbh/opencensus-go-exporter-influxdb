package influxdb

import (
	"fmt"

	"github.com/influxdata/influxdb/client/v2"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

func NewExporter(influxCli client.Client, database string, errorHandler func(error)) view.Exporter {
	return &exporter{influxCli, database, errorHandler}
}

type exporter struct {
	influxCli client.Client
	database string
	errorHandler func(error)
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
		var suffix string

		switch d := row.Data.(type) {
		case *view.CountData:
			fields["value"] = float64(d.Value)
			suffix = ".count"
		case *view.DistributionData:
			fields["min"] = d.Min
			fields["max"] = d.Max
			fields["mean"] = d.Mean
			fields["count"] = d.Count
			suffix = ".histogram"
		case *view.LastValueData:
			fields["value"] = float64(d.Value)
			suffix = ".gauge"
		case *view.SumData:
			fields["value"] = float64(d.Value)
			suffix = ".gauge"
		default:
			e.errorHandler(fmt.Errorf("unknown AggregationData type: %T", row.Data))
			return
		}

		pt, err := client.NewPoint(viewData.View.Name + suffix, convertTags(row.Tags), fields, viewData.End)
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

func convertTags(tags []tag.Tag) map[string]string {
	res := make(map[string]string)
	for _, tag := range tags {
		res[tag.Key.Name()] = tag.Value
	}
	return res
}

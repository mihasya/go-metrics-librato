package librato

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const Operations = "operations"
const OperationsShort = "ops"

type LibratoClient struct {
	Email, Token string
}

// property strings
const (
	// display attributes
	Color             = "color"
	DisplayMax        = "display_max"
	DisplayMin        = "display_min"
	DisplayUnitsLong  = "display_units_long"
	DisplayUnitsShort = "display_units_short"
	DisplayStacked    = "display_stacked"
	DisplayTransform  = "display_transform"
	// special gauge display attributes
	SummarizeFunction = "summarize_function"
	Aggregate         = "aggregate"

	// metric keys
	Name        = "name"
	Period      = "period"
	Description = "description"
	DisplayName = "display_name"
	Attributes  = "attributes"

	// measurement keys
	MeasureTime = "measure_time"
	Source      = "source"
	Value       = "value"

	// special gauge keys
	Count      = "count"
	Sum        = "sum"
	Max        = "max"
	Min        = "min"
	SumSquares = "sum_squares"

	// batch keys
	Counters = "counters"
	Gauges   = "gauges"

	MetricsPostUrl = "https://metrics-api.librato.com/v1/metrics"
)

type Measurement map[string]interface{}
type Metric map[string]interface{}

type Batch struct {
	Gauges      []Measurement `json:"gauges,omitempty"`
	Counters    []Measurement `json:"counters,omitempty"`
	MeasureTime int64         `json:"measure_time"`
	Source      string        `json:"source"`
}

// Seems like librato doesn't like us sending gauges and counters together
// Split the structs to send two json batches.
type GaugeBatch struct {
	Gauges []Measurement `json:"gauges,omitempty"`
}

type CounterBatch struct {
	Counters []Measurement `json:"counters,omitempty"`
}

func (self *LibratoClient) PostMetrics(batch Batch) (err error) {

	var counterErr error
	if len(batch.Counters) > 0 {
		counterErr = self.PostCounters(batch)
	}

	var gaugeErr error
	if len(batch.Gauges) > 0 {
		gaugeErr = self.PostGauges(batch)
	}
	if counterErr != nil || gaugeErr != nil {
		return errors.New(fmt.Sprintf("Metrics post error. Gauge post error: %v Counter post error: %v", gaugeErr, counterErr))
	}
	return
}

func (self *LibratoClient) PostCounters(batch Batch) (err error) {
	var (
		js   []byte
		req  *http.Request
		resp *http.Response
	)

	for i := 0; i < len(batch.Counters); i++ {
		batch.Counters[i]["source"] = batch.Source
		batch.Counters[i]["measure_time"] = batch.MeasureTime
	}
	counterBatch := CounterBatch{batch.Counters}
	if js, err = json.Marshal(counterBatch); err != nil {
		return
	}

	if req, err = http.NewRequest("POST", MetricsPostUrl, bytes.NewBuffer(js)); err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(self.Email, self.Token)

	if resp, err = http.DefaultClient.Do(req); err != nil {
		log.Printf("Error Return. %s\n", err)
		return
	}

	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		body = []byte(fmt.Sprintf("(could not fetch response body for error: %s)", err))
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Unable to post to Librato: %d %s %s", resp.StatusCode, resp.Status, string(body))
	}
	log.Printf("Response from counter data post to Librato: %d %s %s\n", resp.StatusCode, resp.Status, string(body))
	resp.Body.Close()

	return
}

func (self *LibratoClient) PostGauges(batch Batch) (err error) {
	var (
		js   []byte
		req  *http.Request
		resp *http.Response
	)

	for i := 0; i < len(batch.Gauges); i++ {
		batch.Gauges[i]["source"] = batch.Source
		batch.Gauges[i]["measure_time"] = batch.MeasureTime
	}

	gaugeBatch := GaugeBatch{batch.Gauges}
	if js, err = json.Marshal(gaugeBatch); err != nil {
		return
	}
	if req, err = http.NewRequest("POST", MetricsPostUrl, bytes.NewBuffer(js)); err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(self.Email, self.Token)

	if resp, err = http.DefaultClient.Do(req); err != nil {
		log.Printf("Error Return. %s\n", err)
		return
	}

	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		body = []byte(fmt.Sprintf("(could not fetch response body for error: %s)", err))
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Unable to post to Librato: %d %s %s", resp.StatusCode, resp.Status, string(body))
	}
	log.Printf("Response from gauge data post to Librato: %d %s %s\n", resp.StatusCode, resp.Status, string(body))
	resp.Body.Close()

	return
}

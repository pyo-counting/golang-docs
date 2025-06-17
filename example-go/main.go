package main

import (
	"fmt"
	"time"

	"github.com/grafana/loki/v3/clients/pkg/logentry/stages"
	"github.com/grafana/loki/v3/clients/pkg/promtail/api"
	"github.com/grafana/loki/v3/pkg/logproto"
	"github.com/grafana/loki/v3/pkg/util/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
)

var testMultiStageYaml = `
pipeline_stages:
- json:
    expressions:
      log:
      stream:
      time:
- labels:
    stream:
`

var testLogEntry = `{"log":"11.11.11.11 - frank [25/Jan/2000:14:00:01 -0500] \"GET /1986.js HTTP/1.1\" 200 932 \"-\" \"Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 GTB6\"","stream":"stderr","time":"2019-04-30T02:12:41.8443515Z"}`

func loadConfig(yml string) stages.PipelineStages {
	var config map[string]interface{}
	err := yaml.Unmarshal([]byte(yml), &config)
	if err != nil {
		panic(err)
	}
	return config["pipeline_stages"].([]interface{})
}

type noopRegisterer struct{}

func (n noopRegisterer) Register(prometheus.Collector) error     { return nil }
func (n noopRegisterer) Unregister(prometheus.Collector) bool    { return true }
func (n noopRegisterer) MustRegister(cs ...prometheus.Collector) { return }

func main() {
	start := time.Now()
	p, err := stages.NewPipeline(log.Logger, loadConfig(testMultiStageYaml), nil, noopRegisterer{})
	if err != nil {
		panic(err)
	}

	e := stages.Entry{
		Extracted: make(map[string]interface{}),
		Entry: api.Entry{
			Labels: make(model.LabelSet),
			Entry: logproto.Entry{
				Timestamp: time.Now(),
				Line:      testLogEntry,
			},
		},
	}
	fmt.Printf("extracted: %v\n", e.Extracted)
	fmt.Printf("label: %v\n", e.Labels)
	fmt.Printf("line: %v\n", e.Entry.Line)
	c := make(chan stages.Entry, 1)
	c <- e
	o := <-p.Run(c)
	fmt.Printf("extracted: %v\n", o.Extracted)
	fmt.Printf("label: %v\n", o.Labels)
	fmt.Printf("line: %v\n", o.Entry.Line)

	end := time.Now()
	fmt.Println("elapsed time:", end.Sub(start), " / size:", p.Size(), " / name:", p.Name())
}

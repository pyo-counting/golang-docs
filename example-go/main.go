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
- match:
    selector: "{match=\"true\"}"
    stages:
    - docker:
    - regex:
        expression: "^(?P<ip>\\S+) (?P<identd>\\S+) (?P<user>\\S+) \\[(?P<timestamp>[\\w:/]+\\s[+\\-]\\d{4})\\] \"(?P<action>\\S+)\\s?(?P<path>\\S+)?\\s?(?P<protocol>\\S+)?\" (?P<status>\\d{3}|-) (?P<size>\\d+|-)\\s?\"?(?P<referer>[^\"]*)\"?\\s?\"?(?P<useragent>[^\"]*)?\"?$"
    - regex:
        source:     filename
        expression: "(?P<service>[^\\/]+)\\.log"
    - timestamp:
        source: timestamp
        format: "02/Jan/2006:15:04:05 -0700"
    - labels:
        action:
        service:
        status_code: "status"
- match:
    selector: "{match=\"false\"}"
    action: drop
`

func loadConfig(yml string) stages.PipelineStages {
	var config map[string]interface{}
	err := yaml.Unmarshal([]byte(yml), &config)
	if err != nil {
		panic(err)
	}
	return config["pipeline_stages"].([]interface{})
}

func newEntry(ex map[string]interface{}, lbs model.LabelSet, line string, ts time.Time) stages.Entry {
	if ex == nil {
		ex = map[string]interface{}{}
	}
	if lbs == nil {
		lbs = model.LabelSet{}
	}
	return stages.Entry{
		Extracted: ex,
		Entry: api.Entry{
			Labels: lbs,
			Entry: logproto.Entry{
				Timestamp: ts,
				Line:      line,
			},
		},
	}
}

func main() {
	start := time.Now()
	p, err := stages.NewPipeline(log.Logger, loadConfig(testMultiStageYaml), nil, prometheus.DefaultRegisterer)
	if err != nil {
		panic(err)
	}
	end := time.Now()
	fmt.Println("elapsed time:", end.Sub(start), " / size:", p.Size(), " / name:", p.Name())

}

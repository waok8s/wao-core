package v2inferenceprotocol

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/waok8s/wao-core/pkg/predictor"
	"github.com/waok8s/wao-core/pkg/util"
)

type ResponseTimePredictorTemplateData struct {
	App string
}

var tmpl *template.Template

func init() {
	tmpl = template.New("ResponseTimePredictor.Endpoint")
}

type ResponseTimePredictor struct {
	// urlTemplate contains the URL template.
	// E.g., "http://10.0.0.1:8090/v2/models/myModel:{{.App}}/versions/v0.1.0/infer"
	urlTemplate string

	client    *http.Client
	editorFns []util.RequestEditorFn
}

var _ predictor.ResponseTimePredictor = (*ResponseTimePredictor)(nil)

func NewResponseTimePredictor(urlTemplate string, insecureSkipVerify bool, timeout time.Duration, editorFns ...util.RequestEditorFn) *ResponseTimePredictor {
	return &ResponseTimePredictor{
		urlTemplate: urlTemplate,
		client: &http.Client{
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify}},
			Timeout:   timeout,
		},
		editorFns: editorFns,
	}
}

// Endpoint constructs the API endpoint. {{.App}} will be replaced with appName.
func (p *ResponseTimePredictor) Endpoint(appName string) (string, error) {
	data := ResponseTimePredictorTemplateData{
		App: appName,
	}

	t, err := tmpl.Parse(p.urlTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (p *ResponseTimePredictor) Predict(ctx context.Context, appName string, cpuUsage float64) (t float64, err error) {
	url, err := p.Endpoint(appName)
	if err != nil {
		return 0.0, fmt.Errorf("unable to get endpoint URL with appName %s: %w", appName, err)
	}

	body, err := json.Marshal(newInferResponseTimeRequest(appName, cpuUsage))
	if err != nil {
		return 0.0, fmt.Errorf("unable to marshal the request body=%+v err=%w", body, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return 0.0, fmt.Errorf("unable to create HTTP request: %w", err)
	}
	for i, f := range p.editorFns {
		if err := f(ctx, req); err != nil {
			return 0.0, fmt.Errorf("editorFns[%d] got error: %w", i, err)
		}
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return 0.0, fmt.Errorf("unable to send HTTP request: %w", err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		var apiResp inferResponseTimeResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return 0.0, fmt.Errorf("could not decode resp: %w", err)
		}
		if len(apiResp.Outputs) == 0 || len(apiResp.Outputs[0].Data) == 0 {
			return 0.0, fmt.Errorf("invalid response apiResp=%+v", apiResp)
		}
		return apiResp.Outputs[0].Data[0], nil
	default:
		return 0.0, fmt.Errorf("HTTP status=%s", resp.Status)
	}
}

// inferResponseTimeRequest holds a request.
//
// E.g.,
//
//	{
//	  "inputs": [
//	    {
//	      "name": "predict-prob",
//	      "shape": [ 1, 1 ],
//	      "datatype": "FP32",
//	      "data": [ [ 10 ] ]
//	    }
//	  ]
//	}
type inferResponseTimeRequest struct {
	Inputs []struct {
		Name     string      `json:"name"`
		Shape    []int       `json:"shape"`
		Datatype string      `json:"datatype"`
		Data     [][]float32 `json:"data"`
	} `json:"inputs"`
}

var (
	InferResponseTimeRequestInputName = "predict-prob"
)

func newInferResponseTimeRequest(_ string, cpuUsage float64) *inferResponseTimeRequest {
	var (
		name     = InferResponseTimeRequestInputName
		datatype = "FP32"
		shapeX   = 1
		shapeY   = 1
	)
	return &inferResponseTimeRequest{
		Inputs: []struct {
			Name     string      `json:"name"`
			Shape    []int       `json:"shape"`
			Datatype string      `json:"datatype"`
			Data     [][]float32 `json:"data"`
		}{
			{
				Name:     name,
				Shape:    []int{shapeX, shapeY},
				Datatype: datatype,
				Data:     [][]float32{{float32(cpuUsage)}},
			},
		},
	}
}

// inferResponseTimeResponse holds a response.
// Ignore values except outputs[*].data[]
//
// E.g.,
//
//	{
//	  "model_name": "model1",
//	  "model_version": "v0.1.0",
//	  "id": "0dc429d2-bd02-404b-b624-a0fa628e451e",
//	  "parameters": {
//	    "content_type": null,
//	    "headers": null
//	  },
//	  "outputs": [
//	    {
//	      "name": "predict",
//	      "shape": [ 1, 1 ],
//	      "datatype": "FP64",
//	      "parameters": null,
//	      "data": [ 20.0 ]
//	    }
//	  ]
//	}
type inferResponseTimeResponse struct {
	Outputs []struct {
		Data []float64 `json:"data"`
	} `json:"outputs"`
}

package client

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	waov1beta1 "github.com/waok8s/wao-core/api/wao/v1beta1"
	"github.com/waok8s/wao-core/pkg/predictor"
)

func Test_predictorCacheKey(t *testing.T) {
	type args struct {
		valueType     string
		namespace     string
		endpointTerm  *waov1beta1.EndpointTerm
		predictorType predictor.PredictorType
		appName       string
		cpuUsage      float64
		inletTemp     float64
		deltaP        float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "all fields",
			args: args{
				valueType: "foo",
				namespace: "bar",
				endpointTerm: &waov1beta1.EndpointTerm{
					Type:     "baz",
					Endpoint: "http://example.com",
					BasicAuthSecret: &corev1.LocalObjectReference{
						Name: "secret",
					},
					FetchInterval: &metav1.Duration{Duration: time.Second}, // not used
				},
				predictorType: predictor.PredictorType("test"),
				appName:       "hoge",
				cpuUsage:      1.234567890, // %.6f
				inletTemp:     2.000000000, // %.6f
				deltaP:        2.000,       // %.6f
			},
			want: "foo#bar#baz|http://example.com|secret#test#hoge#1.234568#2.000000#2.000000",
		},
		{
			name: "empty fields",
			args: args{
				valueType:     "",
				namespace:     "",
				endpointTerm:  nil,
				predictorType: "",
				appName:       "",
				cpuUsage:      0,
				inletTemp:     0,
				deltaP:        0,
			},
			want: "##||###0.000000#0.000000#0.000000",
		},
		{
			name: "valueTypeWatt",
			args: args{
				valueType: valueTypeWatt,
				namespace: "default",
				endpointTerm: &waov1beta1.EndpointTerm{
					Type:            waov1beta1.TypeV2InferenceProtocol,
					Endpoint:        "http://example.com",
					BasicAuthSecret: nil,
					FetchInterval:   nil,
				},
				predictorType: predictor.TypePowerConsumption,
				// appName:       "", // appName is not used
				cpuUsage:  1.234567890,
				inletTemp: 2.000000000,
				deltaP:    2.000,
			},
			want: "Watt#default#V2InferenceProtocol|http://example.com|#PowerConsumption##1.234568#2.000000#2.000000",
		},
		{
			name: "valueTypePowerConsumptionEndpoint",
			args: args{
				valueType: valueTypePowerConsumptionEndpoint,
				namespace: "default",
				endpointTerm: &waov1beta1.EndpointTerm{
					Type:            waov1beta1.TypeRedfish,
					Endpoint:        "http://example.com",
					BasicAuthSecret: &corev1.LocalObjectReference{Name: "secret"},
					FetchInterval:   nil,
				},
				predictorType: predictor.TypePowerConsumption,
				// appName, cpuUsage, inletTemp, deltaP are not used
			},
			want: "PowerConsumptionEndpoint#default#Redfish|http://example.com|secret#PowerConsumption##0.000000#0.000000#0.000000",
		},
		{
			name: "valueTypeResponseTime",
			args: args{
				valueType: valueTypeResponseTime,
				namespace: "default",
				endpointTerm: &waov1beta1.EndpointTerm{
					Type:            waov1beta1.TypeV2InferenceProtocol,
					Endpoint:        "http://example.com",
					BasicAuthSecret: nil,
					FetchInterval:   nil,
				},
				predictorType: predictor.TypeResponseTime,
				appName:       "hoge",
				cpuUsage:      1.234567890,
				// inletTemp, deltaP are not used
			},
			want: "ResponseTime#default#V2InferenceProtocol|http://example.com|#ResponseTime#hoge#1.234568#0.000000#0.000000",
		},
		{
			name: "valueTypeResponseTimeEndpoint",
			args: args{
				valueType: valueTypeResponseTimeEndpoint,
				namespace: "default",
				endpointTerm: &waov1beta1.EndpointTerm{
					Type:            waov1beta1.TypeRedfish,
					Endpoint:        "http://example.com",
					BasicAuthSecret: &corev1.LocalObjectReference{Name: "secret"},
					FetchInterval:   nil,
				},
				predictorType: predictor.TypeResponseTime,
				appName:       "hoge",
				// cpuUsage, inletTemp, deltaP are not used
			},
			want: "ResponseTimeEndpoint#default#Redfish|http://example.com|secret#ResponseTime#hoge#0.000000#0.000000#0.000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := predictorCacheKey(tt.args.valueType, tt.args.namespace, tt.args.endpointTerm, tt.args.predictorType, tt.args.appName, tt.args.cpuUsage, tt.args.inletTemp, tt.args.deltaP); got != tt.want {
				t.Errorf("predictorCacheKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

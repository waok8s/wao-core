package client

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"k8s.io/client-go/kubernetes"

	waov1beta1 "github.com/waok8s/wao-core/api/wao/v1beta1"
	"github.com/waok8s/wao-core/pkg/predictor"
	"github.com/waok8s/wao-core/pkg/predictor/fromnodeconfig"
)

type CachedPredictorClient struct {
	client kubernetes.Interface

	ttl   time.Duration
	cache sync.Map
}

func NewCachedPredictorClient(client kubernetes.Interface, ttl time.Duration) *CachedPredictorClient {
	return &CachedPredictorClient{
		client: client,
		ttl:    ttl,
	}
}

func predictorCacheKey(valueType string,
	namespace string, endpointTerm *waov1beta1.EndpointTerm, // common
	predictorType predictor.PredictorType, // GetPredictorEndpoint
	appName string, // PredictResponseTime
	cpuUsage float64, // PredictPowerConsumption, PredictResponseTime
	inletTemp float64, // PredictPowerConsumption
	deltaP float64, // PredictPowerConsumption
) string {
	if endpointTerm == nil {
		endpointTerm = &waov1beta1.EndpointTerm{}
	}
	secretName := ""
	if endpointTerm.BasicAuthSecret != nil {
		secretName = endpointTerm.BasicAuthSecret.Name
	}
	ep := fmt.Sprintf("%s|%s|%s", endpointTerm.Type, endpointTerm.Endpoint, secretName)
	return fmt.Sprintf("%s#%s#%s#%s#%s#%.6f#%.6f#%.6f", valueType, namespace, ep, predictorType, appName, cpuUsage, inletTemp, deltaP)
}

const (
	valueTypePowerConsumptionEndpoint = "PowerConsumptionEndpoint"
	valueTypeWatt                     = "Watt"

	valueTypeResponseTimeEndpoint = "ResponseTimeEndpoint"
	valueTypeResponseTime         = "ResponseTime"
)

type predictionCache struct {
	PowerConsumptionEndpoint *waov1beta1.EndpointTerm
	Watt                     float64

	ResponseTimeEndpoint *waov1beta1.EndpointTerm
	ResponseTime         float64

	ExpiredAt time.Time
}

func (c *CachedPredictorClient) do(ctx context.Context, valueType string,
	namespace string, endpointTerm *waov1beta1.EndpointTerm, // common
	predictorType predictor.PredictorType, // GetPredictorEndpoint
	appName string, // PredictResponseTime
	cpuUsage float64, // PredictPowerConsumption, PredictResponseTime
	inletTemp float64, // PredictPowerConsumption
	deltaP float64, // PredictPowerConsumption
) (*predictionCache, error) {

	key := predictorCacheKey(valueType, namespace, endpointTerm, predictorType, appName, cpuUsage, inletTemp, deltaP)
	lg := slog.With("func", "CachedPredictorClient.do", "key", key)

	if v, ok1 := c.cache.Load(key); ok1 {
		if cv, ok2 := v.(*predictionCache); ok2 {
			if cv.ExpiredAt.After(time.Now()) {
				lg.Debug("predictor cache hit")
				return cv, nil
			}
		}
	}
	lg.Debug("predictor cache missed")

	cv := &predictionCache{
		ExpiredAt: time.Now().Add(c.ttl),
	}

	switch valueType {
	case valueTypePowerConsumptionEndpoint:
		prov, err := fromnodeconfig.NewEndpointProvider(c.client, namespace, endpointTerm)
		if err != nil {
			return nil, err
		}
		ep, err := prov.Get(ctx, predictorType)
		if err != nil {
			return nil, err
		}
		cv.PowerConsumptionEndpoint = ep
	case valueTypeWatt:
		pred, err := fromnodeconfig.NewPowerConsumptionPredictor(c.client, namespace, endpointTerm)
		if err != nil {
			return nil, err
		}
		watt, err := pred.Predict(ctx, cpuUsage, inletTemp, deltaP)
		if err != nil {
			return nil, err
		}
		cv.Watt = watt
	case valueTypeResponseTimeEndpoint:
		prov, err := fromnodeconfig.NewEndpointProvider(c.client, namespace, endpointTerm)
		if err != nil {
			return nil, err
		}
		ep, err := prov.Get(ctx, predictorType)
		if err != nil {
			return nil, err
		}
		cv.ResponseTimeEndpoint = ep
	case valueTypeResponseTime:
		pred, err := fromnodeconfig.NewResponseTimePredictor(c.client, namespace, endpointTerm)
		if err != nil {
			return nil, err
		}
		t, err := pred.Predict(ctx, appName, cpuUsage)
		if err != nil {
			return nil, err
		}
		cv.ResponseTime = t
	default:
		return nil, fmt.Errorf("unknown valueType=%s", valueType)
	}

	c.cache.Store(key, cv)

	return cv, nil
}

func (c *CachedPredictorClient) GetPredictorEndpoint(ctx context.Context, namespace string, ep *waov1beta1.EndpointTerm, predictorType predictor.PredictorType) (*waov1beta1.EndpointTerm, error) {
	cv, err := c.do(ctx, valueTypePowerConsumptionEndpoint, namespace, ep, predictorType, "", 0.0, 0.0, 0.0)
	if err != nil {
		return nil, err
	}
	return cv.PowerConsumptionEndpoint, nil
}

func (c *CachedPredictorClient) PredictPowerConsumption(ctx context.Context, namespace string, ep *waov1beta1.EndpointTerm, cpuUsage, inletTemp, deltaP float64) (watt float64, err error) {
	cv, err := c.do(ctx, valueTypeWatt, namespace, ep, "", "", cpuUsage, inletTemp, deltaP)
	if err != nil {
		return 0.0, err
	}
	return cv.Watt, nil
}

func (c *CachedPredictorClient) PredictResponseTime(ctx context.Context, namespace string, ep *waov1beta1.EndpointTerm, appName string, cpuUsage float64) (t float64, err error) {
	cv, err := c.do(ctx, valueTypeResponseTime, namespace, ep, "", appName, cpuUsage, 0.0, 0.0)
	if err != nil {
		return 0.0, err
	}
	return cv.ResponseTime, nil
}

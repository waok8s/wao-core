package wao

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	waov1beta1 "github.com/waok8s/wao-core/api/wao/v1beta1"
	"github.com/waok8s/wao-core/pkg/metrics"
	"github.com/waok8s/wao-core/pkg/metrics/deltap"
	"github.com/waok8s/wao-core/pkg/metrics/inlettemp"
	"github.com/waok8s/wao-core/pkg/util"
)

var (
	Scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(Scheme))
	utilruntime.Must(waov1beta1.AddToScheme(Scheme))
}

// NodeConfigReconciler reconciles a NodeConfig object.
//
// NOTE: This reconciler is used in wao-metrics-adaptor instead of the controller. So RBAC rules below should be applied to wao-metrics-adaptor.
// kubebuilder:rbac:groups=wao.bitmedia.co.jp,resources=nodeconfigs,verbs=get;list;watch;create;update;patch;delete
// kubebuilder:rbac:groups=wao.bitmedia.co.jp,resources=nodeconfigs/status,verbs=get;update;patch
// kubebuilder:rbac:groups=wao.bitmedia.co.jp,resources=nodeconfigs/finalizers,verbs=update
// kubebuilder:rbac:groups=core,namespace=wao-system,resources=secrets,verbs=get
type NodeConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	MetricsCollector *metrics.Collector
	MetricsStore     *metrics.Store
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *NodeConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	lg := log.FromContext(ctx)
	lg.Info("Reconcile")

	var nc waov1beta1.NodeConfig
	err := r.Get(ctx, req.NamespacedName, &nc)
	if errors.IsNotFound(err) {
		r.reconcileNodeConfigDeletion(ctx, req.NamespacedName)
		return ctrl.Result{}, nil
	}
	if err != nil {
		lg.Error(err, "unable to get NodeConfig")
		return ctrl.Result{}, err
	}
	if !nc.DeletionTimestamp.IsZero() {
		r.reconcileNodeConfigDeletion(ctx, req.NamespacedName)
		return ctrl.Result{}, nil
	}

	if err := r.reconcileNodeConfig(ctx, req.NamespacedName, &nc); err != nil {
		lg.Error(err, "unable to reconcile NodeConfig", "obj", &nc)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *NodeConfigReconciler) reconcileNodeConfigDeletion(ctx context.Context, objKey types.NamespacedName) {
	lg := log.FromContext(ctx)
	lg.Info("reconcileNodeConfigDeletion")

	for _, vt := range metrics.ValueTypes {
		r.MetricsCollector.Unregister(metrics.CollectorKey(objKey, vt))
	}
}

type curlLogWriter struct {
	Logger logr.Logger
	Msg    string
}

func (w *curlLogWriter) Write(p []byte) (n int, err error) {
	w.Logger.Info(w.Msg, "curl", string(p))
	return len(p), nil
}

func (r *NodeConfigReconciler) getBasicAuthFromSecret(ctx context.Context, namespace string, ref *corev1.LocalObjectReference) (username, password string) {
	if ref == nil || ref.Name == "" {
		return
	}
	secret := &corev1.Secret{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: ref.Name}, secret); err != nil {
		// TODO: log
		return "", ""
	}

	username = string(secret.Data["username"])
	password = string(secret.Data["password"])

	return
}

func (r *NodeConfigReconciler) reconcileNodeConfig(ctx context.Context, objKey types.NamespacedName, nc *waov1beta1.NodeConfig) error {
	lg := log.FromContext(ctx)
	lg.Info("reconcileNodeConfig")

	inletTempConfig := nc.Spec.MetricsCollector.InletTemp
	switch inletTempConfig.Type {
	case waov1beta1.TypeFake:
		fetchTimeout := time.Second
		respondDelay := 100 * time.Millisecond
		c := &metrics.FakeClient{
			Type:  metrics.ValueInletTemperature,
			Value: 15.5,
			Error: nil,
			Delay: respondDelay,
		}
		r.MetricsCollector.Register(metrics.CollectorKey(objKey, metrics.ValueInletTemperature), c, r.MetricsStore, nc.Spec.NodeName, inletTempConfig.FetchInterval.Duration, fetchTimeout)
	case waov1beta1.TypeRedfish:
		serverType := inlettemp.TypeAutoDetect
		insecureSkipVerify := true
		fetchTimeout := 3 * time.Second
		requestTimeout := fetchTimeout - 1*time.Second
		username, password := r.getBasicAuthFromSecret(ctx, objKey.Namespace, inletTempConfig.BasicAuthSecret)

		requestEditorFns := []util.RequestEditorFn{
			util.WithBasicAuth(username, password),
			util.WithCurlLogger(&curlLogWriter{Logger: lg, Msg: "fetch inletTemp"}),
		}
		c := inlettemp.NewRedfishClient(inletTempConfig.Endpoint, serverType, insecureSkipVerify, requestTimeout, requestEditorFns...)
		r.MetricsCollector.Register(metrics.CollectorKey(objKey, metrics.ValueInletTemperature), c, r.MetricsStore, nc.Spec.NodeName, inletTempConfig.FetchInterval.Duration, fetchTimeout)
	default:
		return fmt.Errorf("unsupported metricsCollector.inletTemp.type: %s", inletTempConfig.Type)
	}

	deltapConfig := nc.Spec.MetricsCollector.DeltaP
	switch deltapConfig.Type {
	case waov1beta1.TypeFake:
		fetchTimeout := time.Second
		respondDelay := 100 * time.Millisecond
		c := &metrics.FakeClient{
			Type:  metrics.ValueDeltaPressure,
			Value: 7.5,
			Error: nil,
			Delay: respondDelay,
		}
		r.MetricsCollector.Register(metrics.CollectorKey(objKey, metrics.ValueDeltaPressure), c, r.MetricsStore, nc.Spec.NodeName, inletTempConfig.FetchInterval.Duration, fetchTimeout)
	case waov1beta1.TypeDPAPI:
		insecureSkipVerify := true
		fetchTimeout := 3 * time.Second
		requestTimeout := fetchTimeout - 1*time.Second
		username, password := r.getBasicAuthFromSecret(ctx, objKey.Namespace, inletTempConfig.BasicAuthSecret)

		requestEditorFns := []util.RequestEditorFn{
			util.WithBasicAuth(username, password),
			util.WithCurlLogger(&curlLogWriter{Logger: lg, Msg: "fetch deltaP"}),
		}
		c := deltap.NewDifferentialPressureAPIClient(deltapConfig.Endpoint, "", nc.Spec.NodeName, "", insecureSkipVerify, requestTimeout, requestEditorFns...)
		r.MetricsCollector.Register(metrics.CollectorKey(objKey, metrics.ValueDeltaPressure), c, r.MetricsStore, nc.Spec.NodeName, inletTempConfig.FetchInterval.Duration, fetchTimeout)
	default:
		return fmt.Errorf("unsupported metricsCollector.deltaP.type: %s", deltapConfig.Type)
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&waov1beta1.NodeConfig{}).
		Complete(r)
}

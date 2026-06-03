package __FEATURE_NAME__

import "github.com/prometheus/client_golang/prometheus"

type metricScope int

const (
	metricScopeFeature metricScope = iota
	metricScopeNamespace
	metricScopeAbsolute
)

type metricSpec struct {
	ID     string
	Scope  metricScope
	Name   string
	Help   string
	Labels []string
}

type metricDescriptors struct {
	order []string
	descs map[string]*prometheus.Desc
}

const metricExampleValue = "example_value"

var featureMetricSpecs = []metricSpec{
	{
		ID:    metricExampleValue,
		Scope: metricScopeFeature,
		Name:  "_example_value",
		Help:  "Example __FEATURE_NAME__ metric emitted by the generated exporter skeleton",
	},
}

func loadMetricDescriptors(featureName string, namespace string, specs []metricSpec) metricDescriptors {
	metrics := metricDescriptors{
		order: make([]string, 0, len(specs)),
		descs: make(map[string]*prometheus.Desc, len(specs)),
	}
	for _, spec := range specs {
		metrics.order = append(metrics.order, spec.ID)
		metrics.descs[spec.ID] = prometheus.NewDesc(
			spec.metricName(featureName, namespace),
			spec.Help,
			spec.Labels,
			nil,
		)
	}
	return metrics
}

func (d metricDescriptors) Describe(ch chan<- *prometheus.Desc) {
	for _, id := range d.order {
		ch <- d.Get(id)
	}
}

func (d metricDescriptors) Get(id string) *prometheus.Desc {
	return d.descs[id]
}

func metricName(featureName string, namespace string, id string) string {
	for _, spec := range featureMetricSpecs {
		if spec.ID == id {
			return spec.metricName(featureName, namespace)
		}
	}
	return id
}

func (s metricSpec) metricName(featureName string, namespace string) string {
	switch s.Scope {
	case metricScopeFeature:
		return featureName + s.Name
	case metricScopeNamespace:
		return namespace + s.Name
	case metricScopeAbsolute:
		return s.Name
	default:
		return s.Name
	}
}

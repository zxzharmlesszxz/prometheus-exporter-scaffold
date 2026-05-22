package exporter

import framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"

func (f *Feature) RuntimeConfig() []any {
	return []any{
		"refresh_interval", framework.NormalizeDuration(f.refreshInterval, defaultRefreshInterval),
	}
}

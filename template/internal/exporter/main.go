package exporter

import framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"

func Main() {
	framework.MainFromProject(NewFeature())
}

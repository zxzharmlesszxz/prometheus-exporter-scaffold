package main

import (
	template "github.com/zxzharmlesszxz/prometheus-template-exporter/exporter"

	"__GO_MODULE__/internal/exporter"
)

const projectName = "__PROJECT_NAME__"
const projectDesc = "__PROJECT_DESC__"

func main() {
	template.MainForProject(
		projectName,
		projectDesc,
		exporter.NewFeature(),
	)
}

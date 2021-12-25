package diag

import (
	"github.com/xjlgod/nebula-database-diagnostic/intrenal/info"
	"strings"
)

func GetDiagResult(allInfo *info.AllInfo) []string {
	phyDiagResult := GetPhyDiag(allInfo.PhyInfo)
	metricsDiagResult := GetMetricsDiag(allInfo.MetricsInfo)

	return []string{strings.Join(phyDiagResult, ""), strings.Join(metricsDiagResult, "")}
}

package diag

import (
	"fmt"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/service"
	"strconv"
)

func GetMetricsDiag(infos []*service.ServiceMetricInfo) (diags []string) {
	diags = append(diags, "metrics nothing.\n ")
	diags = metricsDiag(infos, diags)
	if len(diags) > 1 {
		return diags[1:]
	}
	return diags
}

func metricsDiag(infos []*service.ServiceMetricInfo, diags []string) []string {
	diags = append(diags, "==> diag metrics info:\n ")
	for _, serviceMetricsInfo := range infos {
		componentType := serviceMetricsInfo.Type
		metrics := serviceMetricsInfo.Metrics
		switch componentType {
		case config.GraphService:
			if metrics["num_queries.sum.600"] == service.NotCollect {
				diag := fmt.Sprintf("%v graph servie: num_queries.sum.600 is not collected", serviceMetricsInfo.Name)
				diags = append(diags, diag)
				break
			}
			num, err := strconv.Atoi(metrics["num_queries.sum.600"])
			if err != nil {
				diag := fmt.Sprintf("%v graph servie: num_queries.sum.600 is a wrong number", serviceMetricsInfo.Name)
				diags = append(diags, diag)
				break
			}
			if num > ThresholdIdleNumQuerisSum600 {
				diag := fmt.Sprintf("%v graph servie: num_queries.sum.600 is a wrong number", serviceMetricsInfo.Name)
				diags = append(diags, diag)
				break
			}
			diag := fmt.Sprintf("%v graph servie: num_queries.sum.600 is right", serviceMetricsInfo.Name)
			diags = append(diags, diag)
		case config.MetaService:
			if metrics["heartbeat_latency_us.avg.600"] == service.NotCollect {
				diag := fmt.Sprintf("%v meta servie: heartbeat_latency_us.avg.600 is not collected", serviceMetricsInfo.Name)
				diags = append(diags, diag)
				break
			}
			num, err := strconv.Atoi(metrics["num_queries.sum.600"])
			if err != nil {
				diag := fmt.Sprintf("%v meta servie: heartbeat_latency_us.avg.600 is a wrong number", serviceMetricsInfo.Name)
				diags = append(diags, diag)
				break
			}
			if num < ThresholdIdleHeartbeatLatencyUsAvg600 {
				diag := fmt.Sprintf("%v meta servie: heartbeat_latency_us.avg.600 is a wrong number", serviceMetricsInfo.Name)
				diags = append(diags, diag)
				break
			}
			diag := fmt.Sprintf("%v meta servie: heartbeat_latency_us.avg.600 is right", serviceMetricsInfo.Name)
			diags = append(diags, diag)
		case config.StorageService:
			if metrics["num_lookup_errors.sum.600"] == service.NotCollect {
				diag := fmt.Sprintf("%v storage servie: num_lookup_errors.sum.600 is not collected", serviceMetricsInfo.Name)
				diags = append(diags, diag)
				break
			}
			num, err := strconv.Atoi(metrics["num_lookup_errors.sum.600"])
			if err != nil {
				diag := fmt.Sprintf("%v storage servie: num_lookup_errors.sum.600 is a wrong number", serviceMetricsInfo.Name)
				diags = append(diags, diag)
				break
			}
			if num > ThresholdIdleNumLookupErrorsSum600 {
				diag := fmt.Sprintf("%v storage servie: num_lookup_errors.sum.600 is a wrong number", serviceMetricsInfo.Name)
				diags = append(diags, diag)
				break
			}
			diag := fmt.Sprintf("%v storage servie: num_lookup_errors.sum.600 is right", serviceMetricsInfo.Name)
			diags = append(diags, diag)
		}
	}
	return diags
}
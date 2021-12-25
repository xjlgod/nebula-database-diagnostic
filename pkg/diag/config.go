package diag

import (
	"fmt"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/info/service"
)

func GetConfigDiag(infos []*service.ServiceConfigInfo) (diags []string) {
	
	diags = append(diags, "config nothing.\n")
	diags = checkConfigDiag(infos, diags)
	if len(diags) > 1 {
		return diags[1:]
	}
	return diags

}

func checkConfigDiag(infos []*service.ServiceConfigInfo, diags []string) []string {
	diags = append(diags, "==> diag config info:\n")
	for _, serviceConfigInfo := range infos {
		configs := serviceConfigInfo.Configs
		flag := true
		for _, value := range configs {
			if value == service.NotCollect {
				flag = false
				diag := fmt.Sprintf("check %v config: config is incomplete\n", serviceConfigInfo.Name)
				diags = append(diags, diag)
				break
			}
		}
		if flag {
			diag := fmt.Sprintf("check %v config: config is complete\n", serviceConfigInfo.Name)
			diags = append(diags, diag)
		}
	}
	return diags
}

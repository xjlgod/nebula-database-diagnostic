package diag

import "github.com/xjlgod/nebula-database-diagnostic/pkg/info/service"

func GetConfigDiag([]*service.ServiceConfigInfo) (diags []string) {
	
	diags = append(diags, "physical nothing.\n")

	if len(diags) > 1 {
		return diags[1:]
	}
	return diags

}

func processConfigDiag()  {
	
}

package diag

import "github.com/xjlgod/nebula-database-diagnostic/pkg/info/physical"

func GetPhyDiag(info *physical.PhyInfo) (diags []string) {
	diags = append(diags, "physical nothing.\n ")
	diags = processDiag(info, diags)
	if len(diags) > 1 {
		return diags[1:]
	}
	return diags
}

func processDiag(info *physical.PhyInfo, diags []string) []string {
	diags = append(diags, "==> diag process info:\n ")
	if info.Process.WaitNumber >= ThresholdWaitNumber {
		diags = append(diags, "==> process wait number: process wait number is above the threshold.\n")
		if float64(info.Process.RunNumber/info.CPU.LogicNumber) > ThresholdRunNumber {
			diags = append(diags, "==> process run number: cpu is too busy.\n")
		} else {
			diags = cpuDiag(info, diags)
			diags = memoryDiag(info, diags)
		}
	}
	return diags
}

func cpuDiag(info *physical.PhyInfo, diags []string) []string {
	diags = append(diags, "==> diag cpu info:\n ")
	if float64(info.CPU.IdleTime/100) > ThresholdIdleTime {
		diags = append(diags, "==> diag cpu idle time: idle time percent is too big\n")
	}

	if float64(info.CPU.WaitPercent/100) > ThresholdWaitPercent {
		diags = append(diags, "==> diag cpu wait percent: wait percent is too big. there have io problems\n")
	}

	return diags
}

func memoryDiag(info *physical.PhyInfo, diags []string) []string {
	diags = append(diags, "==> diag memory info:\n ")
	if float64(info.Memory.MemFree/info.Memory.MemTotal) >= ThresholdMemoryFree {
		diags = append(diags, "==> diag memory free: memory free is too small\n")
	}
	return diags
}

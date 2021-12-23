package remote

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetNebulaMetrics(ipAddress string, port int) ([]string, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d/stats?stats=", ipAddress, port))
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	metricStr := string(bytes)
	metrics := strings.Split(metricStr, "\n")
	if len(metrics) > 0 {
		if metrics[len(metrics)-1] == "" {
			metrics = metrics[:len(metrics)-1]
		}
	}

	return metrics, nil
}

func GetNebulaConfigs(ipAddress string, port int) ([]string, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d/flags", ipAddress, port))
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	configStr := string(bytes)
	configs := strings.Split(configStr, "\n")
	if len(configs) > 0 {
		if configs[len(configs)-1] == "" {
			configs = configs[:len(configs)-1]
		}
	}

	return configs, nil
}

func GetNebulaComponentStatus(ipAddress string, port int) ([]string, error) {
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d/status", ipAddress, port))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type nebulaStatus struct {
		GitInfoSha string `json:"git_info_sha"`
		Status     string `json:"status"`
	}

	var status nebulaStatus
	if err := json.Unmarshal(bytes, &status); err != nil {
		return nil, err
	}

	statusMetrics := []string{status.Status}

	return statusMetrics, nil
}

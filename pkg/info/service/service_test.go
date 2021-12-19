package service

import (
	"fmt"
	"testing"
)

func TestMetaExporter(t *testing.T) {
	metaExporter := new(MetaExporter)
	var serviceExporter ServiceExporter = metaExporter
	serviceExporter.Config("192.168.8.169", 19559)
	serviceExporter.BuildMetricMap()
	serviceExporter.Collect()
	withMetricMap := serviceExporter.GetWithMetricMap()
	//循环遍历Map
	for key,value:= range withMetricMap{
		fmt.Printf("%s=>%s\n",key,value)
	}
}

func TestGraphExporter(t *testing.T) {
	metaExporter := new(GraphExporter)
	var serviceExporter ServiceExporter = metaExporter
	serviceExporter.Config("192.168.8.169", 19669)
	serviceExporter.BuildMetricMap()
	serviceExporter.Collect()
	withMetricMap := serviceExporter.GetWithMetricMap()
	//循环遍历Map
	for key,value:= range withMetricMap{
		fmt.Printf("%s=>%s\n",key,value)
	}
}

func TestStorageExporter(t *testing.T) {
	metaExporter := new(StorageExporter)
	var serviceExporter ServiceExporter = metaExporter
	serviceExporter.Config("192.168.8.169", 19779)
	serviceExporter.BuildMetricMap()
	serviceExporter.Collect()
	withMetricMap := serviceExporter.GetWithMetricMap()
	//循环遍历Map
	for key,value:= range withMetricMap{
		fmt.Printf("%s=>%s\n",key,value)
	}
}
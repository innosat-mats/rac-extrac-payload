package common

import (
	"sync"
)

// ExtractData reads Ramses data packages and extract the instrument data.
func ExtractData(callback func(data DataRecord), streamBatch ...StreamBatch) {
	var waitGroup sync.WaitGroup
	records := make(chan DataRecord)
	innosatPackages := make(chan DataRecord)
	dataPackages := make(chan DataRecord)
	go DataPackets(innosatPackages, dataPackages)
	go Packets(records, streamBatch...)

	go func() {
		waitGroup.Add(1)
		for data := range dataPackages {
			callback(data)
		}
		waitGroup.Done()
	}()

	for ramsesPackage := range records {
		innosatPackage, err := DecodeSource(ramsesPackage.Buffer)
		if err != nil {
			ramsesPackage.SourceHeader = innosatPackage.Header
			ramsesPackage.TMHeader = innosatPackage.Payload
			ramsesPackage.Buffer = innosatPackage.ApplicationPayload
			ramsesPackage.Error = err
			innosatPackages <- ramsesPackage
			continue
		}
		ramsesPackage.SourceHeader = innosatPackage.Header
		ramsesPackage.TMHeader = innosatPackage.Payload
		ramsesPackage.Buffer = innosatPackage.ApplicationPayload
		innosatPackages <- ramsesPackage
	}
	close(innosatPackages)
	waitGroup.Wait()
}

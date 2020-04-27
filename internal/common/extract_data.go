package common

import (
	"sync"
)

// ExtractCallback is the type of the callback function
type ExtractCallback func(data DataRecord)

// ExtractFunction is the type of the ExtractData function
type ExtractFunction func(callback ExtractCallback, streamBatch ...StreamBatch)

// ExtractData reads Ramses data packages and extract the instrument data.
func ExtractData(callback ExtractCallback, streamBatch ...StreamBatch) {
	var waitGroup sync.WaitGroup
	ramsesChannel := make(chan DataRecord)
	innosatChannel := make(chan DataRecord)
	aezChannel := make(chan DataRecord)

	go DecodeRamses(ramsesChannel, streamBatch...)
	go DecodeAEZ(aezChannel, innosatChannel)

	go func() {
		waitGroup.Add(1)
		for data := range aezChannel {
			callback(data)
		}
		waitGroup.Done()
	}()

	for record := range ramsesChannel {
		innosatPackage, err := DecodeSource(record.Buffer)
		if err != nil {
			record.Error = err
		}
		record.SourceHeader = innosatPackage.Header
		record.TMHeader = innosatPackage.Payload
		record.Buffer = innosatPackage.ApplicationPayload
		innosatChannel <- record
	}
	close(innosatChannel)
	waitGroup.Wait()
}

package extractors

import (
	"sync"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

// ExtractFunction is the type of the ExtractData function
type ExtractFunction func(callback common.Callback, streamBatch ...StreamBatch)

// ExtractData reads Ramses data packages and extract the instrument data.
func ExtractData(callback common.Callback, streamBatch ...StreamBatch) {
	var waitGroup sync.WaitGroup
	ramsesChannel := make(chan common.DataRecord)
	innosatChannel := make(chan common.DataRecord)
	aezChannel := make(chan common.DataRecord)

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

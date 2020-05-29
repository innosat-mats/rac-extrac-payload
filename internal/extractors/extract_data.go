package extractors

import (
	"sync"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

// ExtractFunction is the type of the ExtractData function
type ExtractFunction func(callback common.Callback, streamBatch ...StreamBatch)

const channelBufferSize int = 1024

// ExtractData reads Ramses data packages and extract the instrument data.
func ExtractData(callback common.Callback, streamBatch ...StreamBatch) {
	var waitGroup sync.WaitGroup
	ramsesChannel := make(chan common.DataRecord, channelBufferSize)
	innosatChannel := make(chan common.DataRecord, channelBufferSize)
	aggregatorChannel := make(chan common.DataRecord, channelBufferSize)
	aezChannel := make(chan common.DataRecord, channelBufferSize)

	go DecodeRamses(ramsesChannel, streamBatch...)
	go Aggregator(aggregatorChannel, innosatChannel)
	go DecodeAEZ(aezChannel, aggregatorChannel)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		for data := range aezChannel {
			callback(data)
		}
	}()

	for record := range ramsesChannel {
		if record.Error != nil {
			innosatChannel <- record
			continue
		}
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

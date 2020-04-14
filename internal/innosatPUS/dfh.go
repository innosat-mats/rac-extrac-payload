package pus

//DataHeader is
type DataHeader interface {
	getPUS() uint8
	getServiceType() uint8
	getServiceSubtype() uint8
	getTime() uint32
}

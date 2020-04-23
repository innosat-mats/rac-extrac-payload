package innosat

// SourcePacketAPIDType for service type
type SourcePacketAPIDType uint8

const (
	// TimeAPID Time Reporting (reserved)
	TimeAPID SourcePacketAPIDType = 0
	// ScmAPID Spacecraft Control Module (SCM) Handler (HW and BSW)
	ScmAPID SourcePacketAPIDType = 10
	// RwAPID Reaction Wheel Equipment Handler
	RwAPID SourcePacketAPIDType = 11
	// StrAPID Star Tracker Equipment Handler
	StrAPID SourcePacketAPIDType = 12
	// MagAPID Magnetometer Equipment Handler
	MagAPID SourcePacketAPIDType = 13
	// MtqAPID Magnetorquer Equipment Handler
	MtqAPID SourcePacketAPIDType = 14
	// GpsAPID GPS Equipment Handler
	GpsAPID SourcePacketAPIDType = 15
	// MpduAPID MPDU Equipment Handler
	MpduAPID SourcePacketAPIDType = 16
	// TcmAPID TCM Equipment Handler
	TcmAPID SourcePacketAPIDType = 17
	// DpcuAPID DPCU Equipment Handler
	DpcuAPID SourcePacketAPIDType = 18
	// SttqAPID Spacecraft Telecommand Scheduler
	SttqAPID SourcePacketAPIDType = 30
	// SysAPID System Core
	SysAPID SourcePacketAPIDType = 31
	// PowAPID Power Core
	PowAPID SourcePacketAPIDType = 32
	// TcsAPID Thermal Core
	TcsAPID SourcePacketAPIDType = 33
	// AcsAPID ACS Core
	AcsAPID SourcePacketAPIDType = 34
	// SfdirAPID Spacecraft FDIR Core
	SfdirAPID SourcePacketAPIDType = 35
	// McmAPID Mission Control Module (MCM) Handler (HW and BSW)
	McmAPID SourcePacketAPIDType = 50
	// MttqAPID Mission Telecommand Scheduler
	MttqAPID SourcePacketAPIDType = 70
	// OrbAPID Orbit Estimation
	OrbAPID SourcePacketAPIDType = 71
	// PlmAPID Payload Management
	PlmAPID SourcePacketAPIDType = 72
	// MainAPID Payload-specific APIDs (TBD)
	MainAPID SourcePacketAPIDType = 100
	// IdleAPID Idle Packet?
	IdleAPID SourcePacketAPIDType = 255
)

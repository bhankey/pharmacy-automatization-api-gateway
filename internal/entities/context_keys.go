package entities

type ContextKey int

const (
	UserID ContextKey = iota
	Email
	DeviceFingerPrint
	RequestID
	PharmacyID
)

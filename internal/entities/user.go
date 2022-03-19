package entities

type User struct {
	ID                int
	Name              string
	Surname           string
	Email             string
	Password          string
	PasswordHash      string
	Role              Role
	UseIPCheck        bool
	IsBlock           bool
	DefaultPharmacyID int
}

// TODO use this for all auth operations.
type UserIdentifyData struct {
	IP          string
	UserAgent   string
	FingerPrint string
}

type Role string

const (
	// All use only for access. Don't save in DB this value.
	All        Role = "all"
	Admin      Role = "admin"
	Apothecary Role = "apothecary"
)

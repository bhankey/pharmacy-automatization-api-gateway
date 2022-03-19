package entities

type Pharmacy struct {
	ID        int
	Name      string
	IsBlocked bool
	Address   Address
}

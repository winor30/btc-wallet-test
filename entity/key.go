package entity

// Key is private or public key for HD wallet
type Key struct {
	Net     uint8
	Type    uint8
	Indexes []uint32
	Value   string
}

var (
	NullKey = Key{}
)

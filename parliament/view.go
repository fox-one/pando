package parliament

type Item struct {
	Key    string
	Value  string
	Action string
}

type Proposal struct {
	Number int64
	Action string
	Info   []Item
	Meta   []Item

	ApprovedCount int
	ApprovedBy    string
}

type Flip struct {
	Number int64
	Info   []Item
}

type FlipStat struct {
	Lot string
	Bid string
	Tab string
	Gem string
	Dai string
}

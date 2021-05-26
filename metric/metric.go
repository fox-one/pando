package metric

type Entry struct {
	Name  string
	Value string
}

type Group struct {
	Name    string
	Entries []Entry
}

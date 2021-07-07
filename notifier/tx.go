package notifier

type TxData struct {
	Action  string
	Message string
	Lines   []string

	// copy to another
	receipts []string
}

func (data *TxData) AddLine(line string) {
	data.Lines = append(data.Lines, line)
}

func (data *TxData) cc(userID string) {
	data.receipts = append(data.receipts, userID)
}

package core

import (
	"testing"
)

func TestEncodeTransactionAction(t *testing.T) {
	var action TransactionAction
	b, err := action.Encode()
	t.Log(len(b), err)
}

package parliament

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func Test_execute(t *testing.T) {
	type args struct {
		name string
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "nil data",
			args: args{
				name: "proposal_passed",
				data: nil,
			},
			want: loadTestData("proposal_passed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := execute(tt.args.name, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func loadTestData(name string) []byte {
	name = fmt.Sprintf("testdata/%s.txt", name)
	b, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}

	return bytes.TrimSpace(b)
}

package db

import (
	"log"
	"testing"
)

func init() {
	Init()
	log.Printf("test setup.")
}

func TestInsert(t *testing.T) {
	type args struct {
		keyVal *Pair
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"test1",
			args{
				&Pair{1, []byte{'f', 'o', '0'}},
			},
			true,
			false, // do not want error
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Insert(tt.args.keyVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

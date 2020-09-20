package db

import (
	"log"
	"testing"
)

func init() {
	Init()
	log.Printf("test setup.")
}

func beforeTestInsert() {
	Init()
}

func afterTestInsert() {

}

func TestInsert(t *testing.T) {
	beforeTestInsert()
	defer afterTestInsert()

	type args struct {
		keyVal *Pair
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"test1", args{&Pair{2, []byte{'f', 'o', 'o'}}}, true, false}, // do not want error
		{"test2", args{&Pair{1, []byte{'b', 'a', 'r'}}}, true, false}, // do not want error
		{"test3", args{&Pair{5, []byte{'b', 'a', 'z'}}}, true, false}, // do not want error
		//{"test4", args{&Pair{3, []byte{'z', 'o', 'e'}}}, true, false}, // do not want error
		//{"test5", args{&Pair{4, []byte{'d', 'o', 't'}}}, true, false}, // do not want error
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
	// dump root node
	log.Printf("root: %+v", root)
}

func SetUpInsertIntoNode() {
	Init()
}

func TearDownInsertIntoNode() {
	//TODO
	//deleteDB
}

func Test_insertIntoNode(t *testing.T) {
	// setup and teardown
	SetUpInsertIntoNode()
	defer TearDownInsertIntoNode()

	type args struct {
		cursor Cursor
		kv     *Pair
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"test_insert_into_empty_node",
			args{
				Cursor{NewNode(), 0},
				&Pair{1, []byte{'a', 'b', 'c'}},
			},
			true, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := insertIntoNode(tt.args.cursor, tt.args.kv)
			if (err != nil) != tt.wantErr {
				t.Errorf("insertIntoNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("insertIntoNode() = %v, want %v", got, tt.want)
			}
		})
		// dump root node
		log.Printf("node: %+v", *tt.args.cursor.Node)
	}
}

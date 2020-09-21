package main

import (
	"fmt"
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
		{"test4", args{&Pair{3, []byte{'z', 'o', 'e'}}}, true, false}, // do not want error
		{"test5", args{&Pair{4, []byte{'d', 'o', 't'}}}, true, false}, // do not want error
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
func TestInsertAndSearch1(t *testing.T) {
	Init()

	for i := 1; i <= 10; i++ {
		str := fmt.Sprintf("foo_%v", i)
		Insert(&Pair{i, []byte(str)})
	}
	p := Search(4)
	//log.Printf("found p: %v, root: %v", string(p.Val), root)
	log.Printf("found p: %v, root: %v", p, root)
}

func TestInsertAndSearch0(t *testing.T) {
	Init()
	kv1 := &Pair{10, []byte("foo")}
	kv2 := &Pair{20, []byte("bar")}
	kv3 := &Pair{30, []byte("baz")}

	var (
		ok  bool
		err error
		p   *Pair
	)
	ok, err = Insert(kv2)
	log.Printf("ok:%v, err:%v", ok, err)
	ok, err = Insert(kv1)
	log.Printf("ok:%v, err:%v", ok, err)
	ok, err = Insert(kv3)
	log.Printf("ok:%v, err:%v", ok, err)

	p = Search(30)
	log.Printf("found p: %v, root: %v", string(p.Val), root)
	p = Search(10)
	log.Printf("found p: %v, root: %v", string(p.Val), root)
	p = Search(20)
	log.Printf("found p: %v, root: %v", string(p.Val), root)

}

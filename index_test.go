package boltcompact

import (
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

func TestCompact(t *testing.T) {
	srcPath := "./data/src.db"
	dstPath := "./data/dst.db"
	err := Compact("", dstPath, nil)
	if err == nil {
		t.Fatal("Src path need check empty")
	}

	srcDB, err := bolt.Open(srcPath, os.ModePerm, nil)
	if err != nil {
		t.Fatal("db open error")
	}
	srcDB.Close()

	err = Compact(srcPath, dstPath, nil)
	if err != nil {
		t.Fatal("Src path need check failed")
	}
}

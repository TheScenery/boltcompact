package boltcompact

import (
	"io/fs"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestCompact(t *testing.T) {
	srcPath := "./data/src.db"
	dstPath := "./data/dst.db"
	err := Compact("", dstPath, nil)
	if err == nil {
		t.Fatal("Src path need check empty")
	}

	srcDB, err := bolt.Open(srcPath, fs.ModePerm, nil)
	if err != nil {
		t.Fatal("db open error")
	}
	srcDB.Close()

	err = Compact(srcPath, dstPath, nil)
	if err != nil {
		t.Fatal("Src path need check failed")
	}
}

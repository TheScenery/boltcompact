# BoltCompact
compact [bolt db](https://github.com/boltdb/bolt) to save storage usage

bolt client command "compact" to lib

## Installing
```shell
go get github.com/TheScenery/boltcompact
```
## How to use
```go
import (
	"github.com/TheScenery/boltcompact"
)

func main() {
	boltcompact.Compact("src.db", "dst.db", nil)
}
```

## A message from the author
If you are using the successor of bold db named [bbolt](https://github.com/etcd-io/bbolt). It's a more featureful version of Bolt, And has contains this feature in it's API. You will not need to use this lib. Just use like below: 
```go
import bolt "go.etcd.io/bbolt"

// Open source database.
src, err := bolt.Open("src.db", 0444, nil)
if err != nil {
	return err
}
defer src.Close()

// Open destination database.
dst, err := bolt.Open("dst.db", fi.Mode(), nil)
if err != nil {
	return err
}
defer dst.Close()

// Run compaction.
txMaxSize := 65536
if err := bolt.Compact(dst, src, txMaxSize); err != nil {
	return err
}
```
Good Luck!
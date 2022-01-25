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

 bolt.Compact("src.db", "dst.db", 65536)
```
Good Luck!
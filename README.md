# BoltCompact
compact bolt db to save storage usage

bolt client command "compact" to lib

## Installing
```shell
go get github.com/TheScenery/BoltCompact
```
## How to use
```go
import (
	boltcompact "github.com/TheScenery/boltcompact"
)

func main() {
	boltcompact.Compact("src.db", "dst.db", nil)
}
```


# BoltCompact
compact bolt db to save storage usage

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


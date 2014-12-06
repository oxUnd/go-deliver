## go-deliver

### install

```bash
$ go get github.com/xiangshouding/go-deliver
```

### use

```go
package main

import (
    "fmt"
    d "github.com/xiangshouding/go-deliver"
    "runtime"
    "path"
)

func main() {
    fmt.Println("Map start");
    _,f,_,_ := runtime.Caller(0)
    dir := path.Dir(f)
    dd := d.NewDeliver(path.Join(dir, "./from"), path.Join(dir, "./to"));
    dd.Release(map[string]string{
        "reg": ".*\\.log",
        "release": "/static/$&",
    }, map[string]string{
        "reg": ".*\\.js",
        "release": "/js/$&",
    })
}
```

### api

@TODO

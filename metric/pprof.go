package metric

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func StartPprof() {
	// http://localhost:6060/debug/pprof/allocs
	// http://localhost:6060/debug/pprof/block
	// http://localhost:6060/debug/pprof/cmdline
	// http://localhost:6060/debug/pprof/goroutine
	// http://localhost:6060/debug/pprof/heap
	// http://localhost:6060/debug/pprof/mutex
	// http://localhost:6060/debug/pprof/profile
	// http://localhost:6060/debug/pprof/threadcreate
	// http://localhost:6060/debug/pprof/trace
	fmt.Println(http.ListenAndServe("localhost:6060", nil))
}

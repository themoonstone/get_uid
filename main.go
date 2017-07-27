package main

import (
	"strconv"
	//"encoding/json"

	"net/http"

	"github.com/sony/sonyflake"

	"github.com/drone/routes"
)

var gl sonyflake.GlobalVal
var gs *sonyflake.GlobalVal

// Initalize the global config
// the default count of goroutine is 10
func InitId() {
	gl.Poolsize = 10
	gs = gl.NewGlobal(gl.Poolsize)
	gs.GenId()
}

func main() {
	InitId()

	mux := routes.New()
	mux.Get("/user/:uid", getuid)
	http.Handle("/", mux)
	http.ListenAndServe(":8080", nil)

}

// the restful api
// call it by curl localhost:port/user/getuid
func getuid(w http.ResponseWriter, r *http.Request) {
	id, err := gs.GetId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//	body, err := json.Marshal(sonyflake.Decompose(id))

	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	s := strconv.Itoa(int(id))
	w.Write([]byte(s))
}

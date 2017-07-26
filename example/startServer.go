package main

import (
	"encoding/json"

	"net/http"

	"github.com/sony/sonyflake"

	"github.com/drone/routes"
)

var gl sonyflake.GlobalVal
var gs *sonyflake.GlobalVal

func InitId() {
	gl.Poolsize = 1
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

func getuid(w http.ResponseWriter, r *http.Request) {
	id, err := gs.GetId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(sonyflake.Decompose(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header()["Content-Type"] = []string{"application/json; charset=utf-8"}
	w.Write(body)
}

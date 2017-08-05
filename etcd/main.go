package main

import (
	"get_uid/pkg/etcd"
	"net/http"

	"github.com/drone/routes"
	"fmt"
)

var etcd_id *etcd.EtcdConfig

//initialize the etcd config
func InitId() {

	ip := "http://0.0.0.0:4001"
	var slice []string = make([]string, 0)
	slice = append(slice, ip)
	etcd_id = &etcd.EtcdConfig{}
	etcd_id.NewEtcdConfig(slice)
}
//
func main() {

	InitId()
	mux := routes.New()
	mux.Get("/user/:etcuid", etcuid)
	http.Handle("/", mux)
	http.ListenAndServe(":9090", nil)
	//fmt.Println("hello")
}
//
// the restful api
// call it by curl localhost:port/user/etcuid
func etcuid(w http.ResponseWriter, r *http.Request) {
	id, err := etcd_id.GetUid()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := fmt.Sprintf("id:%s\n", id)

	w.Write([]byte(s))
}

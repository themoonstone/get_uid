package main

import (
	"fmt"
	"get_uid/pkg/etcd"
	"net/http"

	"github.com/drone/routes"
)

var etcd_id *etcd.EtcdConfig

//initialize the etcd config
func InitId() {
	//	var ip *string
	//	ip = flag.String("IP", "http://127.0.0.1:9099", "the ip of etcd service installed")
	//	flag.Parse()
	//	var validIP = regexp.MustCompile(`^([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3}):[0-9]{1,4}$`)

	//	if validIP.MatchString(*ip) {
	//		*ip = "http://" + *ip
	//	} else {
	//		panic(errors.New(fmt.Sprintf("format error. [Usage ip format]:127.0.0.1")))
	//	}
	ip := "http://0.0.0.0:4001"
	var slice []string = make([]string, 0)
	slice = append(slice, ip)
	etcd_id = &etcd.EtcdConfig{}
	etcd_id.NewEtcdConfig(slice)
}

func main() {

	InitId()
	mux := routes.New()
	mux.Get("/user/:etcuid", etcuid)
	http.Handle("/", mux)
	http.ListenAndServe(":9090", nil)
}

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

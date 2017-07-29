package etcd

import (
	"log"
	"testing"
)

func TestGetUid(t *testing.T) {
	etcd_test := &EtcdConfig{}
	etcd_test.NewEtcdConfig(nil)
	log.Println(etcd_test.GetUid())
}

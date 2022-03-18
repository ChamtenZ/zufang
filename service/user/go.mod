module user

go 1.15

require (
	github.com/asim/go-micro/plugins/registry/consul/v3 v3.7.0
	github.com/asim/go-micro/v3 v3.7.1
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/gomodule/redigo v1.8.8
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/miekg/dns v1.1.47 // indirect
	github.com/tedcy/fdfs_client v0.0.0-20200106031142-21a04994525a // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20220310020820-b874c991c1a5 // indirect
	golang.org/x/tools v0.1.9 // indirect
	google.golang.org/protobuf v1.27.1
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

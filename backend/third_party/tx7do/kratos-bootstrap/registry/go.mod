module github.com/tx7do/kratos-bootstrap/registry

go 1.24.6

replace github.com/tx7do/kratos-bootstrap/api => ../api

require (
	github.com/go-kratos/kratos/v3 v3.0.0
	github.com/tx7do/kratos-bootstrap/api v0.0.44
)

require google.golang.org/protobuf v1.36.12 // indirect

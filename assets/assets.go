package assets

//go:generate rm bindata.go
//go:generate go-bindata -nocompress -pkg=assets ./...
//go:generate gofmt -s -w .

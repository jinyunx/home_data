go run genkey.go
go build oversea.go key.go aes.go
go build local.go key.go aes.go
rm key.go
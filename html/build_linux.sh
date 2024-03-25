go run genkey.go
env GOOS=linux GOARCH=amd64 go build oversea.go key.go aes.go
env GOOS=linux GOARCH=arm GOARM=6 go build local.go key.go aes.go
rm key.go
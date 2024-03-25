package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	key := genKey()
	content := fmt.Sprintf("package main\n\nvar key = \"%s\"", key)
	os.WriteFile("key.go", []byte(content), os.ModePerm)
}

func genKey() string {
	const length = 16
	str := "~`!@#$%^&*()_+=-{}[]|;':,./<>?0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

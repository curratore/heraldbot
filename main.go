package main

const (
	VERSION = "0.1"
)

func main() {
	h, err := NewHerald("TOKEN")
	if err != nil {
		panic(err)
	}
	h.Run()
}

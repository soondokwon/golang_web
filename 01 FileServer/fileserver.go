package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func HelloWorld() {
	var input string

	fmt.Printf("Enter your name : ")
	fmt.Scanln(&input)
	fmt.Printf("Hello, %s\n", input)
}

func MakeFileServer(inPort, inDir string) {
	port := flag.String("p", inPort, "port") // port의 주소
	dir := flag.String("d", inDir, "dir")    // dir의 주소
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*dir)))
	log.Printf("Server dir=[%s], port=[%s]\n", *dir, *port)
	listenInput := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(listenInput, nil))
}

func main() {
	// HelloWorld()
	MakeFileServer("8080", ".")
}

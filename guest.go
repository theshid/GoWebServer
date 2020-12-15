package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type GuestBook struct {
	SignatureCount int
	Signatures     []string
}

func getStrings(filename string) []string {
	var lines []string
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil
	}
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	check(scanner.Err())
	return lines
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	signatures := getStrings("signatures.txt")
	guestBook := GuestBook{
		SignatureCount: len(signatures),
		Signatures:     signatures,
	}
	html, err := template.ParseFiles("view.html")
	err = html.Execute(writer, guestBook)
	check(err)
}

func newHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("new.html")
	check(err)
	err = html.Execute(writer, nil)
	check(err)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	signature := request.FormValue("signature")
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("signatures.txt", options, os.FileMode(int(0777)))
	check(err)
	_, err = fmt.Fprintln(file, signature)
	check(err)
	err = file.Close()
	check(err)
	http.Redirect(writer, request, "/guest", http.StatusFound)

}

func main() {
	http.HandleFunc("/guest", viewHandler)
	http.HandleFunc("/guest/new", newHandler)
	http.HandleFunc("/guest/create", createHandler)
	err := http.ListenAndServe("localhost:9080", nil)
	log.Fatal(err)
}

package main

import (
	"fmt"
	"io/ioutil"
)

// DATA STRUCTURES

type Page struct {
	Title string

	// Why not of type "string"? Because it's what expected by io libraries
	// see `godoc io/ioutil WriteFile`
	Body []byte
}

// This is a function with a "receiver" and that is, p, a pointer to Page.
// This returns a value of type error.
// Why a pointer in receiver? Arguments are by default passed by value.
// Having a pointer saves unnecessary memory expense creating a clone.
func (p *Page) save() error {
	filename := p.Title + ".txt"

	// func WriteFile(filename string, data []byte, perm os.FileMode) error
	//	WriteFile writes data to a file named by filename. If the file does not
	//	exist, WriteFile creates it with permissions perm; otherwise WriteFile
	//	truncates it before writing.

	// What is "0600"? It is unix-style permissions. It means read-write
	// permission for the user only.
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Notice that this does not have a receiver, so it is just a regular function
// instead of being a method to an object/struct.
// Why is it returning a pointer of type Page? Again, like above, it is to
// pass the object instead of creating copies.
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"

	// func ReadFile(filename string) ([]byte, error)
	//	ReadFile reads the file named by filename and returns the contents. A
	//	successful call returns err == nil, not err == EOF. Because ReadFile
	//	reads the whole file, it does not treat an EOF from Read as an error to
	//	be reported.
	body, err := ioutil.ReadFile(filename)

	// Error handling.
	if err != nil {
		return nil, err
	}

	// Notice the dereferencer character, "&".
	// This returns a "tuple" of type (*Page, error)
	return &Page{Title: title, Body: body}, nil
}

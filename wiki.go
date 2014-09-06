package main

import (
	// String format package. Printing, etc.
	"fmt"

	// This package is for templating, for keeping HTML in a separate file,
	// allowing us to change our pages without modifying underlying Go code.
	"html/template"

	// File Handling package.
	"io/ioutil"

	// Networking. Creating server and related methods.
	"net/http"
)

// ENTRY POINT

func main() {
	// Instantiation of new Page
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page")}
	p1.save()

	// Remember that it returns a "tuple" (*Page, error)?
	// Here, we just ignore the second character by using "_"
	p2, _ := loadPage("TestPage")

	// string(p2.Body) casts []byte into string
	fmt.Println(string(p2.Body))

	// SERVER

	// func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
	//	HandleFunc registers the handler function for the given pattern in the
	//	DefaultServeMux. The documentation for ServeMux explains how patterns
	//	are matched.
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	//http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}

// DATA STRUCTURES & METHODS

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

// FUNCTION HELPERS

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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	// func (t *Template) ParseFiles(filenames ...string) (*Template, error)
	//	ParseFiles parses the named files and associates the resulting templates
	//	with t. If an error occurs, parsing stops and the returned template is
	//	nil; otherwise it is t. There must be at least one file.

	// Note that this returns a *template.Template
	t, _ := template.ParseFiles(tmpl + ".html")

	// func (t *Template) Execute(wr io.Writer, data interface{}) error
	//	Execute applies a parsed template to the specified data object, writing
	//	the output to wr. If an error occurs executing the template or writing
	//	its output, execution stops, but partial results may already have been
	//	written to the output writer. A template may be executed safely in
	//	parallel.
	t.Execute(w, p)
}

// http://golang.org/pkg/net/http/#ResponseWriter
// http://golang.org/pkg/net/http/#Request
func viewHandler(w http.ResponseWriter, r *http.Request) {
	// This takes all the characters after "/view/" and assigns it to "title"
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)

	// func Fprintf(w io.Writer, format string, a ...interface{})
	// (n int, err error)
	//	Fprintf formats according to a format specifier and writes to w. It
	//	returns the number of bytes written and any write error encountered.

	// http.ResponseWriter implements the Writer interface so it's possible
	// to "stream" the results to it.
	// type Writer interface { Write(p []byte) (n int, err error) }

	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

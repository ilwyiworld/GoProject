package main

import (
	"html/template"
	"os"
)

type Friend struct {
	Fname  string
}

type Person struct {
	UserName string
	Emails   []string
	Friends  []*Friend
}

func main() {
	t := template.New("fieldname example")
	t, _ = t.Parse("hello {{.Fname}}!")
	p := Friend{Fname : "yiworld"}
	t.Execute(os.Stdout, p)

	f1 := Friend{Fname: "minux.ma"}
	f2 := Friend{Fname: "xushiwei"}
	t1 := template.New("fieldname example")
	t1, _ = t1.Parse(`hello {{.UserName}}!
			{{range .Emails}}
				an email {{.}}
			{{end}}
			{{with .Friends}}
			{{range .}}
				my friend name is {{.Fname}}
			{{end}}
			{{end}}
			`)
	p1 := Person{UserName: "yiworld",
		Emails:  []string{"astaxie@beego.me", "astaxie@gmail.com"},
		Friends: []*Friend{&f1, &f2}}
	t1.Execute(os.Stdout, p1)

}
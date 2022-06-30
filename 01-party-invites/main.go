package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Rsvp struct {
	Name, Email, Phone string
	WillAttend         bool
}

type formData struct {
	*Rsvp
	Errors []string
}

var responses = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
	templateNames := [5]string{"welcome", "invite-form", "thanks", "sorry", "guest-list"}

	for index, templateName := range templateNames {
		template, err := template.ParseFiles("layout.html", templateName+".html")

		if err == nil {
			templates[templateName] = template
			fmt.Println("Загруженный шаблон:", index, templateName)
		} else {
			panic(err)
		}
	}
}

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, nil)
}

func inviteformhandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templates["invite-form"].Execute(writer, formData{
			Rsvp: &Rsvp{}, Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		responseData := Rsvp{
			Name:       request.Form["name"][0],
			Email:      request.Form["email"][0],
			Phone:      request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}
		errors := []string{}

		if responseData.Name == "" {
			errors = append(errors, "Введите Ваше Имя")
		}

		if responseData.Email == "" {
			errors = append(errors, "Введите Ваш Email")
		}

		if responseData.Phone == "" {
			errors = append(errors, "Введите Ваш телефонный номер")
		}

		if len(errors) > 0 {
			templates["invite-form"].Execute(writer, formData{
				Rsvp: &responseData, Errors: errors,
			})
		} else {
			responses = append(responses, &responseData)
			if responseData.WillAttend {
				templates["thanks"].Execute(writer, responseData.Name)
			} else {
				templates["sorry"].Execute(writer, responseData.Name)
			}
		}

	}
}

func guestlistHandler(writer http.ResponseWriter, request *http.Request) {
	templates["guest-list"].Execute(writer, responses)
}

func main() {
	fmt.Println("Hello, Go!")
	loadTemplates()
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/invite-form", inviteformhandler)
	http.HandleFunc("/guest-list", guestlistHandler)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Println(err)
	}
}

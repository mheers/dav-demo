package main

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {

	a := auth{
		username: "marcel",
		password: "password",
	}

	addressBookEntries := []User{
		{
			Firstname: "Marcel",
			Lastname:  "Heers",
			Emails: map[string]string{
				"work": "info@heers.it",
			},
			Phones: map[string]string{
				"work": "+491623300432",
			},
			Unit: "IT",
		},
	}

	events := []CalendarEvent{
		{
			ID:          "123",
			Summary:     "Meeting",
			Location:    "Room 101",
			StartAt:     time.Now(),
			EndAt:       time.Now().Add(time.Hour),
			CreatedAt:   time.Now(),
			Description: "Discuss the new project",
		},
	}

	r := http.NewServeMux()
	r.Handle("/carddav/", a.middleware(NewCardDAVHandler("carddav", addressBookEntries)))
	r.Handle("/caldav/", a.middleware(NewCalDavHandler("caldav", events)))
	r.Handle("/files/", a.middleware(NewWebDavHandler("files", "./files")))

	s := &http.Server{
		Addr:           ":8086",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logrus.Fatal(s.ListenAndServe())
}

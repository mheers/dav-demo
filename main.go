package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// TracingMiddleware logs incoming requests
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start := time.Now()

		// Log the incoming request details
		logrus.Infof("Request Method: %s, URL: %s, RemoteAddr: %s, User-Agent: %s",
			r.Method, r.URL.Path, r.RemoteAddr, r.Header.Get("User-Agent"))

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the duration it took to handle the request
		// logrus.Infof("Completed in %s", time.Since(start))
	})
}

func main() {

	a := auth{
		username: "marcel",
		password: "admin",
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

	r := mux.NewRouter()
	r.Use(TracingMiddleware)

	r.Handle("/carddav/", a.middleware(NewCardDAVHandler("carddav", addressBookEntries)))
	r.Handle("/caldav/", a.middleware(NewCalDavHandler("caldav", events)))
	r.PathPrefix("/files").Handler(a.middleware(NewWebDavHandler("files", "./files")))

	s := &http.Server{
		Addr:           ":8086",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logrus.Fatal(s.ListenAndServe())
}

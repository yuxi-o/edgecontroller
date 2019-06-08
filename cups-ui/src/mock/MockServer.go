package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const addr = ":8090"

// Userplane is a CUPS userplane.
type Userplane struct {
	ID           string        `json:"id,omitempty"`
	UUID         string        `json:"uuid,omitempty"`
	Function     string        `json:"function,omitempty"`
	Config       interface{}   `json:"config,omitempty"`
	Selectors    []interface{} `json:"selectors,omitempty"`
	Entitlements []interface{} `json:"entitlements,omitempty"`
}

func main() { // nolint: gocyclo
	var userplanes []Userplane

	handler := http.NewServeMux()

	handler.HandleFunc("/userplanes/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/userplanes/")

		index, userplane, err := func() (int, *Userplane, error) {
			for i, up := range userplanes {
				if up.ID == id {
					return i, &up, nil
				}
			}
			return 0, nil, errors.New("not found")
		}()
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
		}

		switch r.Method {
		case "GET":
			bytes, err := json.Marshal(&userplane)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
			_, err = w.Write(bytes)
			if err != nil {
				log.Fatal(err)
			}

		case "PATCH":
			var u Userplane
			if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
				fmt.Printf("%+v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			userplanes[index] = u
			w.WriteHeader(http.StatusOK)

		case "DELETE":
			userplanes = append(userplanes[:index], userplanes[index+1:]...)
			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "", http.StatusBadRequest)
		}
	})

	handler.HandleFunc("/userplanes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var u Userplane
			if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
				fmt.Printf("%+v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			u.ID = fmt.Sprintf("%d", r.Int())

			userplanes = append(userplanes, u)

			w.WriteHeader(http.StatusCreated)
			return
		}

		bytes, err := json.Marshal(userplanes)
		if err != nil {
			log.Fatal(err)
		}

		_, err = w.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
	})

	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	log.Fatal(s.ListenAndServe())
}

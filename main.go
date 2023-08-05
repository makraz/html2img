package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type preview struct {
	Name      string `json:"name" validate:"required"`
	Url       string `json:"url" validate:"required"`
	ImageName string `json:"image"`
}

type Exception struct {
	Message string `json:"message"`
}

type Response struct {
	Data string `json:"data"`
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))

// var path = "/var/www/html2png"
var path = "public"

func homeLink(w http.ResponseWriter, r *http.Request) {
	fprintf, err := fmt.Fprintf(w, "Welcome home!")
	fmt.Println(fprintf)
	if err != nil {
		return
	}
}

func createPreview(w http.ResponseWriter, r *http.Request) {
	var newPreview preview

	reqBody, err := io.ReadAll(r.Body)

	if err != nil {
		fprintf, err := fmt.Fprintf(w, "Kindly enter data with the preview title and description only in order to update")
		fmt.Println(fprintf)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	errr := json.Unmarshal(reqBody, &newPreview)
	if errr != nil {
		fmt.Println(errr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithDebugf(log.Printf),
	)
	// defer cancel()

	var buf []byte
	var selector = "#preview"

	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(newPreview.Url),
		chromedp.Screenshot(selector, &buf, chromedp.NodeVisible),
	}); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cancel()

	now := strconv.FormatInt(time.Now().Unix(), 10)
	newPreview.ImageName = newPreview.Name + "." + now + ".png"
	fmt.Println(newPreview.ImageName)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			log.Println(err)
		}
	}

	if err := os.WriteFile(path+"/"+newPreview.ImageName, buf, 0o644); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("wrote " + path + "/" + newPreview.ImageName + ".png")
	log.Printf("wrote " + path + "/" + newPreview.ImageName + ".png")

	w.WriteHeader(http.StatusCreated)

	errrr := json.NewEncoder(w).Encode(newPreview)
	if err != nil {
		fmt.Println(errrr)
		log.Println(errrr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	return
}

func getImage(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	fileBytes, err := os.ReadFile(path + "/" + name)

	if err != nil {
		fmt.Println(err)
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
		// panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")

	write, err := w.Write(fileBytes)
	fmt.Println(write)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	return
}

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	err := json.NewEncoder(w).Encode(map[string]string{"Status": "404", "Message": "This endpoint not found"})
	if err != nil {
		return
	}

	return
}

//type event struct {
//	ID          string `json:"ID"`
//	Title       string `json:"Title"`
//	Description string `json:"Description"`
//}

//type allEvents []event

//var events = allEvents{
//	{
//		ID:          "1",
//		Title:       "Introduction to Golang",
//		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
//	},
//}

//func createEvent(w http.ResponseWriter, r *http.Request) {
//	var newEvent event
//	reqBody, err := io.ReadAll(r.Body)
//	if err != nil {
//		fprintf, err := fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
//		fmt.Println(fprintf)
//		if err != nil {
//			return
//		}
//	}
//
//	errr := json.Unmarshal(reqBody, &newEvent)
//	fmt.Println(errr)
//	if err != nil {
//		return
//	}
//	events = append(events, newEvent)
//	w.WriteHeader(http.StatusCreated)
//
//	errrr := json.NewEncoder(w).Encode(newEvent)
//	fmt.Println(errrr)
//	if err != nil {
//		return
//	}
//}

//func getOneEvent(w http.ResponseWriter, r *http.Request) {
//	eventID := mux.Vars(r)["id"]
//
//	for _, singleEvent := range events {
//		if singleEvent.ID == eventID {
//			err := json.NewEncoder(w).Encode(singleEvent)
//			if err != nil {
//				return
//			}
//		}
//	}
//}

//func getAllEvents(w http.ResponseWriter, r *http.Request) {
//	err := json.NewEncoder(w).Encode(events)
//	if err != nil {
//		return
//	}
//}

//func updateEvent(w http.ResponseWriter, r *http.Request) {
//	eventID := mux.Vars(r)["id"]
//	var updatedEvent event
//
//	reqBody, err := io.ReadAll(r.Body)
//	if err != nil {
//		fprintf, err := fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
//
//		fmt.Println(fprintf)
//
//		if err != nil {
//			return
//		}
//	}
//
//	errr := json.Unmarshal(reqBody, &updatedEvent)
//	fmt.Println(errr)
//	if errr != nil {
//		return
//	}
//
//	for i, singleEvent := range events {
//		if singleEvent.ID == eventID {
//			singleEvent.Title = updatedEvent.Title
//			singleEvent.Description = updatedEvent.Description
//			events = append(events[:i], singleEvent)
//			err := json.NewEncoder(w).Encode(singleEvent)
//			if err != nil {
//				return
//			}
//		}
//	}
//}

//func deleteEvent(w http.ResponseWriter, r *http.Request) {
//	eventID := mux.Vars(r)["id"]
//
//	for i, singleEvent := range events {
//		if singleEvent.ID == eventID {
//			events = append(events[:i], events[i+1:]...)
//			fprintf, err := fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
//
//			fmt.Println(fprintf)
//
//			if err != nil {
//				return
//			}
//		}
//	}
//}

func IsAuthorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")

		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						w.WriteHeader(http.StatusUnauthorized)
						//return nil, fmt.Errorf("there was an error")
						return nil, fmt.Errorf("Invalid authorization token")
					}
					return JwtKey, nil
				})

				if error != nil {
					w.WriteHeader(http.StatusUnauthorized)
					//json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
					//return
				}

				if token.Valid {
					next.ServeHTTP(w, r)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
					//return
				}
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			//json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
			json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/preview", IsAuthorized(createPreview)).Methods("POST")
	router.HandleFunc("/image/{name}", getImage).Methods("GET")
	//router.HandleFunc("/health-check", HealthCheck).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(NotFoundPage)
	router.MethodNotAllowedHandler = http.HandlerFunc(NotFoundPage)

	//router.HandleFunc("/event", createEvent).Methods("POST")
	//router.HandleFunc("/events", getAllEvents).Methods("GET")
	//router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	//router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	//router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	mw "schoolManagement/internal/api/middlewares"
	"time"
)

// Even though the user struct is private (not starting with uppercase), the field values after made public (Name, Age and City).
// This is because, while unmarshalling the field values are accessed by another package (encoding/json), so instead of struct, the field values are made public;
type user struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

// rootHandler - Handler for root route;
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	_, err := w.Write([]byte("Hello Root Route"))
	if err != nil {
		fmt.Println("Error from root handler ", err)
	}
}

// studentsHandler - Handler for students route;
func studentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write([]byte("Hello this is a GET method call to students api!"))
		if err != nil {
			fmt.Println("Error from students handler ", err)
		}
		return
	case http.MethodPost:
		_, err := w.Write([]byte("Hello this is a POST method call to students api!"))
		if err != nil {
			fmt.Println("Error from students handler ", err)
		}
		return
	case http.MethodPatch:
		_, err := w.Write([]byte("Hello this is a PATCH method call to students api!"))
		if err != nil {
			fmt.Println("Error from students handler ", err)
		}
		return
	case http.MethodPut:
		_, err := w.Write([]byte("Hello this is a PUT method call to students api!"))
		if err != nil {
			fmt.Println("Error from students handler ", err)
		}
		return
	case http.MethodDelete:
		_, err := w.Write([]byte("Hello this is a DELETE method call to students api!"))
		if err != nil {
			fmt.Println("Error from students handler ", err)
		}
		return
	default:
		var msg string = "This a " + r.Method + " call to students"
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error from students handler ", err)
		}
		return
	}
}

// teachersHandler - Handler for teachers route;
func teachersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write([]byte("Hello this is a GET method call to teachers api!"))
		if err != nil {
			fmt.Println("Error from teachers handler ", err)
		}
		return
	case http.MethodPost:
		// Parsing form post data;
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
		}

		fmt.Println("Form : ", r.Form)

		response := make(map[string]interface{})
		for key, value := range r.Form {
			response[key] = value[0]
		}

		// Processing raw body;
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusBadRequest)
		}
		defer func() {
			err := r.Body.Close()
			if err != nil {
				fmt.Println("Error closing body")
			}
		}()

		fmt.Println("Processed response (RAW): ", string(body))
		fmt.Println("Processed response (X-WWW-FORM): ", response)

		// If you expect json data then unmarshal it;
		var userInstance user
		err = json.Unmarshal(body, &userInstance)
		if err != nil {
			http.Error(w, "Error parsing body", http.StatusBadRequest)
			return
		}

		fmt.Println("User instance after unmarshalling ", userInstance)

		fmt.Println("Body : ", r.Body)
		fmt.Println("Form : ", r.Form)
		fmt.Println("Header : ", r.Header)
		fmt.Println("Context : ", r.Context())
		fmt.Println("Content length : ", r.ContentLength)
		fmt.Println("Host  : ", r.Host)
		fmt.Println("Method  : ", r.Method)
		fmt.Println("Protocol  : ", r.Proto)
		fmt.Println("Remote address : ", r.RemoteAddr)
		fmt.Println("Request URI : ", r.RequestURI)
		fmt.Println("Request TLS : ", r.TLS)
		fmt.Println("Trailer : ", r.Trailer)
		fmt.Println("Transfer Encoding : ", r.TransferEncoding)
		fmt.Println("User agent : ", r.UserAgent())
		fmt.Println("Port : ", r.URL.Port())

		_, err = w.Write([]byte("Hello this is a POST method call to teachers api!"))
		if err != nil {
			fmt.Println("Error from teachers handler ", err)
		}
		return
	case http.MethodPatch:
		_, err := w.Write([]byte("Hello this is a PATCH method call to teachers api!"))
		if err != nil {
			fmt.Println("Error from teachers handler ", err)
		}
		return
	case http.MethodPut:
		_, err := w.Write([]byte("Hello this is a PUT method call to teachers api!"))
		if err != nil {
			fmt.Println("Error from teachers handler ", err)
		}
		return
	case http.MethodDelete:
		_, err := w.Write([]byte("Hello this is a DELETE method call to teachers api!"))
		if err != nil {
			fmt.Println("Error from teachers handler ", err)
		}
		return
	default:
		var msg string = "This a " + r.Method + " call to teachers"
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error from teachers handler ", err)
		}
		return
	}

}

// execsHandler - Handler for execs route;
func execsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write([]byte("Hello this is a GET method call to execs api!"))
		if err != nil {
			fmt.Println("Error from execs handler ", err)
		}
		return
	case http.MethodPost:
		form := r.Form
		fmt.Println("Form : ", form)
		_, err := w.Write([]byte("Hello this is a POST method call to execs api!"))
		if err != nil {
			fmt.Println("Error from execs handler ", err)
		}
		return
	case http.MethodPatch:
		_, err := w.Write([]byte("Hello this is a PATCH method call to execs api!"))
		if err != nil {
			fmt.Println("Error from execs handler ", err)
		}
		return
	case http.MethodPut:
		_, err := w.Write([]byte("Hello this is a PUT method call to execs api!"))
		if err != nil {
			fmt.Println("Error from execs handler ", err)
		}
		return
	case http.MethodDelete:
		_, err := w.Write([]byte("Hello this is a DELETE method call to execs api!"))
		if err != nil {
			fmt.Println("Error from execs handler ", err)
		}
		return
	default:
		var msg string = "This a " + r.Method + " call to execs"
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error from execs handler ", err)
		}
		return
	}
}

func main() {
	const port string = ":3000"

	cert := "cert.pem"
	key := "key.pem"

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/students", studentsHandler)
	mux.HandleFunc("/teachers/", teachersHandler)
	mux.HandleFunc("/execs", execsHandler)

	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}

	rl := mw.NewRateLimiter(5, time.Minute)

	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		WhiteList:                   []string{"name", "age"},
	}

	// All the other middlewares are passed as argument;
	secureMux := mw.Hpp(hppOptions)(rl.Middleware(mw.CompressionMiddleware(mw.ResponseTimeMiddleware(mw.SecurityHandler(mw.Cors(mux))))))

	server := &http.Server{
		Addr:      port,
		TLSConfig: tlsConfig,
		Handler:   secureMux,
	}
	fmt.Println("Server is running on port ", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatal("Error starting the server : ", err)
	}
}

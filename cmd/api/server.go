package main

import (
	"crypto/tls"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	mw "schoolManagement/internal/api/middlewares"
	"schoolManagement/internal/api/routers"
	"schoolManagement/internal/repositories/sqlconnect"
)

// Even though the user struct is private (not starting with uppercase), the field values after made public (Name, Age and City).
// This is because, while unmarshalling the field values are accessed by another package (encoding/json), so instead of struct, the field values are made public;
type user struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func main() {
	// The godotenv.Load() makes sure that the .env variables are picked instead of system env variable by the os.Getenv()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var port = os.Getenv("API_PORT")
	_, err = sqlconnect.ConnectDb()
	if err != nil {
		log.Fatal(err)
	}

	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}

	// Should be uncommented while using all middlewares;
	//rl := mw.NewRateLimiter(5, time.Minute)
	//
	//hppOptions := mw.HPPOptions{
	//	CheckQuery:                  true,
	//	CheckBody:                   true,
	//	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	//	WhiteList:                   []string{"name", "age"},
	//}

	// All the other middlewares are passed as argument;
	//secureMux := mw.Cors(rl.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHandler(mw.CompressionMiddleware(mw.Hpp(hppOptions)(mw.Cors(mux)))))))
	// An enhanced and efficient way to apply middlewares;
	//secureMux := applyMiddleWares(mux, mw.Hpp(hppOptions), mw.CompressionMiddleware, mw.SecurityHandler, mw.ResponseTimeMiddleware, rl.Middleware, mw.Cors)

	// For this server we will use mw.SecurityHandler alone now;
	router := routers.MainRouter()
	secureMux := mw.JWTMiddleware(mw.SecurityHandler(router))
	//secureMux := (mw.SecurityHandler(router))

	server := &http.Server{
		Addr:      port,
		TLSConfig: tlsConfig,
		Handler:   secureMux,
	}
	fmt.Println("Server is running on port ", port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatal("Error starting the server : ", err)
	}
}

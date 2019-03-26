package main

import (
	"net"
	"github.com/oschwald/geoip2-golang"
	"log"
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	handlePages()
	log.Println("Started Server..")
	panic(http.ListenAndServe(":9000", nil))
}

const (
	API    = "/api"
	ApiKey ="YOUR_API_KEY"
)
func handlePages() {
	http.HandleFunc(API+"/city", city)
	http.HandleFunc(API+"/country", country)
	http.HandleFunc(API+"/postalCode", postalCode)
	http.HandleFunc(API+"/inEuropeanUnion", inEuropeanUnion)
	http.HandleFunc(API+"/longitude", longitude)
	http.HandleFunc(API+"/latitude", latitude)
	http.HandleFunc(API+"/timeZone", timeZone)
	http.HandleFunc(API+"/countryCode", countryCode)
	http.HandleFunc(API+"/continentCode", continentCode)
	http.HandleFunc(API+"/continentName", continentName)
}

func APIKeyError(text string) error {
	return &apiKeyErr{text}
}

type apiKeyErr struct {
	s string
}

func (e *apiKeyErr) Error() string {
	return e.s
}

func validate(request *http.Request, writer http.ResponseWriter, ) (*geoip2.City, error) {
	if request.FormValue("api-key") != ApiKey {
		fmt.Fprintf(writer, "Invalid API Key! Please contact an administrator.")
		return new(geoip2.City), APIKeyError("Invalid API Key! Please contact an administrator.")
	}
	db, err := geoip2.Open("geodb.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ip := request.FormValue("ip")
	record, err := db.City(net.ParseIP(ip))
	log.Printf("REQUEST INCOME FROM - %s",ip)
	return record, err
}

// You should probably add some error checking on the validate function

func country(writer http.ResponseWriter, request *http.Request) {
	record, _ := validate(request, writer)
	lang := request.FormValue("lang")
	if lang == "" {
		lang = "en"
	}
	fmt.Fprintf(writer, record.Country.Names[lang])
}

func city(writer http.ResponseWriter, request *http.Request) {
	record, _ := validate(request, writer)
	lang := request.FormValue("lang")
	if lang == "" {
		lang = "en"
	}
	fmt.Fprintf(writer, record.City.Names[lang])
}

func continentName(writer http.ResponseWriter, request *http.Request) {
	record, _ := validate(request, writer)
	lang := request.FormValue("lang")
	if lang == "" {
		lang = "en"
	}
	fmt.Fprintf(writer, record.Continent.Names[lang])
}

func postalCode(writer http.ResponseWriter, request *http.Request) {
	record, _ := validate(request, writer)
	fmt.Fprintf(writer, record.Postal.Code)
}

func inEuropeanUnion(writer http.ResponseWriter, request *http.Request) {
	record, _ := validate(request, writer)
	fmt.Fprintf(writer, strconv.FormatBool(record.RegisteredCountry.IsInEuropeanUnion))
}

func longitude(writer http.ResponseWriter, request *http.Request) {
	record, _ := validate(request, writer)
	fmt.Fprintf(writer, strconv.FormatFloat(record.Location.Longitude, 'f', 2, 64))
}

func latitude(writer http.ResponseWriter, request *http.Request) {
	record, _ := validate(request, writer)
	fmt.Fprintf(writer, strconv.FormatFloat(record.Location.Latitude, 'f', 2, 64))
}

func timeZone(writer http.ResponseWriter, request *http.Request) {
	record, _ := validate(request,writer)
	fmt.Fprintf(writer, record.Location.TimeZone)
}

func countryCode(writer http.ResponseWriter, request *http.Request) {
	record, _ := validate(request, writer)
	fmt.Fprintf(writer, record.Country.IsoCode)
}

func continentCode(writer http.ResponseWriter, request *http.Request) {
	record, _ := validate(request, writer)
	fmt.Fprintf(writer, record.Continent.Code)
}
package main

import (
	"encoding/json"
	"fmt"
	utils "github.com/TrafficLabel/Go-Utilities"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
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
	DB = "geodb.mmdb"
)
func handlePages() {
	http.HandleFunc(API, all)
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

type JsonResult struct {
	IP string
	Country string
	CountryGeoNameID uint
	CountryISOCode string
	City string
	CityGeoNameID uint
	IsEU bool
	TimeZone string
	Continent string
	ContinentCode string
	ContinentGeoNameID uint
	Latitude float64
	Longitude float64
	ZipCode string
	AccuracyRadius uint16
	MetroCode uint
	IsAnonymousProxy bool
	IsSatelliteProvider bool
}

func APIKeyError(writer http.ResponseWriter, text string) error {
	fmt.Fprintf(writer, "Invalid API Key! Please contact an administrator.")
	return &apiKeyErr{text}
}

type apiKeyErr struct {
	s string
}

func (e *apiKeyErr) Error() string {
	return e.s
}

func handleError(err error) {
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func validate(writer http.ResponseWriter, request *http.Request) (*geoip2.City, error, string) {
	if request.FormValue("api-key") != ApiKey {
		return new(geoip2.City), APIKeyError(writer, "Invalid API Key! Please contact an administrator."), ""
	}
	db, err := geoip2.Open(DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ip := request.FormValue("ip")
	if ip == "" {
		ip = utils.GetRealAddr(request)
	}
	record, err := db.City(net.ParseIP(ip))
	if err != nil {
		handleError(err)
	}
	log.Printf("REQUEST INCOME FROM - %s", ip)
	lang := request.FormValue("lang")
	if lang == "" {
		lang = "en"
	}
	return record, nil, lang
}

func all(writer http.ResponseWriter, request *http.Request) {
	result, err := parseAll(writer, request)
	if err != nil {
		log.Println(err.Error())
	} else {
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(&result)
	}
}

func country(writer http.ResponseWriter, request *http.Request) {
	record, _, lang := validate(writer, request)
	fmt.Fprintf(writer, record.Country.Names[lang])
}

func city(writer http.ResponseWriter, request *http.Request) {
	record, _, lang := validate(writer, request)
	fmt.Fprintf(writer, record.City.Names[lang])
}

func continentName(writer http.ResponseWriter, request *http.Request) {
	record, _, lang := validate(writer, request)
	fmt.Fprintf(writer, record.Continent.Names[lang])
}

func postalCode(writer http.ResponseWriter, request *http.Request) {
	record, _, _ := validate(writer, request)
	fmt.Fprintf(writer, record.Postal.Code)
}

func inEuropeanUnion(writer http.ResponseWriter, request *http.Request) {
	record, _, _ := validate(writer, request)
	fmt.Fprintf(writer, strconv.FormatBool(record.RegisteredCountry.IsInEuropeanUnion))
}

func longitude(writer http.ResponseWriter, request *http.Request) {
	record, _, _ := validate(writer, request)
	fmt.Fprintf(writer, strconv.FormatFloat(record.Location.Longitude, 'f', -1, 64))
}

func latitude(writer http.ResponseWriter, request *http.Request) {
	record, _, _ := validate(writer, request)
	fmt.Fprintf(writer, strconv.FormatFloat(record.Location.Latitude, 'f', -1, 64))
}

func timeZone(writer http.ResponseWriter, request *http.Request) {
	record, _, _ := validate(writer, request)
	fmt.Fprintf(writer, record.Location.TimeZone)
}

func countryCode(writer http.ResponseWriter, request *http.Request) {
	record, _, _ := validate(writer, request)
	fmt.Fprintf(writer, record.Country.IsoCode)
}

func continentCode(writer http.ResponseWriter, request *http.Request) {
	record, _, _ := validate(writer, request)
	fmt.Fprintf(writer, record.Continent.Code)
}

func parseAll(writer http.ResponseWriter, request *http.Request) (*JsonResult, error) {
	record, err, lang := validate(writer, request)
	if err != nil {
		return nil, err
	}
	ip := request.FormValue("ip")
	if ip == "" {
		ip = utils.GetRealAddr(request)
	}
	country := record.Country
	city := record.City
	location := record.Location
	continent := record.Continent
	traits := record.Traits
	var result = JsonResult{
		IP:          ip,
		Country:	  country.Names[lang],
		CountryGeoNameID:	 country.GeoNameID,
		CountryISOCode:	 country.IsoCode,
		City:        city.Names[lang],
		CityGeoNameID:        city.GeoNameID,
		IsEU:        country.IsInEuropeanUnion,
		TimeZone:    location.TimeZone,
		Continent:   continent.Names[lang],
		ContinentCode: continent.Code,
		ContinentGeoNameID: continent.GeoNameID,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		ZipCode:     record.Postal.Code,
		AccuracyRadius: location.AccuracyRadius,
		MetroCode: location.MetroCode,
		IsAnonymousProxy: traits.IsAnonymousProxy,
		IsSatelliteProvider: traits.IsSatelliteProvider,
	}
	return &result, nil
}

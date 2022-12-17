package main

import (
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

type Objec struct {
	CityName    string
	CountryName string
	Latitude    float64
	Longitude   float64
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Ping test
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "healthz")
	})

	// Get user value
	r.GET("/ip/:ip", func(c *gin.Context) {
		var ipp = c.Params.ByName("ip")
		var ok = checkIPAddress(ipp)
		//valiIPV6 := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
		//checkIPAddress(valiIPV6)

		if ok {
			c.JSON(http.StatusOK, maxmindLookup(ipp))
		} else {
			c.JSON(http.StatusOK, gin.H{"error": "error"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func checkIPAddress(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	} else {
		return true
	}
}

func maxmindLookup(ipaddress string) Objec {
	//this assume you have a /database folder in this directory with
	//the appropriate GeoLite2-City.mmdb file available for download at
	//Maxmind
	db, err := geoip2.Open("./cmd/geolocate-ip/tools/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ip := net.ParseIP(ipaddress)
	record, err := db.City(ip)
	return Objec{CityName: record.City.Names["en"], CountryName: record.Country.Names["en"], Latitude: record.Location.Latitude, Longitude: record.Location.Longitude}

}

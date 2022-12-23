package main

import (
	"fmt"
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
	Ipaddress   string
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// Ping test
	r.GET("/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Get ClientIP address
	r.GET("/clientIP", func(c *gin.Context) {
		var ok = checkIPAddress(c.ClientIP())
		//Debugging related to issue where ClientIP did not return actual IP address.
		// https://github.com/gin-gonic/gin/issues/2697
		fmt.Println("The clientIP", c.ClientIP())
		//fmt.Println("The RemoteIP", c.RemoteIP())
		//fmt.Println("The RemoteAddress", c.Request.RemoteAddr)
		//fmt.Println("The Referer", c.Request.Referer())
		//fmt.Println("The Header", c.Request.Header)

		if ok {
			c.JSON(http.StatusOK, maxmindLookup(c.ClientIP()))
		} else {
			c.JSON(http.StatusOK, gin.H{"error": "error"})
		}
	})

	// Get user value
	r.GET("/ip/:ip", func(c *gin.Context) {
		var ipp = c.Params.ByName("ip")
		var ok = checkIPAddress(ipp)

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
	r.StaticFile("/favicon.ico", "./cmd/geolocate-ip/files/favicon.ico")
	r.StaticFile("/", "./cmd/geolocate-ip/files/index.html")
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

	return Objec{CityName: record.City.Names["en"], CountryName: record.Country.Names["en"], Latitude: record.Location.Latitude, Longitude: record.Location.Longitude, Ipaddress: ipaddress}

}

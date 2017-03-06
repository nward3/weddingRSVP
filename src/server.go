package main

import (
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"time"
)

type Rsvp struct {
	Name string `json: "name"`
	NumGuests int `json: "num_Guests"`
	IsAttending bool `json: "IsAttending"`
}

const (
	MongoDBHosts = "ds119370.mlab.com:19370"
	Database = "heroku_dqrb1b90"
	AuthUserName = "admin"
	AuthPassword = "testing123"
)

func main() {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: Database,
		Username: AuthUserName,
		Password: AuthPassword,
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	db := mongoSession.DB(Database)

	r := gin.Default()

	// send index.html so that the static site can be viewed
	r.StaticFS("/", http.Dir("client"))

	/* r.GET("/ring", func(c *gin.Context) { */
	/* 	c.JSON(200, gin.H{ */
	/* 		"message": "pong", */
	/* 	}) */
	/* }) */

	r.POST("/rsvp", func(c *gin.Context) {
		var rsvp Rsvp
		c.BindJSON(&rsvp)
		fmt.Println(rsvp.Name)

		rsvpCollection := db.C("rsvps")

		err = rsvpCollection.Insert(
			&Rsvp{
				Name: rsvp.Name,
				NumGuests: rsvp.NumGuests,
				IsAttending: rsvp.IsAttending})

		if err != nil {
			// handle error
			log.Fatal(err)
			c.JSON(400, gin.H{
				"error": true,
				"message": err,
			})
		} else {
			c.JSON(200, gin.H{
				"error": false,
				"message": "success",
			})
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

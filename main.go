package main

import (
	"context"
	"flag"
	"log"

	"github.com/AlexeyAndryushin/reservations/api"
	"github.com/AlexeyAndryushin/reservations/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userColl = "users"

var config = fiber.Config{
    ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.JSON(map[string]string{"error": err.Error()})
    },
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "The listend address of the API server")
	flag.Parse()


	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	//hanlers initialization
	var (
	userHandler = api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	hotelHandler = api.NewHotelHandler(hotelStore, roomStore)
  app = fiber.New(config)
	apiv1 = app.Group("/api/v1")
	)


	//user hanlders
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	//hotel handlers
	apiv1.Get("/hotel",hotelHandler.HanldeGetHotels)

  app.Listen(*listenAddr)
}

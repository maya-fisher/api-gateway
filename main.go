package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"example.com/m/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	pb "github.com/maya-fisher/birthday-service/proto"
	"google.golang.org/grpc"
)



func corsRouterConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AddExposeHeaders("x-uploadid")
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowWildcard = true
	corsConfig.AllowOrigins = strings.Split("http://localhost*", ",")
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders(
		"x-content-length",
		"authorization",
		"cache-control",
		"x-requested-with",
		"content-disposition",
		"content-range",
		"destination",
		"fileID",
	)

	return corsConfig
}

type person struct {
	name     string
	birthday int64
	userId   string
}

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := grpc.Dial(config.ADDRESS, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewBirthdaysClient(conn)

	r := gin.Default()

	r.Use(
		cors.New(corsRouterConfig()),
	)

	r.PUT("/birthday/:userId", func(c *gin.Context) {

		userId := c.Param("userId")

		person := &pb.Person{
			UserId: userId,
		}

		err = c.Bind(&person) //adds the given object from the body

		fmt.Println("PERSON:", person)

		req := &pb.GetBirthdayRequest{Person: person}
		result, err := client.UpdateBirthdayByIdAndName(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)

	})

	r.DELETE("/birthday/:userId", func(c *gin.Context) {
		userId := c.Param("userId")

		req := &pb.GetByIDRequest{UserId: userId}
		result, err := client.DeleteBirthdayByID(c, req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	r.GET("/birthday/:userId", func(c *gin.Context) {

		userId := c.Param("userId")
		req := &pb.GetByIDRequest{UserId: userId}

		result, err := client.GetBirthdayPersonByID(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)

	})

	r.POST("/birthday", func(c *gin.Context) {
		fmt.Println("in post")
		person := &pb.Person{}
		err := c.Bind(&person)

		fmt.Println(person)

		req := &pb.GetBirthdayRequest{Person: person}
		fmt.Println(req)

		result, err := client.CreateBirthdayPersonBy(c, req)

		if err != nil {
			fmt.Println("err", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	r.Run(config.PORT)
}

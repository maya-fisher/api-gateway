package main 

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/maya-fisher/birthday-service/proto"
	"google.golang.org/grpc"
)

func PostHomePage(c *gin.Context) {
	// var req createPerson
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	fmt.Println("error!!!!")
	// }

	// arg := createPerson{
	// 	name: req.name,
	// 	age: req.age,
	// }

	// fmt.Println("name",req.name, req.age)

	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println((string(value)))

	c.JSON(200, gin.H{
		"message": string(value),
	})

}


const (
	address = "localhost:50053"
)

func main() {

	// Set up a connection to the server.

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewBirthdaysClient(conn)

	r := gin.Default()


	r.PUT("/:id/:birthday", func (c *gin.Context)  {
		id := c.Param("id")
		birthday, err := strconv.ParseInt(c.Param("birthday"), 10, 64)

		person := &pb.Person{
			Name: id,
			Birthday: birthday,
		}
		req := &pb.GetBirthdayRequest{Person: person}
		res, err := client.UpdateBirthdayByIdAndName(c, req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": res, //res should return any id for now 
		})

	})

	r.DELETE("/:id", func (c *gin.Context)  {
		id := c.Param("id")

		req := &pb.GetByIDRequest{Id: id}
		res, err := client.DeleteBirthdayByID(c, req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": res, //res should return nothing for now 
		})
	})


	r.GET("/:id", func (c *gin.Context)  {
		id := c.Param("id")

		// Contact the server and print out its response.

		req := &pb.GetByIDRequest{Id: id}
		res, err := client.GetBirthdayPersonByID(c, req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": res, //res should return nothing for now 
		})
	
	})

	r.POST("/:name/:month/:day/:year", func (c *gin.Context)  {
		// month := c.Param("month")
		// day := c.Param("day")
		// year := c.Param("year")

		// date := time.Date(year, time.Month(mon), day, 0, 0, 0, 0, time.UTC)

		name := c.Param("name")
		bi := c.Param("birthday")
		birthday, err := strconv.ParseInt(bi, 10, 64)
	
		// Contact the server and print out its response.
		if err != nil {
			fmt.Println(err)
		}

		person := &pb.Person{
			Name: name,
			Birthday: birthday,
		}

		req := &pb.GetBirthdayRequest{Person: person}
		fmt.Println(req)
		res, err := client.CreateBirthdayPersonBy(c, req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": res, //res should return any id for now 
		})
	})

	// Run http server

	r.Run(":5050")
}
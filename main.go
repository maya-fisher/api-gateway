package main 

import (
	"log"
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	pb "github.com/maya-fisher/birthday-service/proto"
	"google.golang.org/grpc"
)


const (
	address = "localhost:50054"
	port = ":6060"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewBirthdaysClient(conn)


	r := gin.Default()


	r.PUT("/birthday/:userId", func (c *gin.Context)  {
		userId := c.Param("userId")

		person := &pb.Person{
			UserId: userId,
		}
		err = c.Bind(&person)

		fmt.Println(person)

		req := &pb.GetBirthdayRequest{Person: person}
		result, err := client.UpdateBirthdayByIdAndName(c, req)
		fmt.Println(result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK,result)

	})

	r.DELETE("/birthday/:userId", func (c *gin.Context)  {
		userId := c.Param("userId")

		req := &pb.GetByIDRequest{UserId: userId}
		result, err := client.DeleteBirthdayByID(c, req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK,result)
	})

	r.GET("/birthday/:userId", func (c *gin.Context)  {
		userId := c.Param("userId")
		req := &pb.GetByIDRequest{UserId: userId}
		result, err := client.GetBirthdayPersonByID(c, req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK,result)
	
	})


	r.POST("/birthday", func (c *gin.Context)  {


		person := &pb.Person{}
		err := c.Bind(&person)

		req := &pb.GetBirthdayRequest{Person: person}
		result, err := client.CreateBirthdayPersonBy(c, req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK,result)
	})

	
	r.Run(port)
}
package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"xdp-blocker/structs"

	"github.com/dropbox/goebpf"
	"github.com/gin-gonic/gin"
)

func Block(blacklist goebpf.Map) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var IPData structs.IPBlockReq
		if err := c.BindJSON(&IPData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Data": "Please Provide Correct Data Format",
			})
			return
		}
		go BlockIPAddress(IPData.IPAddress, IPData.Subnet, blacklist)

		c.JSON(200, gin.H{
			"Message": "Blocked IP Address " + IPData.IPAddress,
		})
	}

	return gin.HandlerFunc(fn)

}

func BlockIPAddress(ip string, subnet string, bpfMap goebpf.Map) error {
	IPAddress := fmt.Sprintf("%s/%s", ip, subnet)
	x := rand.NewSource(time.Now().UnixNano())
	y := rand.New(x).Intn(100000)
	err := bpfMap.Insert(goebpf.CreateLPMtrieKey(IPAddress), y)
	if err != nil {
		fmt.Println("Error in inserting: ", err)
		return err
	}
	return nil
}

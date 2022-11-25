package controllers

import (
	"fmt"
	"net/http"
	"xdp-blocker/structs"

	"github.com/dropbox/goebpf"
	"github.com/gin-gonic/gin"
)

func UnBlock(blacklist goebpf.Map) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var IPData structs.IPBlockReq
		if err := c.BindJSON(&IPData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Data": "Please Provide Correct Data Format",
			})
			return
		}
		go UnBlockIPAddress(IPData.IPAddress, IPData.Subnet, blacklist)

		c.JSON(200, gin.H{
			"Message": "UnBlocked IP Address " + IPData.IPAddress,
		})
	}

	return gin.HandlerFunc(fn)

}

func UnBlockIPAddress(ip string, subnet string, bpfMap goebpf.Map) error {
	IPAddress := fmt.Sprintf("%s/%s", ip, subnet)
	err := bpfMap.Delete(goebpf.CreateLPMtrieKey(IPAddress))
	if err != nil {
		fmt.Println("Error in deleting: ", err)
		return err
	}
	return nil
}

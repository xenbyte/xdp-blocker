package structs

type IPBlockReq struct {
	IPAddress string `json:"ipAddress"`
	Subnet    string `json:"subnet"`
}

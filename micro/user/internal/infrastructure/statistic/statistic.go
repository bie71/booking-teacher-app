package statistic

import (
	"auth/internal/config"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
)

type Responses struct {
	TotalBookings int64   `json:"total_bookings"`
	TotalRevenue  float64 `json:"total_revenue"`
	TotalTeachers int64   `json:"total_teachers"`
}

type Client struct {
	Client     *resty.Client
	CfgBooking *config.Service
	CfgTeacher *config.Service
}

func NewStatistic(CfgBooking *config.Service, CfgTeacher *config.Service, restyClient *resty.Client) *Client {

	return &Client{Client: restyClient, CfgBooking: CfgBooking, CfgTeacher: CfgTeacher}
}

func (c *Client) GetTotalBookings(ctx *gin.Context) (Responses, error) {

	var res Responses
	token, _ := ctx.Get("token")

	_, err := c.Client.R().
		SetResult(&res).
		SetAuthToken(cast.ToString(token)).
		Get(fmt.Sprintf("%s/api/v1/total-bookings", c.CfgBooking.Host))
	if err != nil {
		return Responses{}, err
	}

	return res, nil
}

func (c *Client) GetTotalTeachers(ctx *gin.Context) (Responses, error) {

	var res Responses

	token, _ := ctx.Get("token")

	_, err := c.Client.R().
		SetResult(&res).
		SetAuthToken(cast.ToString(token)).
		Get(fmt.Sprintf("%s/api/v1/total-teachers", c.CfgTeacher.Host))
	if err != nil {
		return Responses{}, err
	}
	return res, nil
}

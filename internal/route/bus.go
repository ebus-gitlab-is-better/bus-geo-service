package route

import (
	"bus-geo-service/internal/biz"
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BusRouter struct {
	uc *biz.BusUseCase
}

func NewBusRouter(uc *biz.BusUseCase) *BusRouter {
	return &BusRouter{uc: uc}
}

func (uc *BusRouter) Register(router *gin.RouterGroup) {
	router.POST("/:id", uc.Update)
}

type BusDTO struct {
	DriverID string
	Battery  uint
	Lat      float32
	Lon      float32
}

// @Summary	Send bus data
// @Accept		json
// @Produce	json
// @Tags		chat
//
// @Param		id	path	int	true	"Bus ID"	Format(uint64)
//
// @Param		dto	body	route.BusDTO	true	"dto"
//
//	@Success	200
//
// @Failure	401
// @Failure	403
// @Failure	500
// @Failure	400
// @Failure	404
// @Router		/bus/{id} [post]
func (r *BusRouter) Update(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.Atoi(id)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "parse id error",
		})
		return
	}
	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}
	dto := BusDTO{}

	err = json.Unmarshal(body, &dto)
	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}

	err = r.uc.Update(context.TODO(), &biz.Bus{
		ID:       uint(idUint),
		DriverID: dto.DriverID,
		Battery:  dto.Battery,
		Lat:      dto.Lat,
		Lon:      dto.Lon,
	})

	if err != nil {
		c.AbortWithStatusJSON(400, &gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(200)
}

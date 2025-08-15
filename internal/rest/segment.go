package rest

import (
	"net/http"
	"segment_service/domain"
	"strconv"

	"github.com/labstack/echo/v4"
)

type SegmentService interface {
	GetAll() []domain.Segment
	Delete(segment string) error
	Create(segment string) error
	Distribute(segment string, percentage float64) error
}

type SegmentHandler struct {
	segmentService SegmentService
}

func NewSegmentHandler(e *echo.Echo, ss SegmentService) {
	sh := &SegmentHandler{
		segmentService: ss,
	}
	e.GET("/segments", sh.GetAll)
	e.POST("/segments/:segment", sh.Create)
	e.PUT("/segments/:segment", sh.Distribute)
	e.DELETE("/segments/:segment", sh.Delete)
}

// ListSegments godoc
//
//	@Summary		List segments
//	@Description	get all segments
//	@Tags			segments
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		string
//	@Failure		500	{object}	rest.ErrorResponse
//	@Router			/segments [get]
func (sh *SegmentHandler) GetAll(c echo.Context) error {
	res := sh.segmentService.GetAll()

	return c.JSON(http.StatusOK, res)
}

// CreateSegment godoc
//
//	@Summary		Create segment
//	@Description	create a new segment
//	@Tags			segments
//	@Accept			json
//	@Produce		json
//	@Param			segment	path		string	true	"segment name"
//	@Success		200		{object}	rest.MessageResponse
//	@Failure		404		{object}	rest.ErrorResponse
//	@Failure		500		{object}	rest.ErrorResponse
//	@Router			/segments/{segment} [post]
func (sh *SegmentHandler) Create(c echo.Context) error {
	segment := c.Param("segment")

	err := sh.segmentService.Create(segment)
	if err != nil {
		return echo.NewHTTPError(getStatusCode(err), err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

// DistributeSegment godoc
//
//	@Summary		Distribute segment
//	@Description	distribute segment to a percentage of users
//	@Tags			segments
//	@Accept			json
//	@Produce		json
//	@Param			segment		path		string	true	"segment name"
//	@Param			percentage	query		number	true	"percentage of users to distribute (0-100)"
//	@Success		200			{object}	rest.MessageResponse
//	@Failure		400			{object}	rest.ErrorResponse
//	@Failure		404			{object}	rest.ErrorResponse
//	@Failure		500			{object}	rest.ErrorResponse
//	@Router			/segments/{segment} [put]
func (sh *SegmentHandler) Distribute(c echo.Context) error {
	segment := c.Param("segment")

	percentageString := c.QueryParam("percentage")
	if len(percentageString) == 0 {
		percentageString = "100"
	}

	percentage, err := strconv.ParseFloat(percentageString, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect percentage")
	}

	err = sh.segmentService.Distribute(segment, percentage)
	if err != nil {
		return echo.NewHTTPError(getStatusCode(err), err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

// DeleteSegment godoc
//
//	@Summary		Delete segment
//	@Description	delete segment
//	@Tags			segments
//	@Accept			json
//	@Produce		json
//	@Param			segment	path		string	true	"segment name"
//	@Success		200		{object}	rest.MessageResponse
//	@Failure		404		{object}	rest.ErrorResponse
//	@Failure		500		{object}	rest.ErrorResponse
//	@Router			/segments/{segment} [delete]
func (sh *SegmentHandler) Delete(c echo.Context) error {
	segment := c.Param("segment")

	err := sh.segmentService.Delete(segment)
	if err != nil {
		return echo.NewHTTPError(getStatusCode(err), err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

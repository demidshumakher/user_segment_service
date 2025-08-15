package rest

import (
	"net/http"
	"segment_service/domain"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	GetAll() (map[domain.User][]domain.Segment, error)
	GetSegmentsById(id int) ([]domain.Segment, error)
	ClearUserSegments(id int) error
	AddSegment(id int, segment string) error
	DeleteSegment(id int, segment string) error
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(e *echo.Echo, us UserService) {
	uh := &UserHandler{
		userService: us,
	}
	e.GET("/users", uh.GetAll)
	e.GET("/users/:id", uh.GetSegments)
	e.DELETE("/users/:id", uh.Clear)
	e.PUT("/users/:id/:segment", uh.AddSegment)
	e.DELETE("/users/:id/:segment", uh.DeleteSegment)
}

// ListUsers godoc
//
//	@Summary		List users
//	@Description	get users and their segments
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string][]string
//	@Failure		400	{object}	rest.ErrorResponse
//	@Failure		404	{object}	rest.ErrorResponse
//	@Failure		500	{object}	rest.ErrorResponse
//	@Router			/users [get]
func (uh *UserHandler) GetAll(c echo.Context) error {
	res, err := uh.userService.GetAll()
	if err != nil {
		return echo.NewHTTPError(getStatusCode(err), err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

// ShowUser godoc
//
//	@Summary		show user
//	@Description	show user's segments
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"user id"
//	@Success		200	{array}		string
//	@Failure		400	{object}	rest.ErrorResponse
//	@Failure		404	{object}	rest.ErrorResponse
//	@Failure		500	{object}	rest.ErrorResponse
//	@Router			/users/{id} [get]
func (uh *UserHandler) GetSegments(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	res, err := uh.userService.GetSegmentsById(id)
	if err != nil {
		return echo.NewHTTPError(getStatusCode(err), err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

// ClearSegments godoc
//
//	@Summary		Clear segments
//	@Description	Clear user segments
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"user id"
//	@Success		200	{object}	rest.MessageResponse
//	@Failure		400	{object}	rest.ErrorResponse
//	@Failure		404	{object}	rest.ErrorResponse
//	@Failure		500	{object}	rest.ErrorResponse
//	@Router			/users/{id} [delete]
func (uh *UserHandler) Clear(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	err = uh.userService.ClearUserSegments(id)
	if err != nil {
		return echo.NewHTTPError(getStatusCode(err), err)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

// DeleteSegment godoc
//
//	@Summary		Delete segment
//	@Description	delete segment
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"user id"
//	@Param			segment	path		string	true	"segment's name"
//	@Success		200		{object}	rest.MessageResponse
//	@Failure		400		{object}	rest.ErrorResponse
//	@Failure		404		{object}	rest.ErrorResponse
//	@Failure		500		{object}	rest.ErrorResponse
//	@Router			/users/{id}/{segment} [delete]
func (uh *UserHandler) DeleteSegment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	segment := c.Param("segment")

	err = uh.userService.DeleteSegment(id, segment)
	if err != nil {
		return echo.NewHTTPError(getStatusCode(err), err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

// AddSegment godoc
//
//	@Summary		Add segment
//	@Description	add new segment to user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"user id"
//	@Param			segment	path		string	true	"segment name"
//	@Success		200		{object}	rest.MessageResponse
//	@Failure		400		{object}	rest.ErrorResponse
//	@Failure		404		{object}	rest.ErrorResponse
//	@Failure		500		{object}	rest.ErrorResponse
//	@Router			/users/{id}/{segment} [put]
func (uh *UserHandler) AddSegment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	segment := c.Param("segment")

	err = uh.userService.AddSegment(id, segment)
	if err != nil {
		return echo.NewHTTPError(getStatusCode(err), err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

func getStatusCode(err error) int {
	switch err {
	case domain.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

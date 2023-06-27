package api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/westradev/webbr/webbr"
)

type Server struct {
	db *webbr.Webbr
}

func NewServer(db *webbr.Webbr) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) HandlePostInsert(c echo.Context) error {
	var (
		collName = c.Param("collname")
	)

	var data webbr.M
	if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil {
		return err
	}

	id, err := s.db.Insert(collName, data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, webbr.M{"id": id})
}

func (s *Server) HandleGetQuery(c echo.Context) error {
	records, err := s.db.Find("users", webbr.Filter{})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, records)
}

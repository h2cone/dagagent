package server

import (
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "healthy",
	})
}

func Upload(c echo.Context) error {
	subDir := c.FormValue("subDir")
	filename := c.FormValue("filename")

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dir := filepath.Join(DagFolder, subDir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	if len(filename) == 0 {
		filename = file.Filename
	}
	fileloc := filepath.Join(dir, filename)
	dst, err := os.Create(fileloc)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"fileloc": fileloc,
	})
}

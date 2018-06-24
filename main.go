package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ojrac/opensimplex-go"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"os"
	"time"
)

const SIZE = 1024

func noise(n *opensimplex.Noise, x, y int) float64 {
	return n.Eval2(float64(x)/SIZE-0.5, float64(y)/SIZE-0.5)/2.0 + 0.5
}

func genMap() string {
	now := time.Now().UnixNano()
	n := opensimplex.NewWithSeed(now)
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{SIZE, SIZE}})
	for x := 0; x < SIZE; x++ {
		for y := 0; y < SIZE; y++ {
			h := uint8(noise(n, x, y) * 256)
			img.Set(x, y, color.RGBA{h, h, h, h})
		}
	}

	file, err := os.Create(fmt.Sprintf("static/maps/%d.png", now))
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	png.Encode(file, img)
	return file.Name()
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/", "static/index.html")
	e.Static("/dist", "static/dist")
	e.Static("/static", "static")

	e.GET("/maps/new", func(c echo.Context) error {
		return c.String(http.StatusOK, genMap())
	})

	e.Logger.Fatal(e.Start(":8080"))
}

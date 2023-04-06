package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"net/http"

	"github.com/kbinani/screenshot"
	"github.com/labstack/echo/v4"
)

var (
	Version = "DEV"
)

func main() {
	app := echo.New()
	app.HideBanner = true

	app.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusAccepted, `
		<html>
			<body>
				<img src="/feed" style="height: 100%; width: 100%" />
			</body>
		</html>
		`)
	})

	app.GET("/feed", func(c echo.Context) (err error) {
		c.Response().Header().Set(echo.HeaderContentType, "multipart/x-mixed-replace; boundary=frame")
		c.Response().Header().Set("Cache-Control", "no-cache")
		c.Response().WriteHeader(http.StatusOK)

		for {
			var img *image.RGBA
			buffer := new(bytes.Buffer)

			if img, err = screenshot.CaptureDisplay(0); err != nil {
				return
			}

			if err = png.Encode(buffer, img); err != nil {
				return
			}

			frame := buffer.Bytes()

			c.Response().Write([]byte(fmt.Sprintf(
				"--frame\r\nContent-Type: image/png\r\nContent-Size: %d\r\n\r\n%s\r\n",
				len(frame),
				frame,
			)))
			c.Response().Flush()
		}
	})

	ip := flag.String("ip", "127.0.0.1", "IP address to stream on")
	port := flag.Int("port", 8080, "Port to stream through")

	flag.Parse()
	fmt.Printf("live-feed version %s\n", Version)
	app.Start(fmt.Sprintf("%s:%d", *ip, *port))
}

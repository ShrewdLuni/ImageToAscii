package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/nfnt/resize"
	"github.com/rs/cors"
)

var brightness string = ".'^,:;Il!i><~+_-?][}{1)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

type PixelColor struct {
	R int
	G int
	B int
}

type ColorfulAscii struct {
	Ascii string
	Color PixelColor
}

type Image struct {
	Image string  `json:"image"`
	Color [][]int `json:"color"`
}

func main() {
	mux := http.NewServeMux()
	handler := cors.Default().Handler(mux)

	mux.HandleFunc("/api", func(res http.ResponseWriter, req *http.Request) {

		file, _, err := req.FormFile("image")
		if err != nil {
			fmt.Print(err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Print("ok")
		data, err := json.Marshal(ProcessImage(file))
		if err != nil {
			fmt.Print(err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.Write(data)
	})

	if err := http.ListenAndServe("localhost:8080", handler); err != nil {
		fmt.Println(err.Error())
	}
}

func ProcessImage(file multipart.File) []ColorfulAscii {
	// imgRes, err := http.Get("https://c0.klipartz.com/pngpicture/562/67/gratis-png-tablero-de-geometria-tablero-de-instrumentos-hasta-puff-hexagonal-de-2-caras-geometria.png")
	// if err != nil || imgRes.StatusCode != 200 {
	// 	fmt.Println("Something went wrong")
	// }
	// defer imgRes.Body.Close()

	// img, err := png.Decode(imgRes.Body)

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	ratio := float64(height) / float64(width)
	width = 60
	height = int(float64(float64(width)/1.5) * float64(ratio))
	img = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	result := make([]ColorfulAscii, width*height)
	fmt.Print(width, height, ratio)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.At(x, y)
			color := color.RGBAModel.Convert(pixel).(color.RGBA)
			r := float64(color.R)
			g := float64(color.G)
			b := float64(color.B)
			sum := (((r + g + b) / 3) / 255) * 65

			item := ColorfulAscii{Ascii: string(brightness[int(math.Round(sum))]), Color: PixelColor{R: int(r), G: int(g), B: int(b)}}
			if item.Ascii == "." {
				item.Color.R = 255
				item.Color.G = 255
				item.Color.B = 255
			}
			result = append(result, item)
		}
		result = append(result, ColorfulAscii{Ascii: string("enter"), Color: PixelColor{R: int(0), G: int(0), B: int(0)}})
	}
	return result
}

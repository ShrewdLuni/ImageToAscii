package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"image/png"
	"math"
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
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(ProcessImage())
		if err != nil {
			fmt.Print(err)
		}
		res.Write(data)
	})

	handler := cors.Default().Handler(mux)

	mux.HandleFunc("/api", func(res http.ResponseWriter, req *http.Request) {
		//res.Header().Set("Content-Type", "application/json")
		//image := &Image{Image: ProcessImageNoColor(), Color: [][]int{{255, 255, 255}, {1, 1, 1}}}
		//data, err := json.Marshal(ProcessImageNoColor())
		//if err != nil {
		//	fmt.Print(err)
		//}
		fmt.Fprint(res, ProcessImageNoColor())
	})

	if err := http.ListenAndServe("localhost:8000", handler); err != nil {
		fmt.Println(err.Error())
	}

}

func ProcessImageNoColor() string {

	imgRes, err := http.Get("https://c0.klipartz.com/pngpicture/562/67/gratis-png-tablero-de-geometria-tablero-de-instrumentos-hasta-puff-hexagonal-de-2-caras-geometria.png")
	if err != nil || imgRes.StatusCode != 200 {
		fmt.Println("Something went wrong")
	}
	defer imgRes.Body.Close()

	img, err := png.Decode(imgRes.Body)
	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	ratio := float64(height) / float64(width)
	width = 150
	height = int(float64(float64(width)/1.5) * float64(ratio) * float64(0.5))
	img = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	result := ""
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.At(x, y)
			color := color.RGBAModel.Convert(pixel).(color.RGBA)
			r := float64(color.R)
			g := float64(color.G)
			b := float64(color.B)
			sum := (((r + g + b) / 3) / 255) * 65
			result += string(brightness[int(math.Round(sum))])
		}
		result += string("\n")
	}
	return result
}

func ProcessImage() []ColorfulAscii {

	imgRes, err := http.Get("https://c0.klipartz.com/pngpicture/562/67/gratis-png-tablero-de-geometria-tablero-de-instrumentos-hasta-puff-hexagonal-de-2-caras-geometria.png")
	if err != nil || imgRes.StatusCode != 200 {
		fmt.Println("Something went wrong")
	}
	defer imgRes.Body.Close()

	img, err := png.Decode(imgRes.Body)
	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	ratio := float64(height) / float64(width)
	width = 100
	height = int(float64(float64(width)/1.5) * float64(ratio) * float64(0.5))
	img = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	result := make([]ColorfulAscii, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.At(x, y)
			color := color.RGBAModel.Convert(pixel).(color.RGBA)
			r := float64(color.R)
			g := float64(color.G)
			b := float64(color.B)
			sum := (((r + g + b) / 3) / 255) * 65
			item := ColorfulAscii{Ascii: string(brightness[int(math.Round(sum))]), Color: PixelColor{R: int(r), G: int(g), B: int(b)}}
			result = append(result, item)
		}
		result = append(result, ColorfulAscii{Ascii: string("enter"), Color: PixelColor{R: int(0), G: int(0), B: int(0)}})
	}
	return result
}

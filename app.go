package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/nfnt/resize"
	"github.com/rs/cors"
)

var ASCIIbyBrightness string = ".'^,:;Il!i><~+_-?][}{1)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

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

	mux.HandleFunc("/test", func(res http.ResponseWriter, req *http.Request) {

		file, _, err := req.FormFile("image")
		if err != nil {
			fmt.Print(err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		resolution, err := strconv.Atoi(req.FormValue("resolution"))
		if err != nil {
			fmt.Print(err)
		}
		brightness, err := strconv.Atoi(req.FormValue("brightness"))
		if err != nil {
			fmt.Print(err)
		}

		data, err := json.Marshal(ProcessImage(file, resolution, float64(brightness)/100))
		if err != nil {
			fmt.Print(err)
		}
		res.Write(data)
	})

	if err := http.ListenAndServe("localhost:3001", handler); err != nil {
		fmt.Println(err.Error())
	}

}

func ProcessImageNoColor(file multipart.File) {
	img, err := png.Decode(file)
	fmt.Print(err)
	fmt.Print(img)
}

func ProcessImage(file multipart.File, resolution int, brightness float64) []ColorfulAscii {

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	ratio := float64(height) / float64(width)
	width = resolution
	height = int(float64(float64(width)/1.5) * float64(ratio))
	img = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	result := make([]ColorfulAscii, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.At(x, y)
			color := color.RGBAModel.Convert(pixel).(color.RGBA)
			r := float64(color.R)
			g := float64(color.G)
			b := float64(color.B)
			pixelBrightness := ((r + g + b) / 3)
			sum := (float64(Limit(int(pixelBrightness*brightness), 255, 1)) / 255) * 65

			item := ColorfulAscii{Ascii: string(ASCIIbyBrightness[int(math.Round(sum))]), Color: PixelColor{R: LimitPixel(r * brightness), G: LimitPixel(g * brightness), B: LimitPixel(b * brightness)}}
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

func Limit(value int, maxValue int, minValue int) int {
	if value >= maxValue {
		return maxValue
	}
	if value <= minValue {
		return minValue
	}
	return value
}

func LimitPixel(value float64) int {
	if value >= 255 {
		return 255
	}
	return int(math.Round(value))
}

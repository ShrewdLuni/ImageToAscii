package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
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

	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("There is nothing here,try /image")
	})

	mux.HandleFunc("/image", func(res http.ResponseWriter, req *http.Request) {

		resolution, err := strconv.Atoi(req.FormValue("resolution"))
		if err != nil {
			fmt.Printf("Resolution error: \"%s\"\n", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		brightness, err := strconv.ParseFloat(req.FormValue("brightness"), 64)
		if err != nil {
			fmt.Printf("Brightness error: \"%s\"\n", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		var img image.Image
		var dataError error

		image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
		image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)

		if req.FormValue("isFile") == "true" {
			file, _, err := req.FormFile("image")
			if err != nil {
				fmt.Printf("File error: \"%s\"\n", err)
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			img, _, dataError = image.Decode(file)
			if dataError != nil {
				fmt.Println("Error: Image could not be decoded")
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			imageResponse, err := http.Get(req.FormValue("link"))
			if err != nil || imageResponse.StatusCode == 201 {
				fmt.Println("Problem")
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			defer imageResponse.Body.Close()
			img, _, dataError = image.Decode(imageResponse.Body)
			if dataError != nil {
				fmt.Println("Error: Image could not be decoded")
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		data, err := json.Marshal(ProcessImage(img, resolution, brightness))
		if err != nil {
			fmt.Printf("Data error: \"%s\"\n", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.Write(data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if err := http.ListenAndServe("localhost:3001", handler); err != nil {
		fmt.Println(err.Error())
	}

}

func ProcessImage(img image.Image, resolution int, brightness float64) []ColorfulAscii {
	resolution = Limit(resolution, 1000, 1)
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
			sum := (float64(Limit(int(pixelBrightness*brightness), 255, 1)) / 255) * float64(len(ASCIIbyBrightness)-1)

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

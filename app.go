package main

import (
	"fmt"
	"image/color"
	"image/png"
	"math"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

var brightness string = " .'^,:;Il!i><~+_-?][}{1)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

type PixelColor struct {
	R int
	G int
	B int
}

type ColorfulAscii struct {
	Ascii string
	Color PixelColor
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, ProcessImage())
	})

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err.Error())
	}

}

func ProcessImage() []ColorfulAscii {

	imgRes, err := http.Get("https://upload.wikimedia.org/wikipedia/commons/thumb/9/9c/Nazi_Swastika.svg/langru-200px-Nazi_Swastika.svg.png")
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
	result := make([]ColorfulAscii, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.At(x, y)
			color := color.RGBAModel.Convert(pixel).(color.RGBA)
			r := float64(color.R)
			g := float64(color.G)
			b := float64(color.B)
			sum := (((r + g + b) / 3) / 255) * 66

			result = append(result, ColorfulAscii{
				Ascii: string(brightness[int(math.Round(sum))]),
				Color: PixelColor{
					R: int(r),
					G: int(g),
					B: int(b),
				},
			})
		}
		result = append(result, ColorfulAscii{
			Ascii: string("\n"),
			Color: PixelColor{
				R: int(0),
				G: int(0),
				B: int(0),
			},
		})
	}
	return result
}

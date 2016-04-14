package main

import (
	"crypto/sha1"
	"fmt"
	"image"
	"log"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
)

func generateHash(in string) string {
	h := sha1.New()
	h.Write([]byte(in))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func isEven(n int) bool {
	return n%2 == 0
}

func draw(m *gg.Context, pixels [7][7]bool, fg colorful.Color) {
	offset := m.Width() / 7
	var x float64
	var y float64
	var w float64
	var h float64
	bg, err := colorful.Hex("#F0F0F0")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i]); j++ {
			x = float64(i * offset)
			y = float64(j * offset)
			w = float64((i + 1) * offset)
			h = float64((j + 1) * offset)
			if pixels[i][j] {
				m.SetColor(fg)
				m.DrawRectangle(x, y, w, h)
				m.Fill()
			} else {
				m.SetColor(bg)
				m.DrawRectangle(x, y, w, h)
				m.Fill()
			}
		}
	}
}

// reflects the pixels over
func reflectPixels(pixels [7][7]bool) [7][7]bool {
	// Reflect over the middle line
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			pixels[6-i][j] = pixels[i][j]
		}
	}
	return pixels
}

func getPixels(hash string, flag bool) [7][7]bool {
	pixels := [7][7]bool{}
	for i := 1; i < 6; i++ {
		for j := 1; j < 6; j++ {
			char := hash[i*5+j]
			num, _ := strconv.ParseInt("0x"+string(char), 0, 10)
			if flag {
				pixels[i][j] = isEven(int(num))
			} else {
				pixels[i][j] = !isEven(int(num))
			}
		}
	}
	return pixels
}

func generate(c *cli.Context) {
	txt := c.String("message")
	if txt == "" {
		log.Fatal("fatal: no string provided")
	}
	hash := generateHash(txt)
	colorCode := hash[:6]
	num, _ := strconv.ParseInt("0x"+string(hash[0]), 0, 10)
	flag := isEven(int(num))
	width, height := 420, 420
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	fg, err := colorful.Hex("#" + colorCode)
	if err != nil {
		log.Fatal(err)
	}
	ctx := gg.NewContextForImage(m)
	basePixels := getPixels(hash, flag)
	finalPixels := reflectPixels(basePixels)
	draw(ctx, finalPixels, fg)
	ctx.SavePNG("identicon.png")
}

func main() {
	app := cli.NewApp()
	app.Name = "Identicon"
	app.Usage = "Generate identicons from strings"
	app.Version = "0.0.2"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "message, m",
			Usage: "The message to turn into an identicon",
		},
	}
	app.Action = generate
	app.Run(os.Args)
}

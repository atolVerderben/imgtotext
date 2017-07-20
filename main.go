package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
	"sync"
)

func loadimg(path string) (image.Image, error) {
	fImg1, _ := os.Open(path)
	defer fImg1.Close()
	img, _, err := image.Decode(fImg1)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func loop() {
	fmt.Println("Enter a png filename,tileSize[,outfile] to convert\nType 'exit' to quit")
	reader := bufio.NewReader(os.Stdin)
	rawtext, _ := reader.ReadString('\n')
	text := strings.TrimSpace(rawtext)
	switch text {
	case "exit":
		os.Exit(0)
		break
	default:
		if strings.Contains(text, ",") {
			vals := strings.Split(text, ",")
			if len(vals) < 2 || len(vals) > 3 {
				fmt.Println("Incorrect Format must be type: png,int[,outfile]")
				break
			}
			width, err := strconv.Atoi(strings.TrimSpace(vals[1]))
			if err != nil {
				fmt.Println("Size Not an Integer")
				break
			}
			if len(vals) == 2 {
				convertToAStar(strings.TrimSpace(vals[0]), width, "")
				break
			}
			if len(vals) == 3 {
				convertToAStar(strings.TrimSpace(vals[0]), width, strings.TrimSpace(vals[2]))
				break
			}
			/*height, err := strconv.Atoi(vals[2])
			if err != nil {
				fmt.Println("Height Incorrect Format")
				break
			}*/

		}

		break
	}
}

func convertToAStar(filename string, tileSize int, outputLoc string) {
	if !strings.Contains(filename, ".png") {
		filename += ".png"
	}
	img, err := loadimg(filename)
	if err != nil {
		fmt.Printf("Error: %s could not be opened: %s\n", filename, err.Error())
		return
		//os.Exit(1)
	}
	width, height := img.Bounds().Size().X, img.Bounds().Size().Y

	pixels, err := getPixels(img)
	astarMap := make([]string, width) // Let's make it the max length it could possibly be
	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		return
		//os.Exit(1)
	}
	var wg sync.WaitGroup
	for i := 0; i < height; i += tileSize {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			astarString := ""
			for j := 0; j < width; j += tileSize { // Or tileHeight if not square
				if pixels[i][j].R == 255 && pixels[i][j].G == 255 && pixels[i][j].B == 255 {
					astarString += " "
				} else {
					astarString += "X"
				}
			}
			astarString += "\n"
			astarMap[i] = astarString

		}(i)

	}
	wg.Wait()
	if outputLoc == "" {
		outputLoc = "output/" + filename + "_" + strconv.Itoa(tileSize) + "_astar.txt"
	}
	f, err := os.Create(outputLoc)
	if err != nil {
		fmt.Printf("Error: %s on %s\n", outputLoc, err.Error())
	}
	defer f.Close()

	for _, astarString := range astarMap {
		_, err = f.WriteString(astarString)
	}

	fmt.Printf("%s Successfully Converted to %s\n", filename, outputLoc)
}

func main() {
	// You can register another format here
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	for {
		loop()
	}

}

// Get the bi-dimensional pixel array
func getPixels(img image.Image) ([][]Pixel, error) {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

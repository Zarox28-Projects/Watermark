package engine

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// rotateRGBA rotates an RGBA image by the given angle
func rotateRGBA(src *image.RGBA, angle float64) *image.RGBA {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	cos, sin := math.Cos(angle), math.Sin(angle)

	newW := int(math.Abs(float64(w)*cos)+math.Abs(float64(h)*sin)) + 1
	newH := int(math.Abs(float64(w)*sin)+math.Abs(float64(h)*cos)) + 1
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	cx, cy := float64(w)/2, float64(h)/2
	ncx, ncy := float64(newW)/2, float64(newH)/2

	for y := range newH {
		for x := range newW {
			dx, dy := float64(x)-ncx, float64(y)-ncy
			srcX := int(math.Round(cos*dx + sin*dy + cx))
			srcY := int(math.Round(-sin*dx + cos*dy + cy))
			if srcX >= 0 && srcX < w && srcY >= 0 && srcY < h {
				dst.Set(x, y, src.At(srcX, srcY))
			}
		}
	}
	return dst
}

func ProcessImage(inputPath string, text string, outputPath string) (bool, *string) {
	// Open the input image file
	file, err := os.Open(inputPath)
	if err != nil {
		errStr := err.Error()
		return false, &errStr
	}
	defer file.Close()

	// Decode the input image
	img, _, err := image.Decode(file)
	if err != nil {
		errStr := err.Error()
		return false, &errStr
	}

	// Convert the image to RGBA
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	// Parse the font and create a face
	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		errStr := err.Error()
		return false, &errStr
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{Size: 40, DPI: 72})
	if err != nil {
		errStr := err.Error()
		return false, &errStr
	}
	defer face.Close()

	// Measure the text size and create a temporary image for the stamp
	advance := font.MeasureString(face, text)
	textW := advance.Ceil() + 10
	textH := 60

	// Draw the text onto the temporary image
	tmp := image.NewRGBA(image.Rect(0, 0, textW, textH))
	d := &font.Drawer{
		Dst:  tmp,
		Src:  image.NewUniform(color.White),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(5), Y: fixed.I(textH - 15)},
	}
	d.DrawString(text)

	// Apply transparency to the stamp image
	stamp := image.NewRGBA(image.Rect(0, 0, textW, textH))
	for y := range textH {
		for x := range textW {
			r, _, _, _ := tmp.At(x, y).RGBA()
			a := uint8((r >> 8) * 80 / 255)
			stamp.SetRGBA(x, y, color.RGBA{255, 255, 255, a})
		}
	}

	// Rotate the stamp image and get its bounds
	rotated := rotateRGBA(stamp, -math.Pi/6)
	rw, rh := rotated.Bounds().Dx(), rotated.Bounds().Dy()

	bounds := rgba.Bounds()
	W, H := bounds.Dx(), bounds.Dy()
	spacingX := rw + 80
	spacingY := rh + 60

	// Draw the rotated stamp onto the output image
	for y := -rh; y < H+rh; y += spacingY {
		for x := -rw; x < W+rw; x += spacingX {
			row := y / spacingY
			offsetX := x + (row%2)*spacingX/2
			dst := image.Pt(offsetX, y)
			draw.Draw(rgba, image.Rect(dst.X, dst.Y, dst.X+rw, dst.Y+rh), rotated, image.Point{}, draw.Over)
		}
	}

	// Save the output image
	out, err := os.Create(outputPath)
	if err != nil {
		errStr := err.Error()
		return false, &errStr
	}
	defer out.Close()

	// Encode the output image based on the file extension
	ext := strings.ToLower(filepath.Ext(outputPath))
	switch ext {
	case ".jpg", ".jpeg":
		if err := jpeg.Encode(out, rgba, &jpeg.Options{Quality: 95}); err != nil {
			errStr := err.Error()
			return false, &errStr
		}
	case ".png":
		if err := png.Encode(out, rgba); err != nil {
			errStr := err.Error()
			return false, &errStr
		}
	default:
		if err := png.Encode(out, rgba); err != nil {
			errStr := err.Error()
			return false, &errStr
		}
	}

	return true, nil
}

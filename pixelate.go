package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
)

func Pixelate(output string, faceInfo FaceInfo, multi bool, cnt int) error {
	file, err := os.Open(output)
	if err != nil {
		return err
	}
	defer file.Close()

	originimg, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	file, err = os.Open(output)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	dest := image.NewRGBA(bounds)

	min_y := faceInfo.FaceRectangle.Top - 5
	min_x := faceInfo.FaceRectangle.Left - 5

	h := (faceInfo.FaceRectangle.Height/10)*10 + 10
	w := (faceInfo.FaceRectangle.Width/10)*10 + 10
	block := w / 10

	max_y := min_y + h
	max_x := min_x + w

	for y := min_y; y < max_y; y = y + block {
		for x := min_x; x < max_x; x = x + block {
			var cr, cg, cb float32
			var alpha uint8
			for j := y; j < y+block; j++ {
				for i := x; i < x+block; i++ {
					if i >= 0 && j >= 0 && i < max_x && j < max_y {
						c := color.RGBAModel.Convert(img.At(i, j))
						col := c.(color.RGBA)
						cr += float32(col.R)
						cg += float32(col.G)
						cb += float32(col.B)
						alpha = col.A
					}
				}
			}
			cr = cr / float32(block*block)
			cg = cg / float32(block*block)
			cb = cb / float32(block*block)
			for j := y; j < y+block; j++ {
				for i := x; i < x+block; i++ {
					if i >= 0 && j >= 0 && i < max_x && j < max_y {
						dest.Set(i, j, color.RGBA{uint8(cr), uint8(cg), uint8(cb), alpha})
					}
				}
			}
		}
	}

	startPoint := image.Point{0, 0}

	pixelatedRect := image.Rectangle{startPoint, startPoint.Add(dest.Bounds().Size())}
	originRect := image.Rectangle{image.Point{0, 0}, originimg.Bounds().Size()}

	rgba := image.NewRGBA(originRect)
	draw.Draw(rgba, originRect, originimg, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, pixelatedRect, dest, image.Point{0, 0}, draw.Over)

	dstFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	jpeg.Encode(dstFile, rgba, nil)

	return nil
}

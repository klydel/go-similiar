// Copyright 2013 Adam Peck

package img

import(
	"util"
	"errors"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

func NNAvgLSHash(m image.Image) util.BitSet64 {
	b := m.Bounds()
	r := image.Rect(0, 0, 8, 8)
	i := image.NewGray(r)
	dx, dy := b.Dx() / r.Dx(), b.Dx() / r.Dy()
  for y := r.Min.Y; y < r.Max.Y; y++ {
    for x := r.Min.X; x < r.Max.X; x++ {
      i.Set(x, y, m.At(x * dx, y * dy))
    }
  }

	var sum uint
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
      sum += uint(i.At(x, y).(color.Gray).Y)
    }
  }
	avg := uint8(sum / 64)

	var h util.BitSet64
	d := r.Dy()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if i.At(x, y).(color.Gray).Y > avg {
				h.Set(uint(y * d + x))
			}
		}
	}
	return h
}

type Img struct {
	NNAvgLSHash64 util.BitSet64
}

func Fetch(url string, c *http.Client) (Img, error) {
	r, err := c.Get(url)
	if err != nil {
		return Img{}, err
	}
	if r.StatusCode != http.StatusOK {
		return Img{}, errors.New(http.StatusText(r.StatusCode))
	}
	defer r.Body.Close()

	i, _, err := image.Decode(r.Body)
	if err != nil {
		return Img{}, err
	}
	return Img{NNAvgLSHash(i)}, nil
}

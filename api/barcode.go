package api

import (
	"errors"
	"github.com/hoisie/web"
	"github.com/saintfish/barcode"
	"github.com/saintfish/hoc5lib/model"
	"github.com/saintfish/webutil"
	"image/png"
	"net/http"
	"regexp"
	"strconv"
)

func HandleEan13(ctx *web.Context, code string) {
	ean13, err := barcode.NewEan13(code)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	if ean13.String() != code {
		ctx.Redirect(http.StatusMovedPermanently, ean13.String())
		return
	}
	img := ean13.Encode()
	ctx.ContentType("image/png")
	err = png.Encode(ctx, img)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	return
}

var barcodePattern = regexp.MustCompile("^[0-9]{12}$")

func barcodeRange(start string, count int32, filter func(string) bool) ([]string, error) {
	if len(start) != 12 || !barcodePattern.MatchString(start) {
		return nil, errors.New("Invalid start barcode. Should be 12-digit number.")
	}
	if count < 0 {
		return nil, errors.New("Negative count.")
	}
	if count > 1000 {
		return nil, errors.New("Count too big.")
	}
	if count == 0 {
		return []string{}, nil
	}
	result := []string{}
	curr := start
	for i := 0; i < int(count); i++ {
		var code barcode.Barcode
		var err error
		for {
			code, err = barcode.NewEan13(curr)
			if err != nil {
				return nil, err
			}
			if !filter(code.String()) {
				break
			}
			curr = increment(curr)
			if len(curr) >= 13 {
				return result, nil
			}
		}
		result = append(result, code.String())
		curr = increment(curr)
	}
	return result, nil
}

func increment(s string) string {
	b := []byte(s)
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] < '0' && b[i] > '9' {
			panic("Invalid input")
		}
		if b[i] < '9' {
			b[i]++
			return string(b)
		} else {
			b[i] = '0'
		}
	}
	return "1" + string(b)
}

func HandleBookBarcodeRange(ctx *web.Context, start string, count string) {
	c, err := strconv.ParseInt(count, 10, 32)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	r, err := barcodeRange(start, int32(c), func(b string) bool {
		book, err := model.GetBook(b)
		return err == nil && book != nil
	})
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	webutil.Json(ctx, r)
	return
}

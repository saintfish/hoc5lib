package api

import (
	"bitbucket.org/saintfish/gopdf/pdf"
	"errors"
	"fmt"
	"github.com/hoisie/web"
	"github.com/saintfish/barcode"
	"github.com/saintfish/hoc5lib/model"
	"github.com/saintfish/webutil"
	"image/png"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

func barcodeRange(start string, count, step int32, filter func(string) bool) ([]string, error) {
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
			curr = increment(curr, step)
			if len(curr) >= 13 {
				return result, nil
			}
		}
		result = append(result, code.String())
		curr = increment(curr, step)
	}
	return result, nil
}

func increment(s string, step int32) string {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic("Invalid input")
	}
	i += int64(step)
	return fmt.Sprintf("%012d", i)
}

func HandleBookBarcodeRange(ctx *web.Context, start string, count string, step string) {
	c, err := strconv.ParseInt(count, 10, 32)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	s, err := strconv.ParseInt(step, 10, 32)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
	r, err := barcodeRange(start, int32(c), int32(s), func(b string) bool {
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

func HandleBookBarcodePDF(ctx *web.Context) {
	const kNumBarcodePerPage = 50
	r := strings.Split(ctx.Params["barcodes"], ",")
	for _, code := range r {
		_, err := barcode.NewEan13(code)
		if err != nil {
			webutil.Error(ctx, err)
			return
		}
	}
	doc := pdf.New()
	var canvas *pdf.Canvas
	for i := range r {
		if i%kNumBarcodePerPage == 0 {
			if canvas != nil {
				canvas.Close()
			}
			canvas = doc.NewPage(pdf.USLetterWidth, pdf.USLetterHeight)
			canvas.Translate(0.5*pdf.Inch, 0.5*pdf.Inch)
			for i := 0; i <= 5; i++ {
				canvas.DrawLine(
					pdf.Point{1.5 * pdf.Unit(i) * pdf.Inch, 0},
					pdf.Point{1.5 * pdf.Unit(i) * pdf.Inch, 10 * pdf.Inch})
			}
			for i := 0; i <= 10; i++ {
				canvas.DrawLine(
					pdf.Point{0, pdf.Unit(i) * pdf.Inch},
					pdf.Point{7.5 * pdf.Inch, pdf.Unit(i) * pdf.Inch})
			}
		}
		row := 9 - (i%kNumBarcodePerPage)/5
		col := (i % kNumBarcodePerPage) % 5
		code, _ := barcode.NewEan13(r[i])
		img := code.Encode()
		margin := 0.1 * pdf.Inch
		canvas.DrawImage(img, pdf.Rectangle{
			pdf.Point{1.5*pdf.Unit(col)*pdf.Inch + margin, pdf.Unit(row)*pdf.Inch + margin},
			pdf.Point{1.5*pdf.Unit(col+1)*pdf.Inch - margin, pdf.Unit(row+1)*pdf.Inch - margin},
		})
	}
	if canvas != nil {
		canvas.Close()
	}
	ctx.ContentType("application/pdf")
	err := doc.Encode(ctx)
	if err != nil {
		webutil.Error(ctx, err)
		return
	}
}

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"gitlab.com/xfmoulet/qoi"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("usage: " + os.Args[0] + " directory")
		return
	}
	dir := os.Args[1]
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for _, infile := range files {
		if !strings.HasSuffix(infile.Name(), ".png") {
			continue
		}
		f, err := os.Open(dir + "/" + infile.Name())
		if err != nil {
			fmt.Println("Error opening file: ", err)
			return
		}

		img, _, err := image.Decode(f)
		if err != nil {
			fmt.Println("Error decoding file: ", err)
			return
		}

		var of bytes.Buffer

		// QOI
		start := time.Now()
		qoi.Encode(&of, img)
		enc_qoi_duration := time.Since(start)

		start = time.Now()
		image.Decode(&of)
		dec_qoi_duration := time.Since(start)

		of.Reset()

		// PNG
		start = time.Now()
		png.Encode(&of, img)
		enc_png_duration := time.Since(start)

		start = time.Now()
		image.Decode(bytes.NewBuffer(of.Bytes()))
		dec_png_duration := time.Since(start)

		fmt.Printf("Encoding: QOI %4dms - PNG %4dms - Decoding: QOI %4dms PNG %4dms - %s\n",
			enc_qoi_duration.Milliseconds(),
			enc_png_duration.Milliseconds(),
			dec_qoi_duration.Milliseconds(),
			dec_png_duration.Milliseconds(),
			infile.Name())

	}
}

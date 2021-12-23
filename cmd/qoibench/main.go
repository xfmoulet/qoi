package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xfmoulet/qoi"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("usage: " + os.Args[0] + " directory")
		return
	}
	dir := os.Args[1]

	filepath.WalkDir(dir, func(path string, infile fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if infile.IsDir() {
			return nil
		}

		if !strings.HasSuffix(infile.Name(), ".png") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return err
		}

		img, _, err := image.Decode(f)
		if err != nil {
			fmt.Println("Error decoding file:", err)
			return err
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
			path)

		return nil
	})
}

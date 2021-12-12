package qoi

/*

QOI - The “Quite OK Image” format for fast, lossless image compression

Original version by Dominic Szablewski - https://phoboslab.org
Go version by Xavier-Frédéric Moulet

*/

import (
	"bufio"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"
)

const (
	qoi_INDEX   byte = 0b00_000000
	qoi_RUN_8   byte = 0b010_00000
	qoi_RUN_16  byte = 0b011_00000
	qoi_DIFF_8  byte = 0b10_000000
	qoi_DIFF_16 byte = 0b110_00000
	qoi_DIFF_24 byte = 0b1110_0000
	qoi_COLOR   byte = 0b1111_0000

	qoi_MASK_2 byte = 0b11_000000
	qoi_MASK_3 byte = 0b111_00000
	qoi_MASK_4 byte = 0b1111_0000
)

const qoiMagic = "qoif"

func qoi_COLOR_HASH(r, g, b, a byte) byte {
	return byte(r ^ g ^ b ^ a)
}

type pixel [4]byte

func Decode(r io.Reader) (image.Image, error) {
	cfg, err := DecodeConfig(r)
	if err != nil {
		return nil, err
	}

	b := bufio.NewReader(r)

	img := image.NewNRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))

	var index [64]pixel

	run := 0

	pixels := img.Pix // pixels yet to write
	px := pixel{0, 0, 0, 255}
	for len(pixels) > 0 {
		if run > 0 {
			run--
		} else {

			b1, err := b.ReadByte()
			// read and handle end of file
			if err == io.EOF {
				return img, nil
			}
			if err != nil {
				return nil, err
			}
			switch {
			case (b1 & qoi_MASK_2) == qoi_INDEX:
				px = index[b1^qoi_INDEX]

			case (b1 & qoi_MASK_3) == qoi_RUN_8:
				run = int(b1 & 0x1f)

			case ((b1 & qoi_MASK_3) == qoi_RUN_16):
				b2, err := b.ReadByte()
				if err != nil {
					return nil, err
				}
				run = (((int(b1) & 0x1f) << 8) | int(b2)) + 32

			case ((b1 & qoi_MASK_2) == qoi_DIFF_8):
				px[0] += ((b1 >> 4) & 0x03) - 2
				px[1] += ((b1 >> 2) & 0x03) - 2
				px[2] += (b1 & 0x03) - 2

			case ((b1 & qoi_MASK_3) == qoi_DIFF_16):
				b2, err := b.ReadByte()
				if err != nil {
					return nil, err
				}
				px[0] += (b1 & 0x1f) - 16
				px[1] += (b2 >> 4) - 8
				px[2] += (b2 & 0x0f) - 8

			case ((b1 & qoi_MASK_4) == qoi_DIFF_24):
				b2, err := b.ReadByte()
				if err != nil {
					return nil, err
				}
				b3, err := b.ReadByte()
				if err != nil {
					return nil, err
				}

				px[0] += (((b1 & 0x0f) << 1) | (b2 >> 7)) - 16
				px[1] += ((b2 & 0x7c) >> 2) - 16
				px[2] += (((b2 & 0x03) << 3) | ((b3 & 0xe0) >> 5)) - 16
				px[3] += (b3 & 0x1f) - 16

			case (b1 & qoi_MASK_4) == qoi_COLOR:
				if b1&8 != 0 {
					b2, err := b.ReadByte()
					if err != nil {
						return nil, err
					}
					px[0] = b2
				}
				if b1&4 != 0 {
					b2, err := b.ReadByte()
					if err != nil {
						return nil, err
					}
					px[1] = b2
				}
				if b1&2 != 0 {
					b2, err := b.ReadByte()
					if err != nil {
						return nil, err
					}
					px[2] = b2
				}
				if b1&1 != 0 {
					b2, err := b.ReadByte()
					if err != nil {
						return nil, err
					}
					px[3] = b2
				}
			default:
				px = pixel{255, 0, 255, 255}
			}

			index[int(qoi_COLOR_HASH(px[0], px[1], px[2], px[3]))%len(index)] = px
		}

		// TODO stride ..
		copy(pixels[:4], px[:])
		pixels = pixels[4:] // advance
	}
	return img, nil
}

func Encode(w io.Writer, m image.Image) error {
	minX := m.Bounds().Min.X
	maxX := m.Bounds().Max.X
	minY := m.Bounds().Min.Y
	maxY := m.Bounds().Max.Y

	var out = bufio.NewWriter(w)

	// write header to output
	if err := binary.Write(out, binary.BigEndian, []byte(qoiMagic)); err != nil {
		return err
	}
	// width
	if err := binary.Write(out, binary.BigEndian, uint32(maxX-minX)); err != nil {
		return err
	}
	// height
	if err := binary.Write(out, binary.BigEndian, uint32(maxY-minY)); err != nil {
		return err
	}
	// channels
	if err := binary.Write(out, binary.BigEndian, uint8(4)); err != nil {
		return err
	}
	// 0b0000rgba colorspace
	if err := binary.Write(out, binary.BigEndian, uint8(0)); err != nil {
		return err
	}

	var index [64]pixel
	px_prev := pixel{0, 0, 0, 255}
	run := 0

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			// extract pixel and convert to non-premultiplied
			c := color.NRGBAModel.Convert(m.At(x, y))
			c_r, c_g, c_b, c_a := c.RGBA()
			px := pixel{byte(c_r >> 8), byte(c_g >> 8), byte(c_b >> 8), byte(c_a >> 8)}

			if px == px_prev {
				run++
			}

			last_pixel := x == (maxX-1) && y == (maxY-1)
			if run > 0 && (run == 0x2020 || px != px_prev || last_pixel) {
				if run < 33 {
					out.WriteByte(qoi_RUN_8 | byte(run-1))
				} else {
					run -= 33
					out.WriteByte(qoi_RUN_16 | byte(run>>8))
					out.WriteByte(byte(run & 0xff))
				}
				run = 0
			}

			if px != px_prev {
				var index_pos byte = qoi_COLOR_HASH(px[0], px[1], px[2], px[3]) % 64
				if index[index_pos] == px {
					out.WriteByte(qoi_INDEX | index_pos)
				} else {
					index[index_pos] = px

					vr := int(px[0]) - int(px_prev[0])
					vg := int(px[1]) - int(px_prev[1])
					vb := int(px[2]) - int(px_prev[2])
					va := int(px[3]) - int(px_prev[3])

					if vr > -17 && vr < 16 &&
						vg > -17 && vg < 16 &&
						vb > -17 && vb < 16 &&
						va > -17 && va < 16 {
						switch {
						case va == 0 &&
							vr > -3 && vr < 2 &&
							vg > -3 && vg < 2 &&
							vb > -3 && vb < 2:
							out.WriteByte(qoi_DIFF_8 | byte(((vr+2)<<4)|(vg+2)<<2|(vb+2)))
						case va == 0 &&
							vr > -17 && vr < 16 &&
							vg > -9 && vg < 8 &&
							vb > -9 && vb < 8:
							out.WriteByte(qoi_DIFF_16 | byte(vr+16))
							out.WriteByte(byte(((vg + 8) << 4) | (vb + 8)))
						default:
							out.WriteByte(qoi_DIFF_24 | byte((vr+16)>>1))
							out.WriteByte(byte(((vr + 16) << 7) | ((vg + 16) << 2) | ((vb + 16) >> 3)))
							out.WriteByte(byte(((vb + 16) << 5) | (va + 16)))
						}
					} else {
						mask := qoi_COLOR
						if vr != 0 {
							mask |= 1 << 3
						}
						if vg != 0 {
							mask |= 1 << 2
						}
						if vb != 0 {
							mask |= 1 << 1
						}
						if va != 0 {
							mask |= 1 << 0
						}
						out.WriteByte(mask)
						if vr != 0 {
							out.WriteByte(px[0])
						}
						if vg != 0 {
							out.WriteByte(px[1])
						}
						if vb != 0 {
							out.WriteByte(px[2])
						}
						if va != 0 {
							out.WriteByte(px[3])
						}
					}

				}
			}

			px_prev = px
		}
	}
	binary.Write(out, binary.BigEndian, uint32(0)) // padding

	return out.Flush()
}

func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	var header [4 + 4 + 4 + 1 + 1]byte
	if _, err = io.ReadAtLeast(r, header[:], len(header)); err != nil {
		return
	}

	if string(header[:4]) != qoiMagic {
		return cfg, errors.New("Invalid magic")
	}

	return image.Config{
		Width:      int(binary.BigEndian.Uint32(header[4:])),
		Height:     int(binary.BigEndian.Uint32(header[8:])),
		ColorModel: color.NRGBAModel,
	}, err
}

func init() {
	image.RegisterFormat("qoi", qoiMagic, Decode, DecodeConfig)
}

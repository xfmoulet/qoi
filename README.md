# QOI - The “Quite OK Image” format for fast, lossless image compression

package and small utilities in pure Go, quite OK implementation

See [qoi.h](https://github.com/phoboslab/qoi/blob/master/qoi.h) for
the documentation.

More info at https://phoboslab.org/log/2021/11/qoi-fast-lossless-image-compression

## Performance

Performance is currently around half C version (optimized at `-O3`)

```
$ ./qoibench ../../../images/wallpaper/
Encoding: QOI  184ms - PNG 1382ms - Decoding: QOI   62ms PNG  178ms - 1492858.png
Encoding: QOI  180ms - PNG 1389ms - Decoding: QOI   51ms PNG  156ms - 1492868.png
Encoding: QOI  152ms - PNG  863ms - Decoding: QOI   30ms PNG   84ms - 1492869.png
Encoding: QOI  204ms - PNG 1147ms - Decoding: QOI   71ms PNG  176ms - 1492893.png
Encoding: QOI  316ms - PNG 2697ms - Decoding: QOI  103ms PNG  274ms - EwZDbLoWQAEskRq.png
Encoding: QOI  194ms - PNG 1524ms - Decoding: QOI   62ms PNG  156ms - Hy23XKX.png
Encoding: QOI   43ms - PNG  301ms - Decoding: QOI   15ms PNG   34ms - Screenshot_2021-11-16_13-57-47.png
Encoding: QOI  180ms - PNG 1078ms - Decoding: QOI   45ms PNG  134ms - car.png
``` 

```
$ ./qoibench ../../../images/kodak/
Encoding: QOI   19ms - PNG   69ms - Decoding: QOI    6ms PNG   17ms - kodim01.png
Encoding: QOI   15ms - PNG   96ms - Decoding: QOI    5ms PNG   14ms - kodim02.png
Encoding: QOI   15ms - PNG  143ms - Decoding: QOI    5ms PNG   13ms - kodim03.png
Encoding: QOI   16ms - PNG   87ms - Decoding: QOI    5ms PNG   15ms - kodim04.png
Encoding: QOI   17ms - PNG   68ms - Decoding: QOI    6ms PNG   15ms - kodim05.png
Encoding: QOI   15ms - PNG   87ms - Decoding: QOI    7ms PNG   13ms - kodim06.png
Encoding: QOI   21ms - PNG  111ms - Decoding: QOI    5ms PNG   16ms - kodim07.png
Encoding: QOI   18ms - PNG   67ms - Decoding: QOI    6ms PNG   19ms - kodim08.png
Encoding: QOI   16ms - PNG  100ms - Decoding: QOI    6ms PNG   19ms - kodim09.png
Encoding: QOI   16ms - PNG   93ms - Decoding: QOI    5ms PNG   15ms - kodim10.png
Encoding: QOI   16ms - PNG   93ms - Decoding: QOI    6ms PNG   16ms - kodim11.png
Encoding: QOI   16ms - PNG  110ms - Decoding: QOI    5ms PNG   15ms - kodim12.png
``` 

## Example Usage

- `cmd/qoiconv` converts between png <> qoi
- `cmd/qoibench` bench the en/decoding vs. golang png implementation

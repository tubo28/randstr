# randstr

Generates random string.

```shell
$ go install github.com/tubo28/randstr@latest
$ randstr --help
Usage of $GOPATH/bin/randstr:
  -alnum
    	Digits + Latin alphabet (default)
  -digit
    	Digits
  -graph
    	Digits + Latin alphabet + Symbols
  -pattern string
    	Pattern. Each X is replaced with a random character. (default "XXXXXXXXXXXXXXXX")
$ randstr
6qf8PNUykE5Mb4tH
$ randstr --pattern=XXXXXXXX --digit
27844573
$ randstr --pattern=XXXXXXXX --alnum
7mHiYRd4
$ randstr --pattern=XXXXXXXX --graph
tvK?fgy6
```

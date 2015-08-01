fuzzing
=======

First get go-fuzz:

    go get github.com/dvyukov/go-fuzz/go-fuzz
    go get github.com/dvyukov/go-fuzz/go-fuzz-build

Now build the zip:

    go-fuzz-build github.com/vrischmann/envconfig/fuzz

Now generate a base corpus:

    cd $GOPATH/github.com/vrischmann/envconfig/fuzz
    go run gen-base-corpus/main.go

Now run go-fuzz:

    go-fuzz -bin=./fuzz-fuzz.zip -workdir=.

And now we wait.

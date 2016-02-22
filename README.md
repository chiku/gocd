# gocd

[![Build Status](https://travis-ci.org/chiku/gocd.svg?branch=master)](https://travis-ci.org/chiku/gocd)

Golang client for [Gocd](https://www.go.cd)

Development prerequisites
-------------------------

* Install `make`
* [Install golang](https://golang.org/doc/install) (1.5 or better). Please set `GO15VENDOREXPERIMENT=1` if you are using golang 1.5
* Add `$GOPTAH/bin` to `PATH`
* Install `glide`
`go get github.com/Masterminds/glide`

Running tests
-------------

```shell
make prereqs
make all
```

Contributing
------------

* Fork the project.
* Make your feature addition or bug fix.
* Add tests for it. This is important so I don't break it in a future version unintentionally.
* Commit, but do not mess with the VERSION. If you want to have your own version, that is fine but bump the version in a commit by itself in another branch so I can ignore it when I pull.
* Send me a pull request.

License
-------

This library is released under the MIT license. Please refer to LICENSE for more details.

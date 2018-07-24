# Chrome WebDriver Client for Go

[![GoDoc](https://godoc.org/github.com/radutopala/webdriver?status.svg)](https://godoc.org/github.com/radutopala/webdriver)
[![Travis](https://travis-ci.org/radutopala/webdriver.svg?branch=master)](https://travis-ci.org/radutopala/webdriver)
[![Go Report Card](https://goreportcard.com/badge/github.com/radutopala/webdriver)](https://goreportcard.com/report/github.com/radutopala/webdriver)

## About

This is a [WebDriver] client for [Go][go], supporting the [WebDriver
protocol][webdriver] for Chrome and [ChromeDriver][chromedriver]. 

[webdriver]: https://www.w3.org/TR/webdriver/
[go]: http://golang.org/
[chromedriver]: https://sites.google.com/a/chromium.org/chromedriver/

## Installing

Run

    go get -u github.com/radutopala/webdriver

to fetch the package.

The package requires a working WebDriver installation, which can include recent versions of a web browser being driven by Selenium WebDriver.

## Documentation

The API documentation is at https://godoc.org/github.com/radutopala/webdriver. See [the unit
tests](https://github.com/radutopala/webdriver/blob/master/remote_test.go) for better usage information.

### Downloading and Pack Dependencies

Download and pack the ChromeDriver binaries:

    $ go run download/download.go

You only have to do this once initially and later when version numbers in download.go change.

### Testing Locally

Run the tests:

    $ go test 

* There is one top-level test for Chrome and ChromeDriver.
    
* There are subtests that are shared between both top-level tests.

* To run only the top-level tests, pass:

    * `-test.run=TestChrome`

* To run a specific subtest, pass `-test.run=TestChrome/<subtest>` as
  appropriate. This flag supports regular expressions.

* If the Chrome binaries or the ChromeDriver binary cannot be found, the corresponding tests will be
  skipped.

* The binaries under test can be configured by passing flags to `go
  test`. See the available flags with `go test --arg --help`.

* Add the argument `-test.v` to see detailed output from the test automation framework.

### Testing With Docker

To ensure hermeticity, we also have tests that run under Docker. You will need an installed and running Docker system.

To run the tests under Docker, run:

    $ go test --docker

This will create a new Docker container and run the tests in it. (Note: flags supplied to this invocation are not carried through to the `go test` invocation within the Docker container).

For debugging Docker directly, run the following commands:

    $ docker build -t webdriver testing/
    $ docker run --volume=${GOPATH?}:/code --workdir=/code/src/github.com/radutopala/webdriver -it webdriver bash

## License

This project is licensed under the [MIT][mit] license.

[mit]: https://raw.githubusercontent.com/radutopala/webdriver/master/LICENSE

Please note that this project is a cut-down version of [tebeka/selenium](https://github.com/tebeka/selenium), with ideas from [Symfony/Panther](https://github.com/symfony/panther), targeting only the Chrome WebDriver implementation.

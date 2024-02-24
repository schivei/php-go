# php-go

php-go allows to call Go code from PHP, with minimal code boilerplate.

## Goals:

- Allow to _export_ Go functions and Go constants from Go to PHP
- Be reliable and always safe
- Deploy Go code without re-building the PHP extension

## DOING:

- [X] Support GoLang Http Requests

## Install

You can download this package using "go install".
Then you can run:

```
go install github.com/schivei/php-go
```

When this is finished, go to your project directory and run:

```
php-go
```

You also can pass binary directory as argument:

```
php-go /path/to/bin
```

The extension will be built and placed in the current directory.

Then copy the resulting ``phpgo.so`` to your PHP extensions directory and add ``extension=phpgo.so`` to your php.ini.

Alternatively, you can use ``phpgo.so`` from the ``bin`` directory.

Note: php-go supports PHP 8+ (non-ZTS).

## Usage

#### Exporting Go functions

``` go
package main

import (
  "strings"
  "github.com/schivei/php-go/php"
)

// call php.Export() for its side effects
var _ = php.Export("example", map[string]interface{}{
  "toUpper": strings.ToUpper,
  "takeOverTheWorld": TakeOverTheWorld,
})

func TakeOverTheWorld() {
}

func main() {
}
```

The module can then be compiled as a shared library using `-buildmode c-shared`:

    go build -o example.so -buildmode c-shared .

Note: Go **requires** that the module be a _main_ package with a _main_ function in this mode.

#### Using the module in PHP

``` php
// Create a class from the Go module, and return an instance of it
$module = phpgo_load("/path/to/example.so", "example");

// Call some method
$module->toUpper("foo");
```

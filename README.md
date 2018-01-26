# DynamoDB Storage for OAuth 2.0

> Based on the https://github.com/go-oauth2/mongo

[![License][License-Image]][License-Url]

## Install

``` bash
$ go get -u github.com/contamobi/go-oauth2-dynamodb
```

## Usage

``` go
package main

import (
	"github.com/contamobi/go-oauth2-dynamodb"
	"github.com/contamobi/go-oauth2/manage"
)

func main() {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(
		mongo.NewTokenStore(mongo.NewConfig(
			"mongodb://127.0.0.1:27017",
			"oauth2",
		)),
	)
	// ...
}
```

## MIT License

```
Copyright (c) 2018 Conta.MOBI
```

[License-Url]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg
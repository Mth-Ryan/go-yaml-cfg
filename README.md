# Go Yaml Config
Configure your **Golang** project using a [Yaml](https://yaml.org/) file who accepts env variables on his strings just like a docker-compose file.

### Summary
* [Features](#features)
* [Example](#example)
* [Limitations](#limitations)

## Features
* **Full yaml parsing:** go-yaml-config uses [go-yaml](https://github.com/go-yaml/yaml) for parsing the config files
* **Env Variables support:** all the env variables patterns found on the strings of the config file will be replaced with the respective env or default value, when a default is not provided and the env is not setted, the value will be an empty string.
* **Type Safety:** the config struct type has to be defined on the code, and passed to the config initialization as a generic type argument. In the same way, it is retrieved on the same type as a generic argument. If the type is mismatched, you can handle the error or simply let the application panics with The Must Function.
* **Singleton Pattern**: The initialization guarantees that the config is only initialized once. And all the get operations are always in a constant time complexity.

## Example
First, create a config file for your application on any accessible path with any name, in this case will be on the app root: `./appSettings.yml`, with the follow content:

``` yaml
dbSettings:
  host: "localhost"
  port: 5432
  name: "someDB"
  user: "someUser"
  pass: "${PASS:-somePass}"
```

Then create the corresponding struct in your Go app and initialize the config:
``` go
package main

import (
    	"fmt"
    	conf "github.com/Mth-Ryan/go-yaml-cfg"
)

// Struct Fields must be public
type AppSettings struct {
    	DBSettings struct {
            	Host string `yaml:"host"`
            	Port int	`yaml:"port"`
            	Name string `yaml:"name"`
            	User string `yaml:"user"`
            	Pass string `yaml:"pass"`
    	} `yaml:"dbSettings"`
}

func main() {
    	var config AppSettings
    	var err error

    	if err = conf.InitializeConfigSingleton[AppSettings]("./appSettings.yml"); err != nil {
            	panic(err)
    	}

    	config, err = conf.GetConfig[AppSettings]()
    	if err != nil {
            	panic(err)
    	}
    	fmt.Printf("%+v \n", config)
}

```

If you build and run this project you will get:
``` bash
go build -o test
./test
{DBSettings:{Host:localhost Port:5432 Name:someDB User:someUser Pass:somePass}}
```

Now Overriding the default password with a env variable:
``` bash
PASS="otherPass" ./test
{DBSettings:{Host:localhost Port:5432 Name:someDB User:someUser Pass:otherPass}}
```

## Limitations
* In the current lib state it is not possible to replace env variables inside arrays whatsoever. but we are working on that
* All the fields must be public, this is a limitation of go's reflection lib, there is no way to change this.
* The lib requires go >= 1.18

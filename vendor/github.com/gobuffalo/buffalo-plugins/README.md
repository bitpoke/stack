# buffalo-plugins

A plugin for working with [Buffalo](https://gobuffalo.io) plugins.

## Installation

```bash
$ go get -u -v github.com/gobuffalo/buffalo-plugins
```

## Generating a New Buffalo Plugin

```bash
$ buffalo generate plugin --help
```

## Listing Currently Installed Plugins

```bash
$ buffalo plugins list
```

## Managing Plugins

`buffalo plugins install pkg pkg pkg...`
`buffalo plugins remove pkg pkg pkg...`

Plugins can be configured via the `./config/buffalo-plugins.toml` file that acts as the official list of plugins the application depends on. This is file is optional, until using the `install` and `remove` commands.

With no config file the output of `buffalo plugins list` and `plugins.List()` does not change. Once a configuration file is in place, then that file will dictate the output of those commands.

```bash
$ buffalo plugins list

Bin              |Command                    |Description
---              |---                        |---
buffalo-auth     |buffalo generate auth      |Generates a full auth implementation
buffalo-pop      |buffalo db                 |[DEPRECATED] please use `buffalo pop` instead.
buffalo-goth     |buffalo generate goth-auth |Generates a full auth implementation use Goth
buffalo-goth     |buffalo generate goth      |generates a actions/auth.go file configured to the specified providers.
buffalo-heroku   |buffalo heroku             |helps with heroku setup and deployment for buffalo applications
buffalo-pop      |buffalo destroy model      |Destroys model files.
buffalo-plugins  |buffalo generate plugin    |generates a new buffalo plugin
buffalo-plugins  |buffalo plugins            |tools for working with buffalo plugins
buffalo-pop      |buffalo pop                |A tasty treat for all your database needs
buffalo-trash    |buffalo trash              |destroys and recreates a buffalo app
buffalo-upgradex |buffalo upgradex           |updates Buffalo and/or Pop/Soda as well as your app
```

To add support for the plugin manager, one can either manually edit `./config/buffalo-plugins.toml` or let `buffalo plugins install` create it for you.

```bash
$ buffalo plugins install

go get github.com/gobuffalo/buffalo-pop
./config/buffalo-plugins.toml
```

``` bash
$ cat ./config/buffalo-plugins.toml

[[plugin]]
  binary = "buffalo-pop"
  go_get = "github.com/gobuffalo/buffalo-pop"
```

```bash
$ buffalo plugins list

Bin         |Command               |Description
---         |---                   |---
buffalo-pop |buffalo db            |[DEPRECATED] please use `buffalo pop` instead.
buffalo-pop |buffalo destroy model |Destroys model files.
buffalo-pop |buffalo pop           |A tasty treat for all your database needs
```

The `buffalo-pop` plugin was automatically added because the application in this example is a Buffalo application that uses Pop.

New plugins can be install in bulk with the `install` command

```bash
$ buffalo plugins install github.com/markbates/buffalo-trash github.com/gobuffalo/buffalo-heroku

go get github.com/gobuffalo/buffalo-heroku
go get github.com/gobuffalo/buffalo-pop
go get github.com/markbates/buffalo-trash
./config/buffalo-plugins.toml
```

```bash
$ buffalo plugins list

Bin            |Command               |Description
---            |---                   |---
buffalo-pop    |buffalo db            |[DEPRECATED] please use `buffalo pop` instead.
buffalo-heroku |buffalo heroku        |helps with heroku setup and deployment for buffalo applications
buffalo-pop    |buffalo destroy model |Destroys model files.
buffalo-pop    |buffalo pop           |A tasty treat for all your database needs
buffalo-trash  |buffalo trash         |destroys and recreates a buffalo app
```

``` bash
$ cat ./config/buffalo-plugins.toml

[[plugin]]
  binary = "buffalo-heroku"
  go_get = "github.com/gobuffalo/buffalo-heroku"

[[plugin]]
  binary = "buffalo-pop"
  go_get = "github.com/gobuffalo/buffalo-pop"

[[plugin]]
  binary = "buffalo-trash"
  go_get = "github.com/markbates/buffalo-trash"
```

Finally plugins can be removed with the `remove` command. This only removes them from the config file, not from the users system.

```bash
$ buffalo plugins remove github.com/gobuffalo/buffalo-heroku

./config/buffalo-plugins.toml
```

``` bash
$ cat ./config/buffalo-plugins.toml

[[plugin]]
  binary = "buffalo-pop"
  go_get = "github.com/gobuffalo/buffalo-pop"

[[plugin]]
  binary = "buffalo-trash"
  go_get = "github.com/markbates/buffalo-trash"
```

```bash
$ buffalo plugins list

Bin           |Command               |Description
---           |---                   |---
buffalo-pop   |buffalo db            |[DEPRECATED] please use `buffalo pop` instead.
buffalo-pop   |buffalo destroy model |Destroys model files.
buffalo-pop   |buffalo pop           |A tasty treat for all your database needs
buffalo-trash |buffalo trash         |destroys and recreates a buffalo app
```

## Listening for Plugin Setup Instructions

In Buffalo `v0.13.1-beta.1` events are now emitted with the `buffalo setup` command. The `buffalo-plugins` command will listen for this event and install the necessary plugins for an application. When completed it will emit the `plugins.EvtSetupFinished` event. This event should be listened to by other plugins to run their setup commands.

# Gracenote → FreeDB Proxy

## What is this?
This is a proxy that allows FreeDB aware software to access metadata from the Gracenote CD database.

The root `cddb` package can also be used as a rudimentary Gracenote library.

## Table of Contents
- [Gracenote → FreeDB Proxy](#gracenote--freedb-proxy)
  * [What is this?](#what-is-this)
  * [Configuration](#configuration)
    + [Command line switches](#command-line-switches)
    + [General Environment Variables](#general-environment-variables)
    + [App Engine Environment Variables](#app-engine-environment-variables)
    + [Compile-time configuration](#compile-time-configuration)
  * [Building + Running / Deploying](#building--running--deploying)
    + [Standalone](#standalone)
    + [Google App Engine](#google-app-engine)
  * [Usage](#usage)
    + [Advanced Usage](#advanced-usage)
      - [Parameters](#parameters)

## Configuration
The proxy requires a valid Gracenote Client and User ID.
You can get these at [https://developer.gracenote.com/web-api](https://developer.gracenote.com/web-api).

The proxy provides multiple ways to configure it's settings.
Configuration values are taken with the following priority: `Command line switch > Environment variable > Compile time option`

### Command line switches
| Switch   | Description                            |
|----------|----------------------------------------|
| -address | The address to listen on               |
| -port    | The port to listen on (default "8080") |
| -client  | Gracenote Client ID                    |
| -user    | Gracenote User ID                      |

### General Environment Variables
| Variable             | Description                                                 |
|----------------------|-------------------------------------------------------------|
| ADDR                 | The address to listen on                                    |
| PORT                 | The port to listen on                      |
| CLIENT_KEY           | Gracenote Client ID                                         |
| USER_KEY             | Gracenote User ID                                           |

`CLIENT_KEY` and `USER_KEY` can be configured for App Engine using the `gracenote.yaml` file in `cmd/cddb`.

### App Engine Environment Variables
**Note that the end user typically does not have to set these, they're handled by Google App Engine.**

| Variable             | Description                                                 |
|----------------------|-------------------------------------------------------------|
| USING_APPENGINE      | Whether the program should use App Engine specific features |
| GOOGLE_CLOUD_PROJECT | The Project ID associated with your application             |
| GAE_SERVICE          | The service name of your application                        |
| GAE_VERSION          | The version label of the current application                |

### Compile-time configuration
Go provides a way to set variable values at compile time, this can be used to easily change the default values for variables.
This can be done using the [-X linker flag](https://golang.org/cmd/link/) during building.

This example will set the default address to `127.0.0.1` and the default port to `8888`:
`go build -ldflags "-X main.addr=127.0.0.1 -X main.port=8888"`

| importpath.name                               | Description              |
|-----------------------------------------------|--------------------------|
| main.addr                                     | The address to listen on |
| main.port                                     | The port to listen on    |
| github.com/Hakkin/cddb/cmd/cddb/config.Client | Gracenote Client ID      |
| github.com/Hakkin/cddb/cmd/cddb/config.User   | Gracenote User ID        |

## Building + Running / Deploying
The proxy can be ran as either a standalone program or through Google App Engine.

### Standalone
Simply run `go build` in the `cmd/cddb` directory.
Once built, you can start the proxy by running `cddb`.
You must supply valid Gracenote User and Client IDs (see [configuration](#configuration)).

### Google App Engine
There are multiple configuration files in `cmd/cddb` you must fill out before deploying to Google App Engine.

|  File      |  Description                                             |
|----------------|----------------------------------------------------------|
| gracenote.yaml | Contains your Gracenote Client and User IDs              |
| appengine.json | Contains your App Engine credentials file in JSON format |

Once these are filled out, you can deploy the proxy using `gcloud app deploy` from `cmd/cddb`.

## Usage

Once running, you can use the proxy by settings the FreeDB server in your client to `http://[domain]/cddb`, where `[domain]` is the address or domain the proxy is listening on (see [configuration](#configuration)).

Some clients require additional configuration, pages have been created to help configure various clients:
- [Exact Audio Copy (EAC)](https://github.com/Hakkin/cddb/wiki/Exact-Audio-Copy-Configuration)
- [Mp3tag](https://github.com/Hakkin/cddb/wiki/Mp3tag-Configuration)
- [foobar2000](https://github.com/Hakkin/cddb/wiki/foobar2000-Configuration)

### Advanced Usage
#### Parameters

The proxy allows you to set certain parameters using the URL path:
`http://[domain]/cddb[/language][/country]`

- language
  - Specifies the preferred language for the returned metadata.  
    For the proxy, this only affects the returned genre name.  
    This should be a 3 character [ISO 639-2](https://en.wikipedia.org/wiki/List_of_ISO_639-2_codes) code. 
- country
  - Specifies the country to use for the "genre hierarchy".  
    This should be a 3 character [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) code.  
    Please refer to the examples below and the Gracenote API documentation for the specific use of this parameter.

Examples:

---
The default response:
```
http://[domain]/cddb
Genre: Asian Hip-Hop/Rap
```

---

Setting the language to German:

```
http://[domain]/cddb/ger
Genre: Asiatischer Hip-Hop/Rap
```
This simply translates the genre into the preferred language, but doesn't otherwise modify it.

---

Setting the language to English and the country to Japan:

```
http://[domain]/cddb/eng/jpn
Genre: Hip-Hop/Rap
```
Notice how with the country set to Japan, the genre no longer has the "Asian" specifier.

---

You can read more about the effects of these parameters on Gracenote's WebAPI Documentation: [`language`](https://developer.gracenote.com/sites/default/files/web/webapi/Content/music-web-api/Setting%20the%20Language%20Preference.html) + [`country`](https://developer.gracenote.com/sites/default/files/web/webapi/Content/music-web-api/Specifying%20a%20Country%20Specific.html)
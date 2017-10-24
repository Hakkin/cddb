Gracenote → FreeDB Proxy
-

### What is this?
This is a proxy that allows FreeDB aware software to access metadata from the Gracenote CD database.

The root `cddb` package can also be used as a rudimentary Gracenote library.

### How do I use it?
This proxy can be ran as either a standalone executable (self-hosting) or on Google App Engine.

You can build the standalone executable using the build-\* shell/batch files in the `app/client` directory, or by manually running `go build standalone.go` inside the `app/client` directory.

By default, the standalone version will bind to `127.0.0.1:8080`, you can change this in `app/client/standalone.go`

To deploy on Google App Engine, deploy using `gcloud app deploy` from the `app/client` directory.

Once the server is running, simply set the FreeDB server in your client to:
`http://[domain]/cddb`

For the standalone executable, `[domain]` will be `127.0.0.1:8080` by default, for Google App Engine, it will be the AppSpot domain for your project (or a custom domain if you have it configured).

### Configuration
You must provide your Gracenote Client and User ID in `app/config/config.go` for the service to work.  
You can find more about this at https://developer.gracenote.com/web-api

### Parameters
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

##### Parameter Examples
The default response:

    http://[domain]/cddb
    Genre: Asian Hip-Hop/Rap

---

Setting the language to German:

    http://[domain]/cddb/ger
    Genre: Asiatischer Hip-Hop/Rap
This simply translates the genre into the preferred language, but doesn't otherwise modify it.

---

Setting the language to English and the country to Japan:

    http://[domain]/cddb/eng/jpn
    Genre: Hip-Hop/Rap
Notice how with the country set to Japan, the genre no longer has the "Asian" specifier.

---

You can read more about the effects of these parameters on Gracenote's WebAPI Documentation: [`language`](https://developer.gracenote.com/sites/default/files/web/webapi/Content/music-web-api/Setting%20the%20Language%20Preference.html) + [`country`](https://developer.gracenote.com/sites/default/files/web/webapi/Content/music-web-api/Specifying%20a%20Country%20Specific.html)
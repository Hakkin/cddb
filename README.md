This service proxies FreeDB (CDDB1) requests to Gracenote (CDDB2), and returns the results in a FreeDB friendly format.  
This makes any client that is able to communicate with FreeDB also able to communicate with Gracenote.

## Running
This service can be ran as either a standalone executable (self-hosting) or on Google App Engine  
To run the service standalone, simply run `go build` or `go install` from the `app` directory, then start `app`  
By default, the standalone version will run on port `8080`, you can change this in `app.go`

To deploy on Google App Engine, deploy from the `app` directory.

## Configuration

You must provide your Gracenote Client and User ID in config.go for the service to work.  
You can find more about this at https://developer.gracenote.com/web-api

## Parameters

This service allowed you to configure certain parameters of the requests.  
You can set the *language* and *country* parameters of the Gracenote request using the URL path.

`http://[domain]/cddb/language/country`

*language* should be a 3-length [*ISO 639-2*](https://en.wikipedia.org/wiki/List_of_ISO_639-2_codes) code.  
*country* should be a 3-length [*ISO 3166-1 alpha-3*](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3#Current_codes) code.  
You can read more about the effect of these parameters on Gracenote's WebAPI documentation: [`language`](https://developer.gracenote.com/sites/default/files/web/webapi/Content/music-web-api/Setting%20the%20Language%20Preference.html), [`country`](https://developer.gracenote.com/sites/default/files/web/webapi/Content/music-web-api/Specifying%20a%20Country%20Specific.html)


## Example of parameters

These show a few examples of the affects the parameters have one the responses from Gracenote

Request to */cddb*
```
<ALBUM ORD="1">
	<GN_ID>100670572-A296E2E68B5CB3590EDD5BEC3CE3D6BE</GN_ID>
	<ARTIST>Amos Lee</ARTIST>
	<TITLE>Colours</TITLE>
	<PKG_LANG>ENG</PKG_LANG>
	<GENRE NUM="105245" ID="35493">Western Pop</GENRE>
</ALBUM>
```

Request to */cddb/jpn*
```
<ALBUM ORD="1">
	<GN_ID>100670572-A296E2E68B5CB3590EDD5BEC3CE3D6BE</GN_ID>
	<ARTIST>Amos Lee</ARTIST>
	<TITLE>Colours</TITLE>
	<PKG_LANG>ENG</PKG_LANG>
	<GENRE NUM="105245" ID="35493">ポップ (洋楽)</GENRE>
</ALBUM>
```

Request to */cddb/eng/jpn*
```
<ALBUM ORD="1">
	<GN_ID>100670572-A296E2E68B5CB3590EDD5BEC3CE3D6BE</GN_ID>
	<ARTIST>Amos Lee</ARTIST>
	<TITLE>Colours</TITLE>
	<PKG_LANG>ENG</PKG_LANG>
	<GENRE NUM="70156" ID="28782">General Rock</GENRE>
</ALBUM>
```
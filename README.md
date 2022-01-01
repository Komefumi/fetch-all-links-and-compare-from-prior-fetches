# Fetch All Links And Compare From Prior Fetches

This program can be run with go run . [(space seperated url names)]

A directory called 'fetches' would be made. Each sub directory would of the webpages fetched
Make sure to prefix the links you provide with http:// or https://

Each time a fetch is made, an text file would be created that holds the contents of the fetch

And after all fetches are done, each of the fetches saved for each page will be displayed. With the name (which is really a timestamp)
as well as the number of bytes in that fetched html file

A simple way to run this would be:

```bash
go run . http://golang.org http://gopl.io http://godoc.org
```

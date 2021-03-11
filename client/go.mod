module gitlab.com/256/DebateFrame/client

require (
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/davecgh/go-xdr v0.0.0-20161123171359-e6a2ba005892
	github.com/dennwc/dom v0.2.2-0.20190308181223-8ccb4f24fd8d
	github.com/google/uuid v1.1.1
	github.com/montanaflynn/stats v0.5.0
	github.com/pkg/errors v0.8.1
	gitlab.com/256/WebFrame/dyndom v0.0.0
	gitlab.com/256/WebFrame/waquery v0.0.0
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
)

replace gitlab.com/256/WebFrame/dyndom v0.0.0 => ../../WebFrame/dyndom

replace gitlab.com/256/WebFrame/waquery v0.0.0 => ../../WebFrame/waquery

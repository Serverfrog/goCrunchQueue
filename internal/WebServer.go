package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/hellofresh/health-go/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"mime"
	"os"
	"strings"
)

type GetWidget struct {
	WidgetId string `json:"WidgetId"`
}

var AssetsBasePath, _ = os.Getwd()
var staticResources = make(map[string]*CachedHtml)

const acceptEncoding = "Accept-Encoding"
const contentEncoding = "content-encoding"

func returnAsset(ctx *fasthttp.RequestCtx) {
	if bytes.Contains(ctx.Path(), []byte("..")) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	paths := strings.SplitAfter(string(ctx.Path()), "/assets/")
	if len(paths) != 2 || paths[1] == "" {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	staticResource, exists := staticResources[paths[1]]

	if !exists {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	chooseCompression(ctx, staticResource)

	ext := strings.SplitAfter(paths[1], ".")
	if len(paths) != 2 {
		return
	}
	extension := fmt.Sprintf(".%v", ext[len(ext)-1])

	ctx.SetContentType(mime.TypeByExtension(extension))
}

func chooseCompression(ctx *fasthttp.RequestCtx, cachedHtml *CachedHtml) {
	acceptedEncoding := string(ctx.Request.Header.Peek(acceptEncoding))
	if strings.Contains(acceptedEncoding, "br") {
		ctx.Response.Header.Add(contentEncoding, "br")
		ctx.SetBody(cachedHtml.brotli)
	} else if strings.Contains(acceptedEncoding, "gzip") {
		ctx.Response.Header.Add(contentEncoding, "gzip")
		ctx.SetBody(cachedHtml.gz)

	} else {
		ctx.SetBody(cachedHtml.uncompressed)
	}
}

func createResourceCache() {
	log.Print("Create static Resource Cache")
	assetsPath := fmt.Sprintf("%v/assets/", AssetsBasePath)
	directory := HandleError(os.ReadDir(assetsPath))

	for _, entry := range directory {
		if entry.IsDir() {
			continue
		}
		filePath := fmt.Sprintf("%v/%v", assetsPath, entry.Name())

		byteArray := HandleError(os.ReadFile(filePath))
		staticResources[entry.Name()] = createCache(byteArray)
		log.Infof("Created Cache for: %v", filePath)
	}

}

// api listens on `/api` and accepts a QueueItem
// This will then be put into the Queue, and if all is successfull it will return OK
func addElement(ctx *fasthttp.RequestCtx) {
	var toAdd QueueItem
	HandleErrorB(json.Unmarshal(ctx.PostBody(), &toAdd))
	toAdd = createQueueItem(toAdd.Name, toAdd.CrunchyrollUrl)
	queue.Push(toAdd)
	handleEvent(Event{
		id:   Added,
		item: toAdd,
	})
}

func getAll(ctx *fasthttp.RequestCtx) {
	ctx.SetBody(HandleError(json.Marshal(queue.GetAll())))
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("text/json")
}

func getCurrentProcessed(ctx *fasthttp.RequestCtx) {

}

func buildHandler(version string) *router.Router {
	// add some checks on instance creation
	h, _ := health.New(health.WithComponent(health.Component{
		Name:    "goCrunchQueue",
		Version: version,
	}))

	router := router.New()
	router.GET("/ui", returnAsset)
	router.POST("/api/add", addElement)
	router.GET("/api/all", getAll)
	router.GET("/api/current", getCurrentProcessed)
	router.GET("/status", fasthttpadaptor.NewFastHTTPHandler(h.Handler()))
	return router
}

func wrapHandler(router *router.Router) fasthttp.RequestHandler {
	p := NewPrometheus("fasthttp")
	fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	fastpHandler := p.WrapHandler(router)
	return fastpHandler
}

func StartServer(version string) {
	port := fmt.Sprintf(":%d", configuration.Port)

	createResourceCache()
	log.Printf("Starting server on port %v", port)
	fastpHandler := wrapHandler(buildHandler(version))
	HandleCriticalError(fasthttp.ListenAndServe(port, fastpHandler))
}

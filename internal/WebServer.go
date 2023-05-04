package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/fasthttp/websocket"
	"github.com/hellofresh/health-go/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"mime"
	"os"
	"strconv"
	"strings"
	"time"
)

var AssetsBasePath, _ = os.Getwd()
var staticResources = make(map[string]*CachedHtml)

const acceptEncoding = "Accept-Encoding"
const contentEncoding = "content-encoding"

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod  = (pongWait * 9) / 10
	eventPeriod = 250 * time.Millisecond
)

func returnAsset(ctx *fasthttp.RequestCtx) {
	DebugLog("Serving an Assets under /ui")
	if bytes.Contains(ctx.Path(), []byte("..")) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	DebugLogf("Serving %v", string(ctx.Path()))
	paths := strings.SplitAfter(string(ctx.Path()), "/ui/")
	DebugLogf("Have Split the Paths. %v", paths)
	if len(paths) != 2 || paths[1] == "" {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	resourcePath := fmt.Sprintf("/%v", paths[1])
	staticResource, exists := staticResources[resourcePath]

	if !exists {
		DebugLogf("Resource \"%v\" does not exists in %v", resourcePath, staticResources)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	chooseCompression(ctx, staticResource)

	ext := strings.SplitAfter(resourcePath, ".")
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
	assetsPath := fmt.Sprintf("%v/ui", AssetsBasePath)
	createResourceCacheForDirectory(assetsPath, "")
	DebugLog("Created Resource Cache")
}

func createResourceCacheForDirectory(baseDir string, root string) {

	directory := HandleError(os.ReadDir(baseDir))
	for _, entry := range directory {
		if entry.IsDir() {
			createResourceCacheForDirectory(fmt.Sprintf("%v/%v", baseDir, entry.Name()), fmt.Sprintf("%v/%v", root, entry.Name()))
			continue
		}
		filePath := fmt.Sprintf("%v/%v", baseDir, entry.Name())

		byteArray := HandleError(os.ReadFile(filePath))
		staticResources[fmt.Sprintf("%v/%v", root, entry.Name())] = createCache(byteArray)
		log.Infof("Created Cache for: %v", filePath)
	}
}

// api listens on `/api` and accepts a QueueItem
// This will then be put into the Queue, and if all is successfully it will return OK
func addElement(ctx *fasthttp.RequestCtx) {
	var toAdd QueueItem
	HandleErrorB(json.Unmarshal(ctx.PostBody(), &toAdd))
	toAdd = createQueueItem(toAdd.Name, toAdd.CrunchyrollUrl)
	queue.Push(toAdd)
	eventHandler.handleEvent(Event{
		Id:      Added,
		Item:    toAdd,
		Message: fmt.Sprintf("Item Added via REST API :%v, Name:%v, Url:%v", toAdd.Id, toAdd.Name, toAdd.CrunchyrollUrl),
	})
}

func getAll(ctx *fasthttp.RequestCtx) {
	ctx.SetBody(HandleError(json.Marshal(queue.GetAll())))
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("text/json")
}

func getCurrentProcessed(ctx *fasthttp.RequestCtx) {
	worker.mux.Lock()
	defer worker.mux.Unlock()
	ctx.SetBody(HandleError(json.Marshal(worker.currentItem)))
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("text/json")
}
func getStdLog(ctx *fasthttp.RequestCtx) {
	getLog(ctx, "out")
}
func getErrLog(ctx *fasthttp.RequestCtx) {
	getLog(ctx, "err")
}

func getLog(ctx *fasthttp.RequestCtx, logType string) {
	id := ctx.UserValue("id")
	path := fmt.Sprintf("%v/%v-%v.txt", configuration.LogDestination, id, logType)
	fileContent := HandleError(os.ReadFile(path))
	ctx.SetContentType("text/plain")
	ctx.SetBody(fileContent)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func webSocket(ctx *fasthttp.RequestCtx) {
	err := upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
		var lastMod time.Time
		if n, err := strconv.ParseInt(string(ctx.FormValue("lastMod")), 16, 64); err == nil {
			lastMod = time.Unix(0, n)
		}

		go writer(ws, lastMod)
		reader(ws)
	})

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			log.Println(err)
		}
		return
	}
}
func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	HandleErrorB(ws.SetReadDeadline(time.Now().Add(pongWait)))
	ws.SetPongHandler(func(string) error { HandleErrorB(ws.SetReadDeadline(time.Now().Add(pongWait))); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, lastMod time.Time) {
	pingTicker := time.NewTicker(pingPeriod)
	eventTicker := time.NewTicker(eventPeriod)
	defer func() {
		pingTicker.Stop()
		HandleErrorB(ws.Close())
	}()
	for {
		select {
		case <-pingTicker.C:
			HandleErrorB(ws.SetWriteDeadline(time.Now().Add(writeWait)))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case <-eventTicker.C:
			for event, value := range wsEventListener.getEvents() {
				HandleErrorB(ws.SetWriteDeadline(time.Now().Add(writeWait)))
				HandleErrorB(ws.WriteMessage(websocket.TextMessage, value))
				log.Debugf("Send Event %v for Item %v", event.Id, event.Item.Id)
			}
		}
	}
}

func notFoundHandler(ctx *fasthttp.RequestCtx) {
	DebugLogf("Not found. %v", string(ctx.Path()))
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func redirectToIndexHtml(ctx *fasthttp.RequestCtx) {
	DebugLog("Redirect to Index.html")
	ctx.Redirect("/ui/index.html", fasthttp.StatusMovedPermanently)
}

func buildHandler(version string) *router.Router {
	// add some checks on instance creation
	h, _ := health.New(health.WithComponent(health.Component{
		Name:    "goCrunchQueue",
		Version: version,
	}))

	router := router.New()
	router.POST("/api/add", addElement)
	router.GET("/api/all", getAll)
	router.GET("/api/current", getCurrentProcessed)
	router.GET("/api/log/std/{id}", getStdLog)
	router.GET("/api/log/err/{id}", getErrLog)
	router.ANY("/ws", webSocket)
	router.GET("/status", fasthttpadaptor.NewFastHTTPHandler(h.Handler()))
	router.GET("/ui/", redirectToIndexHtml)
	router.GET("/ui/{filepath:*}", returnAsset)
	router.NotFound = notFoundHandler
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

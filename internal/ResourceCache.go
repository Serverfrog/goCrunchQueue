package internal

import (
	"bytes"
	"compress/gzip"
	"github.com/andybalholm/brotli"
)

type CachedHtml struct {
	brotli       []byte
	gz           []byte
	uncompressed []byte
}

func createCache(rendered []byte) *CachedHtml {
	widgetCache := CachedHtml{uncompressed: rendered}
	widgetCache.brotli = createBrotli(widgetCache)
	widgetCache.gz = createGunzip(widgetCache)
	return &widgetCache
}

func createBrotli(widget CachedHtml) []byte {
	writer := bytes.NewBufferString("")
	brotliWriter := brotli.NewWriter(writer)
	_, err := brotliWriter.Write(widget.uncompressed)
	HandleErrorB(err)
	HandleErrorB(brotliWriter.Flush())
	HandleErrorB(brotliWriter.Close())
	return writer.Bytes()
}
func createGunzip(widget CachedHtml) []byte {
	writer := bytes.NewBufferString("")
	gzipWriter := gzip.NewWriter(writer)
	_, err := gzipWriter.Write(widget.uncompressed)
	HandleErrorB(err)
	HandleErrorB(gzipWriter.Flush())
	HandleErrorB(gzipWriter.Close())
	return writer.Bytes()
}

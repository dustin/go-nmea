package main

import (
	"io"
	"log"
	"os"

	"text/template"

	"github.com/dustin/go-nmea"
)

const kmlHeader = `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
<Folder>`
const kmlPoint = `<Placemark>
  <TimeStamp>{{.TS}}</TimeStamp>
  <Point><coordinates>{{.Lon}},{{.Lat}}</coordinates></Point>
  <gx:flyToMode>smooth</gx:flyToMode>
</Placemark>
`
const kmlFooter = `</Folder></kml>`

const tsFormat = "2006-01-02T15:04:05Z"

var tmpl = template.Must(template.New("").Parse(kmlPoint))

type errRememberer struct {
	w   io.WriteCloser
	err error
}

func (e errRememberer) Write(b []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}

	var n int
	n, e.err = e.w.Write(b)

	return n, e.err
}

func (e errRememberer) Close() error {
	if e.err != nil {
		return e.err
	}
	return e.w.Close()
}

type kmlWriter struct {
	w errRememberer
}

func (k kmlWriter) HandleRMC(m nmea.RMC) {
	tmpl.Execute(k.w, struct {
		Lon, Lat float64
		TS       string
	}{m.Longitude, m.Latitude, m.Timestamp.Format(tsFormat)})
}

func (k kmlWriter) Init() error {
	k.w.Write([]byte(kmlHeader))
	return k.w.err
}

func (k kmlWriter) Close() error {
	k.w.Write([]byte(kmlFooter))
	return k.w.Close()
}

func main() {
	h := kmlWriter{errRememberer{w: os.Stdout}}
	h.Init()
	err := nmea.Process(os.Stdin, h, func(s string, err error) error {
		if s != "" && err != nil {
			log.Printf("On %q: %v", s, err)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error processing stuff: %v", err)
	}
	h.Close()
}

package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"sync/atomic"
)

type HostData struct {
	Hostname   string
	Count      int64
	Interfaces []net.Interface
	Addresses  map[string]string
}

type RequestData struct {
	RemoteAddr string
}

type TemplateData struct {
	H HostData
	R RequestData
}

const housecat = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Yes, It's Caturday</title>
  </head>
  <body>
		<h2>Yes, It's Caturday!</h2>
		<p><strong>HINT:</strong> It's almost ALWAYS caturday.</p>
    <p>
    <img class="" src="//thecatapi.com/api/images/get?format=src&amp;type=gif" alt="Lovely kitten"></p>
    <table>
      <tr>
        <td><strong>Hostname</strong></td>
        <td>{{ .H.Hostname }}</td>
      </tr>
      <tr>
        <td><strong>Count renders</strong></td>
        <td>{{ .H.Count }}</td>
      </tr>
      <tr>
        <td><strong>Remote</strong></td>
        <td>{{ .R.RemoteAddr }}</td>
      </tr>
    </table>
		<br>
		<div>
			<strong>Interfaces</strong>
			<table>
			<tr>
			<th>Name</th>
			<th>HwAddr</th>
			<th>MTU</th>
			<th>Index</th>
			</tr>
			{{ range $i := .H.Interfaces }}
				<tr>
				<td>{{ $i.Name }}</td>
				<td>{{ $i.HardwareAddr }}</td>
				<td>{{ $i.MTU }}</td>
				<td>{{ $i.Index }}</td>
				</tr>
			{{ end }}
			</table>
		</div>
		<br>
		<div>
			<strong>Addresses</strong>
			<table>
			{{ range $addr, $if := .H.Addresses }}
				<tr>
					<td>{{ $if }}</td>
					<td>{{ $addr }}</td>
				</tr>
			{{ end }}
			</table>
		</div>
		<p>Meeeeow :3 - you can use this in your tests too, see: <a href="https://github.com/fntlnz/caturday">https://github.com/fntlnz/caturday</a></p>
  </body>
</html>
`

var (
	hostData HostData
	healthy  int32
)

func ips() map[string]string {
	ips := map[string]string{}
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ips[ip.String()] = i.Name
		}
	}
	return ips
}

func handler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&hostData.Count, 1)

	reqData := RequestData{RemoteAddr: r.RemoteAddr}
	t := template.New("kittens")
	t, _ = t.Parse(housecat)
	t.Execute(w, TemplateData{H: hostData, R: reqData})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func main() {
	atomic.StoreInt32(&healthy, 0)
	log.Println("Starting caturday...")

	hostname, _ := os.Hostname()
	ifaces, _ := net.Interfaces()
	addrs := ips()
	hostData = HostData{Hostname: hostname, Count: 0, Interfaces: ifaces, Addresses: addrs}

	http.HandleFunc("/", handler)
	http.HandleFunc("/healthz", healthHandler)

	atomic.StoreInt32(&healthy, 1)
	log.Println("Initializing the HTTP server")
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		atomic.StoreInt32(&healthy, 0)
		log.Printf("Error while starting caturday: %v\n", err)
	}

	atomic.StoreInt32(&healthy, 0)
	log.Println("caturday has shut down")
}

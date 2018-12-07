package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/influxdata/platform/kit/prom"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type HostData struct {
	Hostname   string
	Counter    prometheus.Counter
	Interfaces []net.Interface
	Addresses  map[string]string
}

type RequestData struct {
	RemoteAddr string
}

type TemplateData struct {
	H     HostData
	R     RequestData
	Count int
}

const houserawcat = `
YES! It's Caturday!
             *     ,MMM8&&&.            *
                  MMMM88&&&&&    .
                 MMMM88&&&&&&&
     *           MMM88&&&&&&&&
                 MMM88&&&&&&&&
                 'MMM88&&&&&&'
                   'MMM8&&&'      *
          |\___/|
          )     (             .              '
         =\     /=
           )===(       *
          /     \
          |     |
         /       \
         \       /
  _/\_/\_/\__  _/_/\_/\_/\_/\_/\_/\_/\_/\_/\_
  |  |  |  |( (  |  |  |  |  |  |  |  |  |  |
  |  |  |  | ) ) |  |  |  |  |  |  |  |  |  |
  |  |  |  |(_(  |  |  |  |  |  |  |  |  |  |
  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |
  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |

Hostname: {{ .H.Hostname }}
Count: {{ .Count }}
Remote Addr: {{ .R.RemoteAddr }}

Interfaces
-----------
{{ range $i := .H.Interfaces }}
{{ $i.Name }}  {{ $i.HardwareAddr }} {{ $i.MTU }} {{ $i.Index }}
{{ end }}


Addresses
----------
{{ range $addr, $if := .H.Addresses }}
{{ $if }} {{ $addr }}
{{ end }}


Meeeeow :3 - you can use this in your tests too, see: https://github.com/fntlnz/caturday
`

const housecat = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Yes, It's Caturday</title>
		<link href="data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAACoopsAAP9ZAG5JIQD/AMMAEQD/AH11awAKBgEA/1EAAABz/wDmkdkA/O7eAAD/9wD/xPYAkcXrAAAAAAAAAAAAoiIiIiIiIioiIiIioiIiIiIqIiIiIiIiIiIiIiIiKiKiIiKiIiIiIjMzMzMAIAIid3d3d9zFVVIRERcRzMlVkru7ERvMxVVSiIiIiMzGzGJERERE3MzNIiIiIiIiIiIioiKiKiIiIqIiIiIiIiIiIiIiIiIiKiIiKiIioiIiIiIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" rel="icon" type="image/x-icon" />
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
        <td>{{ .Count }}</td>
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
	hostData.Counter.Inc()
	cv := counterValue(hostData.Counter)
	reqData := RequestData{RemoteAddr: r.RemoteAddr}
	t := template.New("kittens")
	t, _ = t.Parse(housecat)
	t.Execute(w, TemplateData{H: hostData, R: reqData, Count: cv})
	log.Printf("received request: %s - counter: %d", r.RemoteAddr, cv)
}

func rawHandler(w http.ResponseWriter, r *http.Request) {
	hostData.Counter.Inc()
	cv := counterValue(hostData.Counter)
	reqData := RequestData{RemoteAddr: r.RemoteAddr}
	t := template.New("kittensraw")
	t, _ = t.Parse(houserawcat)
	t.Execute(w, TemplateData{H: hostData, R: reqData, Count: cv})
	log.Printf("received request: %s - counter: %d", r.RemoteAddr, cv)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("KO"))
}

func counterValue(counter prometheus.Counter) int {
	dm := &dto.Metric{}
	counter.Write(dm)
	return int(dm.Counter.GetValue())
}

func main() {
	reg := prom.NewRegistry()
	reg.MustRegister(prometheus.NewGoCollector())
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "caturday",
		Name:      "requests_count",
		Help:      "Number of requests going to a caturday instance",
	})
	reg.MustRegister(counter)

	log.Println("Starting caturday...")

	hostname, _ := os.Hostname()
	ifaces, _ := net.Interfaces()
	addrs := ips()
	hostData = HostData{Hostname: hostname, Counter: counter, Interfaces: ifaces, Addresses: addrs}

	http.HandleFunc("/", handler)
	http.HandleFunc("/raw", rawHandler)
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/health", healthHandler)
	http.Handle("/metrics", reg.HTTPHandler())

	atomic.StoreInt32(&healthy, 1)
	log.Println("Initializing the HTTP server")
	if err := http.ListenAndServe(":8080", nil); err != http.ErrServerClosed {
		atomic.StoreInt32(&healthy, 0)
		log.Printf("Error while starting caturday: %v\n", err)
	}

	atomic.StoreInt32(&healthy, 0)
	log.Println("caturday has shut down")
}

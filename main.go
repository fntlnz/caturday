package main

import (
	"html/template"
	"net"
	"net/http"
	"os"
)

type HostData struct {
	Hostname   string
	Count      int
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
    <p>
    <img class="" src="http://thecatapi.com/api/images/get?format=src&type=gif" alt="Lovely kitten"></p>
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
  </body>
</html>
`

var hostData HostData

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
	hostData.Count = hostData.Count + 1

	reqData := RequestData{RemoteAddr: r.RemoteAddr}
	t := template.New("kittens")
	t, _ = t.Parse(housecat)
	t.Execute(w, TemplateData{H: hostData, R: reqData})
}

func main() {
	hostname, _ := os.Hostname()
	ifaces, _ := net.Interfaces()
	addrs := ips()
	hostData = HostData{Hostname: hostname, Count: 0, Interfaces: ifaces, Addresses: addrs}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

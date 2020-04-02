package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

//func main() {
//	fs := http.FileServer(http.Dir("."))
//	http.Handle("/", fs)
//	log.Println("Listening on :3000...")
//	err := http.ListenAndServe(":3000", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	byts, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(byts))
	fmt.Fprintf(w, "%s", file)

}

var file = `
<html>
<body>
<script src="https://cdn.jsdelivr.net/npm/mermaid@8.4.0/dist/mermaid.min.js"></script>
<script>mermaid.initialize({startOnLoad:true});</script>

Here is one mermaid diagram:
<div class="mermaid">
    graph TD
    A[Client] --> B[Load Balancer]
    B --> C[Server1]
    B --> D[Server2]
</div>

And here is another:
<div class="mermaid">
    graph TD
    A[Client] -->|tcp_123| B(Load Balancer)
    B -->|tcp_456| C[Server1]
    B -->|tcp_456| D[Server2]
</div>
</body>
</html>`

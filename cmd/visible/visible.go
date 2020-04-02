// Command visible is a chromedp example demonstrating how to wait until an
// element is visible.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/chromedp/cdproto/dom"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

var flagPort = flag.Int("port", 8544, "port")
var str string

func main() {
	flag.Parse()
	var output string
	flag.StringVar(&output, "o", "./", "Output directory of documentation")
	flag.Parse()
	filename := flag.Arg(0)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	newTemplate, err := template.New("").Parse(indexHTML)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	if err := newTemplate.Execute(&buf, string(file)); err != nil {
		panic(err)
	}
	println(buf.String())
	// run server
	go testServer(fmt.Sprintf(":%d", *flagPort))

	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	err = chromedp.Run(ctx, visible(fmt.Sprintf("http://localhost:%d", *flagPort), buf.String()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)
}

func visible(host string, mermaidString string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(host),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(makeVisibleScript).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			str, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	}
}

const (
	makeVisibleScript = `setTimeout(function() {
	document.querySelector('#box1').style.display = '';
}, 3000);`
)

// testServer is a simple HTTP server that serves a static html page.
func testServer(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(res, indexHTML)
	})
	return http.ListenAndServe(addr, mux)
}

const indexHTML = `<html>
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

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ts := httptest.NewServer(foo{})
	defer ts.Close()

	const expr = `(function(d, id, v) {
		var b = d.querySelector('body');
		var el = d.createElement('div');
		el.id = id;
		el.innerText = v;
		b.insertBefore(el, b.childNodes[0]);
	})(document, %q, %q);`

	var nodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			s := fmt.Sprintf(expr, "thing", "a new thing!")
			_, exp, err := runtime.Evaluate(s).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
		chromedp.WaitVisible(`#thing`),
	); err != nil {
		panic(err)
	}

	fmt.Println("Document tree:")
	fmt.Print(nodes[0].Dump("  ", "  ", false))
}

type foo struct {
}

func (foo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"

	"char-list/internal/rules"
	"char-list/internal/sheet"
)

//go:embed assets web
var files embed.FS

var registry = rules.NewRegistry(rules.DnD5e{}, rules.Shadowdark{})

var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"dict": func(kv ...any) map[string]any {
		m := map[string]any{}
		for i := 0; i+1 < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
		return m
	},
	"json": func(v any) template.JS {
		b, _ := json.Marshal(v)
		return template.JS(b)
	},
}).ParseFS(files, "web/templates/*.html"))

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	music := flag.String("music", "music", "directory with background music served at /static/music/")
	flag.Parse()
	sheet.Assets = files

	mux := http.NewServeMux()
	static, _ := fs.Sub(files, "web/static")
	mux.Handle("GET /static/music/", http.StripPrefix("/static/music/", http.FileServer(http.Dir(*music))))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))
	mux.HandleFunc("GET /{$}", index)
	mux.HandleFunc("GET /sheet/{ed}", sheetPage)
	mux.HandleFunc("POST /api/{ed}/random", apiRandom)
	mux.HandleFunc("POST /api/{ed}/derive", apiDerive)
	mux.HandleFunc("POST /pdf/{ed}", pdfHandler)
	mux.HandleFunc("GET /pdf/{ed}/random", pdfRandom)

	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func index(w http.ResponseWriter, r *http.Request) {
	render(w, "index.html", map[string]any{"Editions": registry.All()})
}

func sheetPage(w http.ResponseWriter, r *http.Request) {
	ed, ok := registry.Get(r.PathValue("ed"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	render(w, "sheet_"+ed.ID()+".html", map[string]any{"Ed": ed, "Data": ed.Data()})
}

func render(w http.ResponseWriter, name string, data any) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		log.Printf("template %s: %v", name, err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

type randomReq struct {
	Char   json.RawMessage `json:"char"`
	Fields []string        `json:"fields"`
}

func apiRandom(w http.ResponseWriter, r *http.Request) {
	ed, ok := registry.Get(r.PathValue("ed"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	var req randomReq
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	c, logs, err := ed.Random(req.Char, req.Fields)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	raw, _ := json.Marshal(c)
	d, _ := ed.Derive(raw)
	writeJSON(w, map[string]any{"char": c, "rolls": logs, "derived": d})
}

func apiDerive(w http.ResponseWriter, r *http.Request) {
	ed, ok := registry.Get(r.PathValue("ed"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	raw, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	d, err := ed.Derive(raw)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, d)
}

func pdfHandler(w http.ResponseWriter, r *http.Request) {
	ed, ok := registry.Get(r.PathValue("ed"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	raw, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	servePDF(w, r, ed, raw)
}

func pdfRandom(w http.ResponseWriter, r *http.Request) {
	ed, ok := registry.Get(r.PathValue("ed"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	c, _, err := ed.Random(nil, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	raw, _ := json.Marshal(c)
	servePDF(w, r, ed, raw)
}

func servePDF(w http.ResponseWriter, r *http.Request, ed rules.Edition, raw json.RawMessage) {
	var buf bytes.Buffer
	if err := ed.PDF(raw, &buf); err != nil {
		log.Printf("pdf %s: %v", ed.ID(), err)
		http.Error(w, err.Error(), 500)
		return
	}
	var meta struct {
		Name string `json:"name"`
	}
	json.Unmarshal(raw, &meta)
	if meta.Name == "" {
		meta.Name = "character"
	}
	disp := "inline"
	if r.URL.Query().Get("download") == "1" {
		disp = "attachment"
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("%s; filename*=UTF-8''%s.pdf", disp, url.PathEscape(meta.Name)))
	w.Write(buf.Bytes())
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

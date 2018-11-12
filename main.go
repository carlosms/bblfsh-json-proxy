package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	bblfsh "gopkg.in/bblfsh/client-go.v3"

	"gopkg.in/src-d/go-cli.v0"
	log "gopkg.in/src-d/go-log.v1"
)

// version will be replaced automatically by the CI build.
// See https://github.com/src-d/ci/blob/v1/Makefile.main#L56
var (
	name    = "bblfsh-json-proxy"
	version = "undefined"
	build   = "undefined"
)

var app = cli.New(name, version, build, "bblfsh json proxy")

func main() {
	app.AddCommand(&ServeCommand{})

	app.RunMain()
}

type ServeCommand struct {
	cli.PlainCommand `name:"serve" short-description:"serve the app" long-description:"starts serving the application"`
	cli.LogOptions   `group:"Log Options"`
	Host             string `long:"host" env:"BBLFSH_JSON_HOST" default:"0.0.0.0" description:"IP address to bind the HTTP server"`
	Port             int    `long:"port" env:"BBLFSH_JSON_PORT" default:"8095" description:"Port to bind the HTTP server"`
	ServerURL        string `long:"bblfsh" env:"BBLFSH_JSON_SERVER_URL" default:"127.0.0.1:9432" description:"Address where bblfsh server is listening"`
}

func (c *ServeCommand) Execute(args []string) error {
	c.initLog()

	log.With(log.Fields{"version": version, "build": build}).
		Infof("listening on %s:%d", c.Host, c.Port)

	err := c.listenHTTP()
	log.Errorf(err, "")
	return err
}

func (c *ServeCommand) initLog() {
	if c.LogFields == "" {
		bytes, err := json.Marshal(log.Fields{"app": name})
		if err != nil {
			panic(err)
		}
		c.LogFields = string(bytes)
	}

	log.DefaultFactory = &log.LoggerFactory{
		Level:       c.LogLevel,
		Format:      c.LogFormat,
		Fields:      c.LogFields,
		ForceFormat: c.LogForceFormat,
	}
	log.DefaultFactory.ApplyToLogrus()

	log.DefaultLogger = log.New(nil)
}

func (c *ServeCommand) listenHTTP() error {
	addressHTTP := fmt.Sprintf("%s:%d", c.Host, c.Port)

	httpListener, err := net.Listen("tcp", addressHTTP)
	if err != nil {
		return fmt.Errorf("error creating http listener: %s", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/version", c.handleVersion)
	mux.HandleFunc("/languages", c.handleLanguages)
	mux.HandleFunc("/parse", c.handleParse)

	hs := &http.Server{
		Handler: mux,
	}

	if err = hs.Serve(httpListener); err != nil {
		return fmt.Errorf("error starting http server: %s", err)
	}

	return nil
}

func jsonError(w http.ResponseWriter, code int, err error) {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func (c *ServeCommand) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cli, err := bblfsh.NewClient(c.ServerURL)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err)
		return
	}

	resp, err := cli.NewVersionRequest().Do()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	build := ""
	if !resp.Build.IsZero() {
		build = resp.Build.Format(time.RFC3339)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"version": resp.Version,
		"build":   build,
	})
}

func (c *ServeCommand) handleLanguages(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cli, err := bblfsh.NewClient(c.ServerURL)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err)
		return
	}

	resp, err := cli.NewSupportedLanguagesRequest().Do()
	if err != nil {
		jsonError(w, 0, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type uastMode = string

const (
	native    uastMode = "native"
	annotated uastMode = "annotated"
	semantic  uastMode = "semantic"
)

type parseRequest struct {
	Language string   `json:"language"`
	Filename string   `json:"filename"`
	Content  string   `json:"content"`
	Filter   string   `json:"filter"`
	Mode     uastMode `json:"mode"`
}

func (c *ServeCommand) handleParse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req parseRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err)
		return
	}

	cli, err := bblfsh.NewClient(c.ServerURL)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err)
		return
	}

	var mode bblfsh.Mode
	switch req.Mode {
	case native:
		mode = bblfsh.Native
	case annotated:
		mode = bblfsh.Annotated
	case semantic:
		mode = bblfsh.Semantic
	case "":
		mode = bblfsh.Semantic
	default:
		jsonError(w, http.StatusBadRequest, fmt.Errorf(`invalid "mode" %q; it must be one of "native", "annotated", "semantic"`, req.Mode))
		return
	}

	resp, lang, err := cli.NewParseRequest().
		Language(req.Language).
		Filename(req.Filename).
		Content(req.Content).
		Mode(mode).
		UAST()

	if bblfsh.ErrSyntax.Is(err) {
		jsonError(w, http.StatusBadRequest, fmt.Errorf("error parsing UAST: %s", err))
		return
	}
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err)
		return
	}

	if req.Filter != "" {
		jsonError(w, http.StatusNotImplemented, fmt.Errorf("filter is not yet implemented"))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"uast":     resp,
		"language": lang,
	})
}

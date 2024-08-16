package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func main() {
	handler := caddyhttp.StaticResponse{
		StatusCode: caddyhttp.WeakString(strconv.Itoa(http.StatusOK)),
		Body:       "hello world\n",
	}

	app := caddyhttp.App{
		Servers: map[string]*caddyhttp.Server{
			"veritas": {
				Listen: []string{":8080"},
				Routes: caddyhttp.RouteList{
					{
						MatcherSetsRaw: []caddy.ModuleMap{
							{
								"host": caddyconfig.JSON(caddyhttp.MatchHost{"localhost"}, nil),
							},
						},
						HandlersRaw: []json.RawMessage{
							caddyconfig.JSONModuleObject(handler, "handler", "static_response", nil),
						},
					},
				},
			},
		},
	}

	cfg := caddy.Config{
		AppsRaw: caddy.ModuleMap{
			"http": caddyconfig.JSON(app, nil),
		},
	}

	err := caddy.Run(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}

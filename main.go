// An Antresol addon (contract 1, docs/addons/contract.md in the core): the
// server starts this program with its settings in ANTRESOL_* variables; it
// serves HTTP on the Unix socket ANTRESOL_ADDON_SOCKET and must answer
// GET /health within 15 s. The app reaches /api/… through the core's proxy
// (/api/addons/<id>/x/…); the program reaches the core through
// addonkit.Client, and the core calls its hooks (/hooks/<name>).
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antresol-app/antresol/addonkit"
)

func main() {
	env, err := addonkit.FromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(env.DataDir, 0o700); err != nil {
		log.Fatalf("data dir: %v", err)
	}
	os.Remove(env.Socket) // a stale socket of an earlier run
	ln, err := net.Listen("unix", env.Socket)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	ready := &addonkit.Ready{}
	srv := &http.Server{Handler: newHandler(env, addonkit.NewClient(env), ready), ReadHeaderTimeout: 10 * time.Second}

	// the core stops an addon with SIGTERM and kills it 10 s later
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		shut, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shut)
	}()
	// A long start (a data import) goes here, before Set, while /health
	// already answers {"ok": true, "ready": false}; the core API accepts
	// the token meanwhile. This addon has nothing to prepare.
	ready.Set()
	log.Printf("%s: serving on %s", env.ID, env.Socket)
	if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("serve: %v", err)
	}
}

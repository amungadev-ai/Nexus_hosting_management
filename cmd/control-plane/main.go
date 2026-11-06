// Command nexus-control-plane is the panel's HTTP control plane (walking
// skeleton): health, capability report, guided-router steps, static-site
// provisioning and audit-chain inspection — all backed by the Store interface
// so SQLite and PostgreSQL remain equivalent (DEC-001).
package main

import (
        "context"
        "crypto/rand"
        "encoding/hex"
        "encoding/json"
        "errors"
        "flag"
        "fmt"
        "log"
        "net/http"
        "os"
        "path/filepath"
        "regexp"
        "strconv"
        "strings"
        "sync"
        "time"

        "github.com/amungadev-ai/Nexus_hosting_management/internal/audit"
        "github.com/amungadev-ai/Nexus_hosting_management/internal/capability"
        "github.com/amungadev-ai/Nexus_hosting_management/internal/store"
        "github.com/amungadev-ai/Nexus_hosting_management/internal/store/postgres"
        "github.com/amungadev-ai/Nexus_hosting_management/internal/store/sqlite"
)

const version = "0.1.0-skeleton"

func envOr(key, def string) string {
        if v := os.Getenv(key); v != "" {
                return v
        }
        return def
}

func main() {
        addr := flag.String("addr", envOr("NEXUS_LISTEN_ADDR", "127.0.0.1:8080"), "control-plane listen address")
        siteAddr := flag.String("site-addr", envOr("NEXUS_SITE_ADDR", "127.0.0.1:8081"), "static-site listen address (empty disables)")
        driver := flag.String("driver", envOr("NEXUS_DB_DRIVER", "sqlite"), "storage driver: sqlite | postgres")
        dsn := flag.String("dsn", envOr("NEXUS_DB_DSN", ""), "driver DSN (default: <data-dir>/nexus.db for sqlite)")
        dataDir := flag.String("data-dir", envOr("NEXUS_DATA_DIR", "/var/lib/nexuspanel"), "state directory")
        flag.Parse()

        if *dsn == "" {
                if *driver == "sqlite" {
                        *dsn = filepath.Join(*dataDir, "nexus.db")
                } else {
                        log.Fatal("-dsn is required for the postgres driver")
                }
        }
        if err := os.MkdirAll(filepath.Join(*dataDir, "sites"), 0o750); err != nil {
                log.Fatalf("create data dir: %v", err)
        }

        var st store.Store
        var err error
        switch *driver {
        case "sqlite":
                st, err = sqlite.Open(*dsn)
        case "postgres":
                st, err = postgres.Open(*dsn)
        default:
                log.Fatalf("unknown driver %q (sqlite|postgres)", *driver)
        }
        if err != nil {
                log.Fatal(err)
        }
        defer st.Close()

        s := &server{st: st, dataDir: *dataDir, router: newRouterState()}

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        payload := fmt.Sprintf(`{"driver":%q,"addr":%q,"version":%q}`, *driver, *addr, version)
        _, err = s.st.AppendAudit(ctx, "system", "server.started", payload)
        cancel()
        if err != nil {
                log.Fatalf("audit server.started: %v", err)
        }

        if *siteAddr != "" {
                siteMux := http.NewServeMux()
                siteMux.Handle("/", &siteServer{dataDir: *dataDir})
                go func() {
                        log.Printf("site server listening on %s", *siteAddr)
                        if err := http.ListenAndServe(*siteAddr, siteMux); err != nil {
                                log.Fatalf("site server: %v", err)
                        }
                }()
        }

        mux := http.NewServeMux()
        registerRoutes(mux, s)
        log.Printf("control plane %s listening on %s (driver=%s, data-dir=%s)", version, *addr, *driver, *dataDir)
        log.Fatal(http.ListenAndServe(*addr, mux))
}

type server struct {
        st      store.Store
        dataDir string
        router  *routerState
}

func registerRoutes(mux *http.ServeMux, s *server) {
        mux.HandleFunc("GET /healthz", s.handleHealth)
        mux.HandleFunc("GET /readyz", s.handleReady)
        mux.HandleFunc("GET /api/v1/capabilities", s.handleCapabilities)
        mux.HandleFunc("GET /api/v1/audit/tail", s.handleAuditTail)
        mux.HandleFunc("GET /api/v1/audit/verify", s.handleAuditVerify)
        mux.HandleFunc("GET /api/v1/router/steps", s.handleRouterSteps)
        mux.HandleFunc("POST /api/v1/router/steps/{id}/complete", s.handleRouterComplete)
        mux.HandleFunc("GET /api/v1/sites", s.handleListSites)
        mux.HandleFunc("POST /api/v1/sites", s.handleProvisionSite)
        mux.HandleFunc("GET /api/v1/sites/{id}", s.handleGetSite)
        mux.HandleFunc("GET /s/{hostname}/", s.handleServeSiteByPath)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(code)
        _ = json.NewEncoder(w).Encode(v)
}

func (s *server) audit(action, payload string) error {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        _, err := s.st.AppendAudit(ctx, "control-plane", action, payload)
        return err
}

func (s *server) handleHealth(w http.ResponseWriter, _ *http.Request) {
        writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "version": version})
}

func (s *server) handleReady(w http.ResponseWriter, r *http.Request) {
        if err := s.st.Ping(r.Context()); err != nil {
                writeJSON(w, http.StatusServiceUnavailable, map[string]any{"status": "unavailable", "error": err.Error()})
                return
        }
        writeJSON(w, http.StatusOK, map[string]any{"status": "ready"})
}

func (s *server) handleCapabilities(w http.ResponseWriter, _ *http.Request) {
        writeJSON(w, http.StatusOK, capability.Evaluate(capability.Detect()))
}

func (s *server) handleAuditTail(w http.ResponseWriter, r *http.Request) {
        n := 20
        if v := r.URL.Query().Get("n"); v != "" {
                parsed, err := strconv.Atoi(v)
                if err != nil || parsed <= 0 {
                        writeJSON(w, http.StatusBadRequest, map[string]any{"error": "n must be a positive integer"})
                        return
                }
                n = parsed
        }
        if n > 500 {
                n = 500
        }
        recs, err := s.st.AuditTail(r.Context(), n)
        if err != nil {
                writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
                return
        }
        if recs == nil {
                recs = []audit.Record{}
        }
        writeJSON(w, http.StatusOK, map[string]any{"records": recs})
}

func (s *server) handleAuditVerify(w http.ResponseWriter, r *http.Request) {
        if err := s.st.VerifyAuditChain(r.Context()); err != nil {
                writeJSON(w, http.StatusInternalServerError, map[string]any{"valid": false, "error": err.Error()})
                return
        }
        writeJSON(w, http.StatusOK, map[string]any{"valid": true})
}

// --- guided router (DEC-011) ---

type routerStep struct {
        ID    string `json:"id"`
        Title string `json:"title"`
        Done  bool   `json:"done"`
}

var skeletonSteps = []struct{ ID, Title string }{
        {"install", "Panel installed and services running"},
        {"capabilities", "Host capabilities detected and install gates passed"},
        {"router", "Guided router wired to the SPA"},
        {"site", "First static site provisioned"},
        {"audit", "Audit chain verified end to end"},
}

type routerState struct {
        mu   sync.Mutex
        done map[string]bool
}

func newRouterState() *routerState {
        return &routerState{done: map[string]bool{"install": true}} // the installer completed this step by definition
}

func (rs *routerState) list() []routerStep {
        rs.mu.Lock()
        defer rs.mu.Unlock()
        steps := make([]routerStep, 0, len(skeletonSteps))
        for _, st := range skeletonSteps {
                steps = append(steps, routerStep{ID: st.ID, Title: st.Title, Done: rs.done[st.ID]})
        }
        return steps
}

func (rs *routerState) complete(id string) bool {
        rs.mu.Lock()
        defer rs.mu.Unlock()
        for _, st := range skeletonSteps {
                if st.ID == id {
                        rs.done[id] = true
                        return true
                }
        }
        return false
}

func (s *server) handleRouterSteps(w http.ResponseWriter, _ *http.Request) {
        writeJSON(w, http.StatusOK, map[string]any{"steps": s.router.list()})
}

func (s *server) handleRouterComplete(w http.ResponseWriter, r *http.Request) {
        id := r.PathValue("id")
        if !s.router.complete(id) {
                writeJSON(w, http.StatusNotFound, map[string]any{"error": "unknown step " + id})
                return
        }
        payload, _ := json.Marshal(map[string]string{"step": id})
        if err := s.audit("router.step_completed", string(payload)); err != nil {
                writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
                return
        }
        writeJSON(w, http.StatusOK, map[string]any{"steps": s.router.list()})
}

// --- sites ---

var hostnameLabel = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)

func validHostname(h string) bool {
        if h == "" || len(h) > 253 || strings.Contains(h, "..") {
                return false
        }
        for _, label := range strings.Split(h, ".") {
                if !hostnameLabel.MatchString(label) {
                        return false
                }
        }
        return true
}

func randHex(n int) string {
        b := make([]byte, n)
        _, _ = rand.Read(b)
        return hex.EncodeToString(b)
}

func (s *server) handleListSites(w http.ResponseWriter, r *http.Request) {
        sites, err := s.st.ListSites(r.Context())
        if err != nil {
                writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
                return
        }
        if sites == nil {
                sites = []store.Site{}
        }
        writeJSON(w, http.StatusOK, map[string]any{"sites": sites})
}

func (s *server) handleGetSite(w http.ResponseWriter, r *http.Request) {
        site, err := s.st.GetSite(r.Context(), r.PathValue("id"))
        if errors.Is(err, store.ErrNotFound) {
                writeJSON(w, http.StatusNotFound, map[string]any{"error": "site not found"})
                return
        }
        if err != nil {
                writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
                return
        }
        writeJSON(w, http.StatusOK, site)
}

func (s *server) handleProvisionSite(w http.ResponseWriter, r *http.Request) {
        var req struct {
                Hostname string `json:"hostname"`
        }
        if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&req); err != nil {
                writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid JSON body"})
                return
        }
        hostname := strings.ToLower(strings.TrimSpace(req.Hostname))
        if !validHostname(hostname) {
                writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid hostname"})
                return
        }

        root := filepath.Join(s.dataDir, "sites", hostname)
        if err := os.MkdirAll(root, 0o750); err != nil {
                writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
                return
        }
        indexPath := filepath.Join(root, "index.html")
        if _, err := os.Stat(indexPath); os.IsNotExist(err) {
                page := fmt.Sprintf(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><title>%s</title></head>
<body><h1>%s</h1><p>Provisioned by the NexusPanel walking skeleton.</p></body></html>
`, hostname, hostname)
                if err := os.WriteFile(indexPath, []byte(page), 0o640); err != nil {
                        writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
                        return
                }
        }

        now := time.Now().UTC()
        site := store.Site{
                ID:        randHex(8),
                Hostname:  hostname,
                RootPath:  root,
                State:     store.SiteActive,
                CreatedAt: now,
                UpdatedAt: now,
        }
        if err := s.st.CreateSite(r.Context(), site); err != nil {
                if errors.Is(err, store.ErrDuplicate) {
                        writeJSON(w, http.StatusConflict, map[string]any{"error": "hostname already provisioned"})
                        return
                }
                writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
                return
        }
        payload, _ := json.Marshal(map[string]string{"id": site.ID, "hostname": site.Hostname})
        if err := s.audit("site.provisioned", string(payload)); err != nil {
                writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
                return
        }
        writeJSON(w, http.StatusCreated, site)
}

// handleServeSiteByPath serves provisioned sites under /s/{hostname}/ on the
// control plane (convenience for environments without host-based DNS).
func (s *server) handleServeSiteByPath(w http.ResponseWriter, r *http.Request) {
        hostname := r.PathValue("hostname")
        if !validHostname(hostname) {
                http.NotFound(w, r)
                return
        }
        root := filepath.Join(s.dataDir, "sites", hostname)
        if fi, err := os.Stat(root); err != nil || !fi.IsDir() {
                http.NotFound(w, r)
                return
        }
        http.StripPrefix("/s/"+hostname, http.FileServer(http.Dir(root))).ServeHTTP(w, r)
}

// siteServer serves provisioned sites on their own listener, routed by the
// Host header. TLS is added with ACME once a throwaway domain is available
// (walking-skeleton spec, deferred list) — the seam is this handler.
type siteServer struct {
        dataDir string
}

func (ss *siteServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        host := strings.ToLower(strings.SplitN(r.Host, ":", 2)[0])
        if !validHostname(host) {
                http.NotFound(w, r)
                return
        }
        root := filepath.Join(ss.dataDir, "sites", host)
        if fi, err := os.Stat(root); err != nil || !fi.IsDir() {
                http.NotFound(w, r)
                return
        }
        http.FileServer(http.Dir(root)).ServeHTTP(w, r)
}

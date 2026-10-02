package main

import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "log/slog"
 "net/http"
 "os"
 "sync"
 "time"

 "github.com/previewdock/previewdock/internal/deployments"
 "github.com/previewdock/previewdock/internal/github"
)

type deploymentStore struct { sync.Mutex; items map[string]deployments.Deployment }
func (s *deploymentStore) create(d deployments.Deployment) { s.Lock(); defer s.Unlock(); s.items[d.ID] = d }
func (s *deploymentStore) list() []deployments.Deployment { s.Lock(); defer s.Unlock(); out:=make([]deployments.Deployment,0,len(s.items)); for _, d:=range s.items { out=append(out,d) }; return out }

func main() {
 logger:=slog.New(slog.NewJSONHandler(os.Stdout,nil)); secret:=os.Getenv("GITHUB_WEBHOOK_SECRET")
 if secret=="" { logger.Warn("GITHUB_WEBHOOK_SECRET is unset; webhook endpoint rejects all requests") }
 store:=&deploymentStore{items:map[string]deployments.Deployment{}}
 mux:=http.NewServeMux()
 mux.HandleFunc("GET /healthz", func(w http.ResponseWriter,r *http.Request){ jsonResponse(w,200,map[string]string{"status":"ok"}) })
 mux.HandleFunc("GET /readyz", func(w http.ResponseWriter,r *http.Request){ jsonResponse(w,200,map[string]string{"status":"ready"}) })
 mux.HandleFunc("GET /api/v1/deployments", func(w http.ResponseWriter,r *http.Request){ jsonResponse(w,200,store.list()) })
 mux.HandleFunc("POST /api/v1/webhooks/github", func(w http.ResponseWriter,r *http.Request){
  defer r.Body.Close(); body,err:=github.ReadPayload(r.Body,1<<20); if err!=nil { apiError(w,400,"invalid_payload",err.Error()); return }
  if secret=="" || !github.VerifySignature(secret,r.Header.Get("X-Hub-Signature-256"),body) { apiError(w,401,"invalid_signature","webhook signature verification failed"); return }
  if r.Header.Get("X-GitHub-Event")!="pull_request" { jsonResponse(w,202,map[string]string{"status":"ignored"}); return }
  ev,err:=github.ParsePullRequest(body); if err!=nil { apiError(w,400,"invalid_payload",err.Error()); return }
  if ev.Action!="opened" && ev.Action!="synchronize" && ev.Action!="reopened" && ev.Action!="closed" { jsonResponse(w,202,map[string]string{"status":"ignored"}); return }
  if ev.Action=="closed" { jsonResponse(w,202,map[string]string{"status":"cleanup_queued"}); return }
  if ev.Fork { jsonResponse(w,202,map[string]string{"status":"approval_required"}); return }
  d:=deployments.New(ev.Repository,ev.Number,ev.SHA); store.create(d)
  logger.Info("deployment queued", "deployment_id",d.ID,"repository",d.Repository,"pr",d.PullRequest)
  jsonResponse(w,202,d)
 })
 srv:=&http.Server{Addr:env("API_ADDR",":8080"),Handler:requestID(mux),ReadHeaderTimeout:5*time.Second,ReadTimeout:15*time.Second,WriteTimeout:15*time.Second,IdleTimeout:60*time.Second}
 logger.Info("previewdock api listening", "addr",srv.Addr); if err:=srv.ListenAndServe(); !errors.Is(err,http.ErrServerClosed){logger.Error("server stopped", "error",err)}
 _=context.Background()
}
func env(k,d string)string{if v:=os.Getenv(k);v!=""{return v};return d}
func jsonResponse(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_ = json.NewEncoder(w).Encode(v)}
func apiError(w http.ResponseWriter,status int,code,msg string){jsonResponse(w,status,map[string]any{"error":map[string]string{"code":code,"message":msg}})}
func requestID(next http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){id:=r.Header.Get("X-Request-ID");if id==""{sum:=sha256.Sum256([]byte(time.Now().String()));id=hex.EncodeToString(sum[:])[:16]};w.Header().Set("X-Request-ID",id);next.ServeHTTP(w,r)})}

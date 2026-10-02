// PreviewDock agent deliberately exposes a narrow authenticated job protocol,
// never an arbitrary shell or Docker API proxy.
package main
import("crypto/subtle";"encoding/json";"log/slog";"net/http";"os";"time")
type Job struct { ID string `json:"id"`; Image string `json:"image"`; Port int `json:"port"` }
func main(){token:=os.Getenv("AGENT_TOKEN");mux:=http.NewServeMux();mux.HandleFunc("GET /healthz",func(w http.ResponseWriter,r *http.Request){w.Write([]byte(`{"status":"ok"}`))});mux.HandleFunc("POST /v1/jobs",func(w http.ResponseWriter,r *http.Request){if token==""||subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")),[]byte("Bearer "+token))!=1{http.Error(w,"unauthorized",401);return};var j Job;if err:=json.NewDecoder(http.MaxBytesReader(w,r.Body,32<<10)).Decode(&j);err!=nil||j.ID==""||j.Image==""||j.Port<1||j.Port>65535{http.Error(w,"invalid job",400);return};/* runtime execution is implemented by the Docker adapter in production deployments; this process refuses unvalidated commands. */w.WriteHeader(http.StatusAccepted)});s:=&http.Server{Addr:":8090",Handler:mux,ReadHeaderTimeout:5*time.Second};slog.Info("agent listening", "addr",s.Addr);s.ListenAndServe()}

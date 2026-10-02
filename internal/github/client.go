// Package github contains the narrowly scoped GitHub App client used by the controller.
package github

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const defaultAPIURL = "https://api.github.com"

type AppClient struct { appID string; key *rsa.PrivateKey; baseURL string; http *http.Client; now func() time.Time }
type InstallationToken struct { Token string `json:"token"`; ExpiresAt time.Time `json:"expires_at"` }
type Repository struct { ID int64 `json:"id"`; FullName string `json:"full_name"`; Private bool `json:"private"` }

func NewAppClient(appID, privateKeyPEM, baseURL string, h *http.Client) (*AppClient, error) {
	if appID == "" || privateKeyPEM == "" { return nil, fmt.Errorf("github app id and private key are required") }
	b, _ := pem.Decode([]byte(privateKeyPEM)); if b == nil { return nil, fmt.Errorf("invalid github app private key pem") }
	k, err := x509.ParsePKCS1PrivateKey(b.Bytes); if err != nil { if any, e := x509.ParsePKCS8PrivateKey(b.Bytes); e == nil { var ok bool; k, ok = any.(*rsa.PrivateKey); if !ok { return nil, fmt.Errorf("github app key must be RSA") } } else { return nil, fmt.Errorf("parse github app key: %w", err) } }
	if baseURL == "" { baseURL = defaultAPIURL }; if h == nil { h = &http.Client{Timeout: 15 * time.Second} }
	return &AppClient{appID:appID,key:k,baseURL:strings.TrimRight(baseURL,"/"),http:h,now:time.Now},nil
}
func (c *AppClient) appJWT() (string,error) { now:=c.now().UTC(); claims:=map[string]any{"iat":now.Add(-30*time.Second).Unix(),"exp":now.Add(9*time.Minute).Unix(),"iss":c.appID}; return signJWT(claims,c.key) }
func signJWT(claims map[string]any,key *rsa.PrivateKey)(string,error){ h:=base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`)); p,err:=json.Marshal(claims);if err!=nil{return "",err}; body:=h+"."+base64.RawURLEncoding.EncodeToString(p); sum:=sha256.Sum256([]byte(body)); sig,err:=rsa.SignPKCS1v15(rand.Reader,key,crypto.SHA256,sum[:]);if err!=nil{return "",err};return body+"."+base64.RawURLEncoding.EncodeToString(sig),nil }
func (c *AppClient) installationToken(ctx context.Context, installationID int64)(InstallationToken,error){ jwt,err:=c.appJWT();if err!=nil{return InstallationToken{},err};var out InstallationToken;err=c.do(ctx,http.MethodPost,"/app/installations/"+strconv.FormatInt(installationID,10)+"/access_tokens",jwt,nil,&out);return out,err }
func (c *AppClient) Repositories(ctx context.Context, installationID int64)([]Repository,error){ token,err:=c.installationToken(ctx,installationID);if err!=nil{return nil,err};var out struct{Repositories []Repository `json:"repositories"`};err=c.do(ctx,http.MethodGet,"/installation/repositories?per_page=100","token "+token.Token,nil,&out);return out.Repositories,err }
func (c *AppClient) do(ctx context.Context,method,path,authorization string,in,out any)error{var body io.Reader;if in!=nil{b,err:=json.Marshal(in);if err!=nil{return err};body=strings.NewReader(string(b))};req,err:=http.NewRequestWithContext(ctx,method,c.baseURL+path,body);if err!=nil{return err};req.Header.Set("Accept","application/vnd.github+json");req.Header.Set("Authorization","Bearer "+authorization);if strings.HasPrefix(authorization,"token "){req.Header.Set("Authorization",authorization)};if in!=nil{req.Header.Set("Content-Type","application/json")};res,err:=c.http.Do(req);if err!=nil{return err};defer res.Body.Close();if res.StatusCode<200||res.StatusCode>299{b,_:=io.ReadAll(io.LimitReader(res.Body,4096));return fmt.Errorf("github api %s: %s",res.Status,strings.TrimSpace(string(b)))};if out!=nil{return json.NewDecoder(res.Body).Decode(out)};return nil }

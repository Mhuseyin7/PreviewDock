package github

import("crypto/hmac";"crypto/sha256";"crypto/subtle";"encoding/hex";"encoding/json";"fmt";"io")
func ReadPayload(r io.Reader,max int64)([]byte,error){b,err:=io.ReadAll(io.LimitReader(r,max+1));if err!=nil{return nil,err};if int64(len(b))>max{return nil,fmt.Errorf("payload exceeds maximum size")};return b,nil}
func VerifySignature(secret,header string,body []byte) bool { const prefix="sha256="; if len(header)<len(prefix)||header[:len(prefix)]!=prefix{return false}; supplied,err:=hex.DecodeString(header[len(prefix):]);if err!=nil{return false};m:=hmac.New(sha256.New,[]byte(secret));m.Write(body);return subtle.ConstantTimeCompare(supplied,m.Sum(nil))==1 }
type PullRequestEvent struct{Action string `json:"action"`; Number int `json:"number"`; PullRequest struct{Head struct{SHA string `json:"sha"`; Repo *struct{} `json:"repo"`} `json:"head"`} `json:"pull_request"`; Repository struct{FullName string `json:"full_name"`} `json:"repository"`; Fork bool}
func ParsePullRequest(b []byte)(PullRequestEvent,error){var e PullRequestEvent;if err:=json.Unmarshal(b,&e);err!=nil{return e,err};e.Fork=e.PullRequest.Head.Repo==nil;if e.Number<1||e.Repository.FullName==""||e.PullRequest.Head.SHA==""{return e,fmt.Errorf("pull request payload lacks required repository, number, or commit")};return e,nil}

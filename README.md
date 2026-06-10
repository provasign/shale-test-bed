# shale-test-bed

A deliberately small Go HTTP service used to demonstrate
[Shale](https://github.com/provasign/shale) — agent PR evidence — on real
pull requests. Each PR in this repo exercises a different part of the Shale
card: full session evidence, coverage gaps, recorded check failures, and the
no-evidence nudge.

```
POST /login   {"user": "...", "password": "..."}
GET  /health
```

Run: `go run .` · Test: `go test ./...`

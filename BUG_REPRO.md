# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	venue-reservation/cmd/venue-server	[no test files]
?   	venue-reservation/internal/cli	[no test files]
ok  	venue-reservation/internal/access	0.001s
ok  	venue-reservation/internal/export	0.002s
?   	venue-reservation/internal/model	[no test files]
?   	venue-reservation/internal/seed	[no test files]
ok  	venue-reservation/internal/httpapi	0.009s
ok  	venue-reservation/internal/integration	0.027s
--- FAIL: TestBusiness001Regression (0.01s)
    business_regression_test.go:22: unexpected first page: [{ID:RSV-0003 VenueID:A-101 SlotID:A-101-SLOT-01 Applicant:member Purpose:session RSV-0003 Scope:team Status:pending Notes: CreatedAt:sequence UpdatedAt:sequence Revision:1} {ID:RSV-0004 VenueID:A-101 SlotID:A-101-SLOT-01 Applicant:member Purpose:session RSV-0004 Scope:team Status:pending Notes: CreatedAt:sequence UpdatedAt:sequence Revision:1}]
FAIL
FAIL	venue-reservation/internal/service	0.037s
ok  	venue-reservation/internal/store	0.014s
ok  	venue-reservation/internal/timeline	0.001s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/venue-server): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/venue-server): exit `0`

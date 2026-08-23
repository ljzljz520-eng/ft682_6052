# ft682_6052

Go project for module `venue-reservation`.

## Standard commands

```bash
go build ./...
go test -count=1 ./...
```

## Run

```bash
go run ./cmd/venue-server
```

## Docker validation

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh my-go-task linux/arm64
./build_benzhi_docker.sh my-go-task linux/amd64
docker run -it my-go-task:latest
```

## Known initial failures

See `BUG_REPRO.md` for the exact command and output captured during packaging.

## Application

This pure-Go service provides venue and time-slot catalog APIs, reservation submission and paged search, detail history, review, collaboration, archive, reports, import, and CSV export. Start it with:

```bash
go run ./cmd/venue-server -address :8080 -db venue.db -seed=true
```

Tomcat is not required for this Go HTTP service. In a Tomcat-centered environment, run the Go process beside Tomcat and proxy the Tomcat-facing route to `127.0.0.1:8080`; the application persists to the configured bbolt database path.

`TestPersistenceSurvivesReopen` verifies that Venue, TimeSlot, Reservation, AuditEvent, CollaborationNote, and ArchiveRecord survive a close and reopen. `TestBusiness001Regression` intentionally remains red to reproduce the injected page-offset boundary defect.

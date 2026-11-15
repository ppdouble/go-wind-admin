# Kratos v3 compatibility ports

`tx7do/` contains narrowly scoped source copies of the upstream tx7do modules
used by the backend. Their original module paths are retained by `backend/go.mod`
`replace` directives. The copies are needed because those upstream versions still
import Kratos v2; mixing v2 and v3 in the same server registers
`errors/errors.proto` twice and panics before startup.

The changes in the copies are import migrations, the v3 slog adapter, the v3
`grpc.NewClient` call, and an adapter preserving the Aegis BBR rate limiter.
Each module keeps its upstream license. Do not remove a port until its upstream
release can build with Kratos v3 and the server dependency graph contains no v2.

`pkg/v2logcompat` keeps the source level logging interface for these copies;
it does not import Kratos v2 at runtime. `pkg/geolitecompat` reuses the original
GeoLite asset package without copying its 60 MB database into this repository.

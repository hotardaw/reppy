// defines my Go project's module path and dependencies, like a package.json
// go.sum contains cryptographic checksums of specific versions of dependencies, ensuring reproducible builds. like a package-lock.json
module go-reppy/backend

go 1.22.2

require (
	github.com/golang-jwt/jwt/v4 v4.5.1
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.31.0
)

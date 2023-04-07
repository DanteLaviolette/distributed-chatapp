module go.violettedev.com/eecs4222/login

go 1.20

require (
	github.com/gofiber/fiber/v2 v2.42.0
	go.mongodb.org/mongo-driver v1.11.2
	golang.org/x/crypto v0.7.0
)

require (
	go.violettedev.com/eecs4222/shared/auth v1.0.0
	go.violettedev.com/eecs4222/shared/constants v1.0.0
	go.violettedev.com/eecs4222/shared/database v1.0.0
)

replace (
	go.violettedev.com/eecs4222/shared/auth => ../shared/auth
	go.violettedev.com/eecs4222/shared/constants => ../shared/constants
	go.violettedev.com/eecs4222/shared/database => ../shared/database
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/fasthttp/websocket v1.5.1 // indirect
	github.com/gofiber/websocket/v2 v2.1.4 // indirect
	github.com/golang-jwt/jwt/v5 v5.0.0-rc.1 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/savsgio/dictpool v0.0.0-20221023140959-7bf2e61cea94 // indirect
	github.com/savsgio/gotils v0.0.0-20220530130905-52f3993e8d6d // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.44.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)
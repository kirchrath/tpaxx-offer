# tPaxxOffers

## pre requirements
`go swagger (https://goswagger.io)`

`sqlite3`

## generate api
`swagger generate server -A tPaxxOffers`

## build
`go build ./cmd/t-paxx-offers-server`

## run
`./t-paxx-offers-server --port=9988`

## prepare data
`./infxprepare.py /path/to/infx1/file.if > offers.sql`

`sqlite3 offers.db < offers.sql`
package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"tPaxxOffers/models"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"database/sql"
	"tPaxxOffers/restapi/operations"
	"tPaxxOffers/restapi/operations/offers"

	// import to provide sqlite connector
	_ "github.com/mattn/go-sqlite3"
)

//go:generate swagger generate server --target .. --name tPaxxOffers --spec ../swagger.yaml

func configureFlags(api *operations.TPaxxOffersAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func allOffers(params offers.FindOffersParams) (result []*models.Offer) {
	result = make([]*models.Offer, 0)
	var argCounter = 1

	var ANDS []string
	var ANDPARAMS []interface{}

	if params.Code != nil {
		ANDS = append(ANDS, fmt.Sprintf("hotelcode = ?%d", argCounter))
		ANDPARAMS = append(ANDPARAMS, *params.Code)
		argCounter++
	} else if params.Destination != nil {
		ANDS = append(ANDS, fmt.Sprintf("outdest = ?%d", argCounter))
		ANDPARAMS = append(ANDPARAMS, *params.Destination)
		argCounter++
	}

	if params.Source != nil {
		ANDS = append(ANDS, fmt.Sprintf("outsource = ?%d", argCounter))
		ANDPARAMS = append(ANDPARAMS, *params.Source)
		argCounter++
	}

	if params.Passengers != nil {
		ANDS = append(ANDS, fmt.Sprintf("belegung = ?%d", argCounter))
		ANDPARAMS = append(ANDPARAMS, *params.Passengers)
		argCounter++
	}

	if params.Duration != nil {
		ANDS = append(ANDS, fmt.Sprintf("duration = ?%d", argCounter))
		ANDPARAMS = append(ANDPARAMS, *params.Duration)
		argCounter++
	}

	// if params.From != nil {
	// 	ANDS = append(ANDS, fmt.Sprintf(" = ?%d", argCounter))
	// 	ANDPARAMS = append(ANDPARAMS, *params.From)
	// 	argCounter++
	// }

	// if params.To != nil {
	// 	ANDS = append(ANDS, fmt.Sprintf(" = ?%d", argCounter))
	// 	ANDPARAMS = append(ANDPARAMS, *params.To)
	// 	argCounter++
	// }

	db, err := sql.Open("sqlite3", "./offers.db")
	checkErr(err)
	defer db.Close()

	var queryString string
	if argCounter > 1 {
		queryString = "SELECT * FROM offers WHERE " + strings.Join(ANDS, " AND ") + " order by preis LIMIT 10"
	} else {
		queryString = "SELECT * FROM offers order by preis LIMIT 10"
	}

	println(queryString)
	rows, err := db.Query(queryString, ANDPARAMS...)
	checkErr(err)
	defer rows.Close()

	var offerID int32
	var outsource string
	var outdest string
	var start string
	var duration int32
	var hotelcode string
	var accommodation string
	var catering string
	var carrier string
	var operator string
	var category int32
	var tourtype string
	var bmin int32
	var bmax int32
	var vmin int32
	var vmax int32
	var belegung int32
	var amount int32
	var currency string

	for rows.Next() {
		err = rows.Scan(
			&offerID, &outsource, &outdest, &start, &duration, &hotelcode,
			&accommodation, &catering, &carrier, &operator, &category, &tourtype,
			&bmin, &bmax, &vmin, &vmax, &belegung, &amount, &currency)
		checkErr(err)

		offer := models.Offer{ID: &offerID,
			Outsource:     outsource,
			Outdest:       outdest,
			Start:         start,
			Duration:      duration,
			Hotelcode:     hotelcode,
			Accommodation: accommodation,
			Catering:      catering,
			Carrier:       carrier,
			Operator:      operator,
			Category:      category,
			Tourtype:      tourtype,
			Bmin:          bmin,
			Bmax:          bmax,
			VMAX:          vmax,
			VMIN:          vmin,
			Belegung:      belegung,
			Amount:        amount,
			Currency:      currency}

		result = append(result, &offer)
	}

	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func configureAPI(api *operations.TPaxxOffersAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.OffersFindOffersHandler = offers.FindOffersHandlerFunc(func(params offers.FindOffersParams) middleware.Responder {
		return offers.NewFindOffersOK().WithPayload(allOffers(params))
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

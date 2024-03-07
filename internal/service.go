package internal

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"link_click_counter/config"
)

func GetCounterStats(ctx *fasthttp.RequestCtx, dbpool *pgxpool.Pool) {

	code := ctx.UserValue("code")

	var count int
	err := dbpool.QueryRow(ctx, "select count from URL where code = $1", code).Scan(&count)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			logrus.Errorf("URL with code %s not found :%v", code, err)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Errorf("QueryRow failed: %v\n", err)
		return
	}

	jsonCount, err := json.Marshal(count)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Errorf("QueryRow failed: %v\n", err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(jsonCount)

}

func RedirectURL(ctx *fasthttp.RequestCtx, dbpool *pgxpool.Pool) {

	code := ctx.UserValue("code")

	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Errorf("Failed to begin transaction: %v", err)
		return
	}
	defer tx.Rollback(context.Background())

	var url string
	if err := dbpool.QueryRow(context.Background(), "select url from URL where code = $1", code).Scan(&url); err != nil {
		if url == "" {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			logrus.Errorf("URL with code %s not found :%v", code, err)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		logrus.Errorf("Failed to get URL from database: %v", err)
		return
	}
	_, err = dbpool.Exec(context.Background(), "update URL set count = count + 1 where code = $1", code)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Errorf("Failed to update counter: %v", err)
		return
	}

	_, err = dbpool.Exec(context.Background(), "insert into Statistic (url_id) select URL.id from URL where URL.code = $1", code)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Errorf("Failed to insert into stats: %v", err)
		return
	}

	err = tx.Commit(context.Background())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Errorf("Failed to commit transaction: %v", err)
		return
	}

	ctx.Redirect(url, fasthttp.StatusFound)
}

func GetCounters(ctx *fasthttp.RequestCtx, dbpool *pgxpool.Pool) {
	var err error
	var rows pgx.Rows

	if name := ctx.QueryArgs().Peek("name"); len(name) > 0 {

		nameStr := string(name)

		rows, err = dbpool.Query(ctx, "select * from URL where name = $1", nameStr)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			logrus.Errorf("Query failed: %v\n", err)
			return
		}
	} else {
		rows, err = dbpool.Query(ctx, "select * from URL")
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			logrus.Errorf("Query failed: %v\n", err)
			return
		}
	}

	counters, err := pgx.CollectRows(rows, pgx.RowToStructByName[config.Counter])
	if counters == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		logrus.Error("URL not found")
		return
	}
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Error("QueryRow failed: %v\n", err)
		return
	}

	jsonCounters, err := json.Marshal(counters)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Errorf("QueryRow failed: %v", err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(jsonCounters)
	return
}

func CreateCounter(ctx *fasthttp.RequestCtx, dbpool *pgxpool.Pool) {

	var requestData struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	}
	reqData := &requestData

	if err := json.Unmarshal(ctx.PostBody(), &reqData); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		logrus.Errorf("Failed to parse: %v", err)
		return
	}

	if reqData.URL == "" || reqData.Name == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		logrus.Errorf("URL or name is empty")

		return
	}

	code := uuid.New()

	if _, err := dbpool.Exec(context.Background(), "insert into URL(url, name, code) values($1, $2, $3)", reqData.URL, reqData.Name, code); err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23505" {
			logrus.Errorf("URL %s already exist", reqData.URL)

			ctx.SetStatusCode(fasthttp.StatusConflict)
			return

		}
		if pgErr.Code == "23514" {
			logrus.Errorf("URL %s is not valid", reqData.URL)

			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logrus.Errorf("Failed to insert into database: %v", err)

		return
	}

	idStr := code.String()
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBodyString(idStr)
}

package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/laciferin2024/url-shortner.go/enums"

	"github.com/hiroBzinga/bun"
	"github.com/hiroBzinga/bun/dialect/pgdialect"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

type QueryHook struct {
	log *log.Logger
}

//goland:noinspection ALL
func (h *QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	h.log.Debugf("[BUN] %v\n", event.Query)
	return ctx
}

//goland:noinspection ALL
func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	// h.log.Infof(" %v %v", time.Since(event.StartTime), event.Query) //FIXME doesn't look nice
	h.log.Infof("[BUN] %v\n", time.Since(event.StartTime))
	err := event.Err
	if err != nil {
		h.log.Errorln(err)
	}
}

func newPostgressDB(database *db) (db *bun.DB) {

	conf := database.Conf

	driverName := "postgres"

	dbHost := conf.GetString(enums.POSTGRESQL_HOST)
	dbPort := conf.GetString(enums.POSTGRESQL_PORT)
	dbUser := conf.GetString(enums.POSTGRESQL_USER)
	dbPass := conf.GetString(enums.POSTGRESQL_PASSWORD)
	dbName := conf.GetString(enums.POSTGRESQL_DB)

	dbUrl := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`, dbHost, dbPort, dbUser, dbPass, dbName)

	if sqlDB, err := sql.Open(driverName, dbUrl); err == nil {

		log.Infoln("opened connection to the db")

		err = sqlDB.Ping()

		if err != nil {
			log.Errorln("unable to ping sqlDB")
			//return
			panic(err)
		}

		log.Infoln("pinged the db")

		replicas := getReplicas(conf, driverName)

		dbResolver := newResolver(sqlDB, replicas...)
		dbResolverInst := handleResolver(dbResolver)

		db = bun.NewDB(dbResolverInst, pgdialect.New(), bun.WithDiscardUnknownColumns())

		/*dbResolver.Exec(fmt.Sprintf("set timezone to '%s'", conf.GetString(enums.TIMEZONE)))
		dbResolver.Exec("select 1")*/
		// log.Infoln("successfully tried resolvers")

		db.SetConnMaxLifetime(120 * time.Second)
		db.SetMaxIdleConns(20)
		db.SetMaxOpenConns(200)

	} else {
		panic(fmt.Sprintf("Failed to connect to the sqlDB: %s\n", err.Error()))
	}

	if conf.GetString(enums.MODE) == enums.DEVELOPMENT {

		database.Log.SetFormatter(&log.TextFormatter{
			ForceColors: true,
			ForceQuote:  true,
		})

	} else {

		database.Log.SetFormatter(&log.TextFormatter{
			ForceColors: false,
			ForceQuote:  true,
		})

	}

	db.AddQueryHook(&QueryHook{&database.Log})

	_, err := db.ExecContext(context.TODO(), fmt.Sprintf("set timezone to '%s'", conf.GetString(enums.TIMEZONE)))
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect to sqlDB due to - %s", err))
		return
	}
	return
}

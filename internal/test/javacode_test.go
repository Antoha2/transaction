package test

import (
	"fmt"
	"javacode/internal/config"
	"javacode/internal/repository"
	"javacode/internal/service"
	"javacode/pkg/logger"
	"javacode/pkg/logger/sl"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	transport "javacode/internal/transport/http"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

var (
	once          sync.Once
	TransInstance *transport.ApiImpl
	ServInstance  *transport.Service
)

func NewTransSingleton() *transport.ApiImpl {
	once.Do(func() {
		trans := InitTest()
		TransInstance = trans
	})

	return TransInstance
}

func InitTest() *transport.ApiImpl {

	cfg := config.MustLoad()
	slogger := logger.SetupLogger(cfg.Env)
	dbx := MustInitDb(cfg)

	rep := repository.NewRep(cfg, slogger, dbx)
	serv := service.NewServ(cfg, slogger, rep)
	trans := transport.NewApi(cfg, slogger, serv)

	_, err := dbx.DB.Exec("TRUNCATE javacode")
	if err != nil {
		panic(err)
	}
	_, err = dbx.DB.Exec("INSERT INTO javacode (sum) VALUES (1000)")
	if err != nil {
		panic(err)
	}

	return trans
}

func MustInitDb(cfg *config.Config) *sqlx.DB {

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.Dbname,
		cfg.DBConfig.Sslmode,
	)

	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		slog.Warn("failed to parse config", sl.Err(err))
		os.Exit(1)
	}

	// Make connections
	dbx, err := sqlx.Open("pgx", stdlib.RegisterConnConfig(connConfig))
	if err != nil {
		slog.Warn("failed to create connection db", sl.Err(err))
		os.Exit(1)
	}

	err = dbx.Ping()
	if err != nil {
		slog.Warn("error to ping connection pool", sl.Err(err))
		os.Exit(1)
	}
	slog.Info(fmt.Sprintf("Подключение к базе данных на http://127.0.0.1:%v\n", cfg.DBConfig.Port))
	return dbx
}

//test ChangeSum
func TestChangeSum(t *testing.T) {
	tests := []struct {
		name      string
		inputSum  string
		inputRole string
		want      string
		wantErr   bool
	}{

		{"1", `{"sum":500}`, "admin", "{\"sum\":500}\n", false},                 //500
		{"2", `{"sum":500}`, "client", "{\"sum\":1000}\n", false},               //1000
		{"3", `{"sum":1500}`, "admin", `"message": "Insufficient funds"`, true}, //err
		{"4", `{"sum":"1000"}`, "client", `"message": "Bad Request"`, true},     //err
		{"5", `{sum:1000}`, "client", `"message": "Bad Request"`, true},         //err
		{"6", `{sum:1000}`, "user", `"message": "access denied"`, true},         //err
		{"7", `{"sum":"1000"}`, "", `"message": "access denied"`, true},         //err
		{"8", ``, "", ` "message": "Bad Request"`, true},                        //err

	}

	s := NewTransSingleton()
	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.inputSum))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("User-Role", tc.inputRole)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := s.ChangingSumHandler(c)

			if (err != nil) != tc.wantErr { // если ошибка не нил , и не ждем ошибку
				t.Fatalf("error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if (err != nil) && tc.wantErr { // если ошибка не нил , и ждем ошибку
				return
			}
			if !reflect.DeepEqual(rec.Body.String(), tc.want) { //если нет ошибки , то сравниваем значения
				t.Fatalf("expected: %v, got: %v", rec.Body.String(), tc.want)
			}
		})
	}
}

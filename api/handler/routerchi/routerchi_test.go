package routerchi

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"

	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/handler"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/link"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/pages"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/db/postgres"
)

func clearDB() (*pgx.Conn, error) {
	UrlExample := postgres.UrlExample
	conn, err := pgx.Connect(context.Background(), UrlExample)
	if err != nil {
		return nil, err
	}

	defer conn.Close(context.Background())

	sqlStatement := "truncate links, stat;"

	_, err = conn.Exec(context.Background(), sqlStatement)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getShortID() (int, error) {
	UrlExample := postgres.UrlExample
	conn, err := pgx.Connect(context.Background(), UrlExample)
	if err != nil {
		return 0, err
	}

	defer conn.Close(context.Background())

	sqlStatement := "select (shorturl) from links;"

	shorturl := 0

	err = conn.QueryRow(context.Background(), sqlStatement).Scan(&shorturl)
	if err != nil {
		return 0, fmt.Errorf("error while createShortLink, err: %w", err)
	}

	_, err = conn.Exec(context.Background(), sqlStatement)
	if err != nil {
		return 0, err
	}

	return shorturl, nil
}

func TestGetDefaultPage(t *testing.T) {
	rr := httptest.NewRecorder()

	db, err := postgres.NewPgDB()
	if err != nil {
		t.Error("db connect error: ", err)
		return
	}

	ls := link.NewLinks(db)
	h := handler.NewHandlers(ls)

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("creating zap error: ", err)
		return
	}
	defer logger.Sync()
	log := logger.Sugar()

	rh := NewRouterChi(h, log)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rh.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedBody := fmt.Sprintf(pages.DefaultPageContent, "", "", "")

	if rr.Body.String() != expectedBody {
		t.Error("returned body is wrong")
	}

}

func TestPostLongLink(t *testing.T) {
	if _, err := clearDB(); err != nil {
		t.Error("clearDB error: ", err)
	}

	rr := httptest.NewRecorder()

	db, err := postgres.NewPgDB()
	if err != nil {
		t.Error("db connect error: ", err)
		return
	}

	ls := link.NewLinks(db)
	h := handler.NewHandlers(ls)

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("creating zap error: ", err)
		return
	}
	defer logger.Sync()
	log := logger.Sugar()

	rh := NewRouterChi(h, log)

	postParamValue := "https://github.com/"
	postParams := "lurl=" + postParamValue

	buf := bytes.NewBufferString(postParams)

	req, err := http.NewRequest("POST", "/", buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rh.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	shortID, err := getShortID()
	if err != nil {
		t.Error("error on getShortID: ", err)
	}
	shortIDstr := strconv.Itoa(shortID)

	expectedBody := fmt.Sprintf(pages.DefaultPageContent, postParamValue, pages.ShortLinkUrl+shortIDstr, pages.StatUrl+shortIDstr)

	if rr.Body.String() != expectedBody {
		t.Error("returned body is wrong")
	}
}

func TestGetShortLink(t *testing.T) {
	if _, err := clearDB(); err != nil {
		t.Error("clearDB error: ", err)
	}

	rr := httptest.NewRecorder()

	db, err := postgres.NewPgDB()
	if err != nil {
		t.Error("db connect error: ", err)
		return
	}

	ls := link.NewLinks(db)
	h := handler.NewHandlers(ls)

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("creating zap error: ", err)
		return
	}
	defer logger.Sync()
	log := logger.Sugar()

	rh := NewRouterChi(h, log)

	postParamValue := "https://github.com/"
	postParams := "lurl=" + postParamValue

	buf := bytes.NewBufferString(postParams)

	req, err := http.NewRequest("POST", "/", buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rh.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	shortID, err := getShortID()
	if err != nil {
		t.Error("error on getShortID: ", err)
	}
	shortIDstr := strconv.Itoa(shortID)

	// send get
	req, err = http.NewRequest("GET", "/sl/"+shortIDstr, nil)
	if err != nil {
		t.Fatal(err)
	}

	rh.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetStatLink(t *testing.T) {
	if _, err := clearDB(); err != nil {
		t.Error("clearDB error: ", err)
	}

	rr := httptest.NewRecorder()

	db, err := postgres.NewPgDB()
	if err != nil {
		t.Error("db connect error: ", err)
		return
	}

	ls := link.NewLinks(db)
	h := handler.NewHandlers(ls)

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("creating zap error: ", err)
		return
	}
	defer logger.Sync()
	log := logger.Sugar()

	rh := NewRouterChi(h, log)

	postParamValue := "https://github.com/"
	postParams := "lurl=" + postParamValue

	buf := bytes.NewBufferString(postParams)

	req, err := http.NewRequest("POST", "/", buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rh.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	shortID, err := getShortID()
	if err != nil {
		t.Error("error on getShortID: ", err)
	}
	shortIDstr := strconv.Itoa(shortID)

	// send get to sl
	req, err = http.NewRequest("GET", "/sl/"+shortIDstr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rh.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// send get to stat
	req, err = http.NewRequest("GET", "/stat/"+shortIDstr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rh.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

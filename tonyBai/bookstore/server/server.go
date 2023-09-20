package server

import (
	"bookstore/server/middleware"
	"bookstore/store"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type BookStoreServer struct {
	s   store.Store
	srv *http.Server
}

func NewBookStoreServer(addr string, s store.Store) *BookStoreServer { // 这里是一种设计方法，传入一个接口参数，返回一个具体类型
	bSrv := &BookStoreServer{
		s: s,
		srv: &http.Server{
			Addr: addr,
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/book", bSrv.createBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", bSrv.updateBookHandler).Methods("POST")
	router.HandleFunc("/book/{id}", bSrv.getBookHandler).Methods("GET")
	router.HandleFunc("/book", bSrv.getAllBookHandler).Methods("GET")
	router.HandleFunc("/book/{id}", bSrv.deleteBookHandler).Methods("DELETE")

	// 验证路由和记录路由日志
	bSrv.srv.Handler = middleware.Logging(middleware.Validating(router))

	return bSrv
}

func (bs *BookStoreServer) createBookHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var book store.Book
	if err := dec.Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := bs.s.Create(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (bs *BookStoreServer) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "no id found in request", http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(r.Body)
	var book store.Book
	if err := dec.Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book.Id = id
	if err := bs.s.Update(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func (bs *BookStoreServer) getBookHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "no id fount in request", http.StatusBadRequest)
		return
	}

	book, err := bs.s.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response(w, book)
}

func (bs *BookStoreServer) getAllBookHandler(w http.ResponseWriter, r *http.Request) {
	books, err := bs.s.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response(w, books)
}

func (bs *BookStoreServer) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "no id fount in request", http.StatusBadRequest)
		return
	}

	err := bs.s.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("删除成功"))
}

func response(w http.ResponseWriter, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (bs *BookStoreServer) ListenAndServe() (<-chan error, error) {
	var err error
	errChan := make(chan error)
	go func() {
		err = bs.srv.ListenAndServe()
		errChan <- err
	}()
	select {
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second):
		return errChan, nil
	}
}

func (bs *BookStoreServer) ShutDown(ctx context.Context) error {
	return bs.srv.Shutdown(ctx)
}

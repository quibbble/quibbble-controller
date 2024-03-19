package quibbble_controller

import (
	"io"
	"log"
	"net/http"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.serveMux.ServeHTTP(w, r)
}

func (c *Controller) createHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	snapshot, err := qgn.Parse(string(body))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sanitizeSnapshot(snapshot)
	if err := validateSnapshot(snapshot); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	found, err := c.find(snapshot.Tags[qgn.KeyTag], snapshot.Tags[qgn.IDTag])
	if found {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := c.create(snapshot); err != nil {
		log.Println(err.Error())
		c.delete(snapshot.Tags[qgn.KeyTag], snapshot.Tags[qgn.IDTag])
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(http.StatusText(http.StatusCreated)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

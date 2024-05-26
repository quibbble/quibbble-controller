package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}

func (c *Controller) createHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	snapshot, err := qgn.Parse(string(body))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sanitizeSnapshot(snapshot)
	if err := validateSnapshot(snapshot); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	key, id := snapshot.Tags[qgn.KeyTag], snapshot.Tags[qgn.IDTag]
	if found := c.find(key, id); found {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	// Check long term store to see if a snapshot exists.
	if _, err := c.lookup(key, id); err == nil {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	// If this is a new game then increment game count.
	if err := c.increment(key); err != nil {
		log.Println(err.Error())
	}

	if err := c.create(snapshot); err != nil {
		log.Println(err.Error())
		c.delete(key, id)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(http.StatusText(http.StatusCreated)))
}

func (c *Controller) loadHandler(w http.ResponseWriter, r *http.Request) {
	key, id := r.URL.Query().Get(qgn.KeyTag), r.URL.Query().Get(qgn.IDTag)
	if key == "" {
		http.Error(w, "missing 'key' url query", http.StatusBadRequest)
		return
	} else if id == "" {
		http.Error(w, "missing 'id' url query", http.StatusBadRequest)
		return
	}

	if found := c.find(key, id); found {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
		return
	}

	snapshot, err := c.lookup(key, id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := c.create(snapshot); err != nil {
		log.Println(err.Error())
		c.delete(key, id)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (c *Controller) deleteHandler(w http.ResponseWriter, r *http.Request) {
	key, id := r.URL.Query().Get(qgn.KeyTag), r.URL.Query().Get(qgn.IDTag)
	if key == "" {
		http.Error(w, "missing 'key' url query", http.StatusBadRequest)
		return
	} else if id == "" {
		http.Error(w, "missing 'id' url query", http.StatusBadRequest)
		return
	}

	if found := c.find(key, id); !found {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Store the game into long term storage
	// for future retrieval and play.
	if err := c.store(key, id); err != nil {
		log.Println(err.Error())
	}

	if err := c.delete(key, id); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (c *Controller) activityHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := c.activity()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	raw, err := json.Marshal(stats)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(raw)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

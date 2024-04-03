package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/quibbble/quibbble-controller/games"
	qs "github.com/quibbble/quibbble-controller/internal/server"
	"github.com/quibbble/quibbble-controller/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Stats struct {
	LiveGameCount   map[string]int `json:"live_game_count"`
	LivePlayerCount map[string]int `json:"live_player_count"`
}

func (c *Controller) stats() (*Stats, error) {
	stats := Stats{
		LiveGameCount:   make(map[string]int),
		LivePlayerCount: make(map[string]int),
	}

	for _, builder := range games.Builders {
		key := builder.GetInformation().Key
		stats.LiveGameCount[key] = 0
		stats.LivePlayerCount[key] = 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	l, err := c.clientset.CoreV1().Pods(k8s.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", k8s.Component, k8s.GameComponent),
	})
	if err != nil {
		return nil, err
	}
	names := make([]string, 0)
	for _, it := range l.Items {
		names = append(names, it.Name)
	}
	for _, name := range names {
		url := fmt.Sprintf("http://%s.%s/active", name, k8s.Namespace)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var active qs.Active
		if err := json.Unmarshal(body, &active); err != nil {
			return nil, err
		}
		key, _ := k8s.KeyID(name)
		stats.LiveGameCount[key] += 1
		stats.LivePlayerCount[key] += active.PlayerCount
	}

	return &stats, nil
}

package webui

import (
	"net/http"
	"strconv"
	"time"

	"github.com/huohen92/anifeed/store"
)

type historyItemView struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	ConfigName string    `json:"config_name"`
	ConfigURL  string    `json:"config_url"`
	TaskHash   string    `json:"task_hash"`
	CreatedAt  time.Time `json:"created_at"`
}

type historyResponse struct {
	Items []historyItemView `json:"items"`
	Total int               `json:"total"`
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 && n <= 1000 {
			limit = n
		}
	}
	name := r.URL.Query().Get("name")

	st, err := store.NewWithPath(s.databasePath())
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "打开数据库失败: "+err.Error())
		return
	}
	defer st.Close()

	items, err := st.ListHistory(limit, name)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "读取历史记录失败: "+err.Error())
		return
	}
	resp := historyResponse{Items: make([]historyItemView, 0, len(items)), Total: len(items)}
	for _, item := range items {
		resp.Items = append(resp.Items, historyItemView{
			ID:         item.ID,
			Title:      item.Title,
			ConfigName: item.ConfigName,
			ConfigURL:  item.ConfigUrl,
			TaskHash:   item.TaskHash,
			CreatedAt:  item.CreatedAt,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

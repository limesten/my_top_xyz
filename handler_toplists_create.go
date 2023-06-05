package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
)

type Toplist struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	UserID      int           `json:"user_id"`
	CreatedAt   time.Time     `json:"created_at"`
	Items       []ToplistItem `json:"items"`
}

type ToplistItem struct {
	ID          int    `json:"id"`
	ListId      int    `json:"listId"`
	Rank        int    `json:"rank"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (t Toplist) ToDBToplist() database.Toplist {
	dbItems := make([]database.ToplistItem, len(t.Items))
	for i, item := range t.Items {
		dbItems[i] = item.ToDBToplistItem()
	}

	return database.Toplist{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		UserID:      t.UserID,
		CreatedAt:   t.CreatedAt,
		Items:       dbItems,
	}
}

func (t ToplistItem) ToDBToplistItem() database.ToplistItem {
	return database.ToplistItem{
		ID:          t.ID,
		ListID:      t.ListId,
		Rank:        t.Rank,
		Title:       t.Title,
		Description: t.Description,
	}
}

func (cfg apiConfig) handlerToplistsCreate(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Id int `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	var toplist Toplist
	err := decoder.Decode(&toplist)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	dbToplist := toplist.ToDBToplist()

	insertedToplist, err := cfg.DB.InsertToplist(dbToplist)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Error occurred when creating new toplist")
		return
	}

	respondWithJSON(w, http.StatusCreated, resp{
		Id: insertedToplist.ID,
	})
}

package web

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/n-kurasawa/slack-bot/internal/bot"
	"github.com/n-kurasawa/slack-bot/internal/db"
)

type Handler struct {
	store    *db.Store
	template *template.Template
}

func NewHandler(store *db.Store) (*Handler, error) {
	tmpl, err := template.ParseGlob("internal/web/templates/*.html")
	if err != nil {
		return nil, err
	}

	return &Handler{
		store:    store,
		template: tmpl,
	}, nil
}

func (h *Handler) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.handleIndex)
	mux.HandleFunc("/add", h.handleAdd)
	mux.HandleFunc("/delete", h.handleDelete)

	// 静的ファイル用のハンドラー（必要に応じて）
	// mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/web/static"))))

	return mux
}

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	images, err := h.store.ListImages()
	if err != nil {
		http.Error(w, "画像一覧の取得に失敗しました", http.StatusInternalServerError)
		return
	}

	data := struct {
		Images []bot.Image
	}{
		Images: images,
	}

	if err := h.template.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "テンプレートの実行に失敗しました", http.StatusInternalServerError)
	}
}

func (h *Handler) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "不正なリクエストメソッドです", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームデータの解析に失敗しました", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	url := r.FormValue("url")

	if name == "" || url == "" {
		http.Error(w, "名前とURLは必須です", http.StatusBadRequest)
		return
	}

	if err := h.store.SaveImage(name, url); err != nil {
		http.Error(w, "画像の保存に失敗しました", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "不正なリクエストメソッドです", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォームデータの解析に失敗しました", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "不正なID形式です", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteImage(id); err != nil {
		http.Error(w, "画像の削除に失敗しました", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

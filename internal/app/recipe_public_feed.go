package app

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Public/Social Feed screen recipe (Phase G): a server-rendered, reverse
// chronological activity feed. Each item is a Card composed from the Avatar,
// Badge and Button primitives, filtered by server-side Tabs (?view=), paginated
// with the standalone pagination partial, and backed by POST+303 reactions plus
// a flash toast. The loading state is a documented server-rendered Skeleton
// placeholder (the feed has no client loading phase) and the empty state reuses
// the shared primitive. The recipe introduces no new primitives.

// recipeFeedPageSize controls the server-side pagination slice.
const recipeFeedPageSize = 5

// recipeFeedViews is the closed vocabulary of feed views (server-side Tabs).
var recipeFeedViews = []string{"for-you", "following", "new"}

var recipeFeedViewLabels = map[string]string{
	"for-you":   "For you",
	"following": "Following",
	"new":       "New",
}

// recipeFeedKinds is the closed vocabulary of post kinds.
var recipeFeedKinds = []string{"text", "image", "link", "announcement"}

var recipeFeedKindLabels = map[string]string{
	"text":         "Text",
	"image":        "Image",
	"link":         "Link",
	"announcement": "Announcement",
}

// recipeFeedFollowingAuthors is the closed "following" set driving the
// following view.
var recipeFeedFollowingAuthors = map[string]bool{
	"Alicia R.": true,
	"Carla M.":  true,
}

// recipeFeedItem is one post in the feed (the item IS the event).
type recipeFeedItem struct {
	ID        string
	Author    string
	Kind      string
	Body      string
	Timestamp time.Time
	Likes     int
	Comments  int
	New       bool
}

// recipeFeedStore is the in-memory mock store behind the recipe. A single
// mutex guards the slice and the flash toast; the dataset lives only in server
// memory (like the WhatsApp demo store).
type recipeFeedStore struct {
	mu    sync.Mutex
	seq   int
	items []recipeFeedItem
	toast *toastView
}

// feedDemoStore is the shared demo store for the recipe routes.
var feedDemoStore = newRecipeFeedStore()

// resetRecipeFeedStore restores the demo store to its seed state. Tests use it
// to stay deterministic regardless of execution order.
func resetRecipeFeedStore() {
	feedDemoStore = newRecipeFeedStore()
}

// newRecipeFeedStore seeds the mock dataset with timestamps relative to the
// caller's clock, so the reverse-chronological order is deterministic.
func newRecipeFeedStore() *recipeFeedStore {
	now := time.Now()
	return &recipeFeedStore{
		items: []recipeFeedItem{
			{ID: "post-01", Author: "Alicia R.", Kind: "text", Body: "The new design system guide is live: reusable tokens, server-first workflows, zero component JavaScript.", Timestamp: now.Add(-12 * time.Minute), Likes: 24, Comments: 5, New: true},
			{ID: "post-02", Author: "Dev Ops", Kind: "link", Body: "Rolling the ops queue recipe out to staging today — avatar, tone badges, POST+303 transitions.", Timestamp: now.Add(-45 * time.Minute), Likes: 9, Comments: 2, New: true},
			{ID: "post-03", Author: "Carla M.", Kind: "announcement", Body: "Maintenance window this Saturday 02:00–04:00 UTC. The feed may be briefly unavailable.", Timestamp: now.Add(-3 * time.Hour), Likes: 31, Comments: 8},
			{ID: "post-04", Author: "Bob T.", Kind: "text", Body: "Skeleton states are the first thing I check in a review now — empty flashes before data are a real smell.", Timestamp: now.Add(-7 * time.Hour), Likes: 17, Comments: 4},
			{ID: "post-05", Author: "Alicia R.", Kind: "image", Body: "The avatar primitive shipped with scoped tokens and a proper aria-hidden contract.", Timestamp: now.Add(-5 * time.Hour), Likes: 12, Comments: 1, New: true},
			{ID: "post-06", Author: "Ops Lead", Kind: "link", Body: "Pagination is now a standalone partial — the data table footer and the recipes share it.", Timestamp: now.Add(-1 * time.Hour), Likes: 5, Comments: 0},
		},
	}
}

// snapshot returns a copy of every item so handlers never hold a live reference
// into the store's slice.
func (s *recipeFeedStore) snapshot() []recipeFeedItem {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]recipeFeedItem, len(s.items))
	copy(out, s.items)
	return out
}

func (s *recipeFeedStore) get(id string) (recipeFeedItem, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, it := range s.items {
		if it.ID == id {
			return it, true
		}
	}
	return recipeFeedItem{}, false
}

func (s *recipeFeedStore) addLike(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].Likes++
			return true
		}
	}
	return false
}

func (s *recipeFeedStore) delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, it := range s.items {
		if it.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}
	return false
}

// setToast stores a flash toast rendered inline on the next list render (the
// POST+303 contract: the mutation redirects, the following GET shows the
// transient result and consumes it).
func (s *recipeFeedStore) setToast(t toastView) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.toast = &t
}

func (s *recipeFeedStore) takeToast() *toastView {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.toast == nil {
		return nil
	}
	t := s.toast
	s.toast = nil
	return t
}

// ----- view models -----

// skeletonView is the production view model for the shared "skeleton" partial:
// a server-driven loading placeholder (Avatar draws the avatar+title+lines
// layout, Lines is the slice of placeholder text lines to render). The partial
// renders role="status" plus an sr-only label. Lines is a slice (not a count)
// because Go templates can only range over slices/maps/chans, not counts.
type skeletonView struct {
	Avatar bool
	Label  string
	Lines  []int
}

// recipeFeedView is the feed screen: app shell + flash toast slot + Tabs (view
// selector) + the feed panel (Card list + empty state + standalone pagination)
// + the documented Skeleton loading placeholder + the remote refresh form.
type recipeFeedView struct {
	Meta         metaView
	ThemeClass   string
	Title        string
	Description  string
	ViewValue    string
	Views        []recipeFeedViewOption
	Items        []recipeFeedItemView
	Caption      string
	Pagination   *paginationView
	EmptyState   emptyStateView
	Skeleton     *skeletonView
	FlashToast   *toastView
	Refreshed    bool
	RefreshToast *toastView
}

// recipeFeedViewOption is one server-side Tab (a real link, aria-current on the
// active one).
type recipeFeedViewOption struct {
	Value  string
	Label  string
	Href   string
	Active bool
}

// recipeFeedItemView is one rendered post Card.
type recipeFeedItemView struct {
	ID           string
	Author       string
	Avatar       avatarView
	KindLabel    string
	RelativeTime string
	Body         string
	Likes        int
	Comments     int
	New          bool
	ReactAction  string
}

// ----- handlers -----

func (s *server) recipePublicFeedList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	view := newRecipeFeedView(q.Get("view"), q.Get("page"), feedDemoStore.takeToast())
	applyRequestTheme(r, view)

	if strings.EqualFold(r.Header.Get("HX-Request"), "true") {
		s.renderRecipeTemplate(w, http.StatusOK, "recipe-public-feed-panel", view)
		return
	}
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-public-feed", view)
}

// recipePublicFeedReact records a like as a real POST+303 and flashes a
// transient toast on the following render.
func (s *server) recipePublicFeedReact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, ok := feedDemoStore.get(id)
	if !ok {
		s.recipeFeedNotFound(w, "Post not found", "The post you are trying to react to does not exist or has been removed.")
		return
	}
	feedDemoStore.addLike(id)
	feedDemoStore.setToast(newToast("success", "recipe-pf-react-result", fmt.Sprintf("You liked %s's post.", item.Author)))
	http.Redirect(w, r, "/recipes/public-feed", http.StatusSeeOther)
}

// recipePublicFeedRefresh completes a remote refresh for the feed. With HTMX it
// returns the refresh form fragment and an HX-Trigger raising gelium:toast;
// without JavaScript it re-renders the full page with the feed refreshed and a
// persistent inline toast, exactly like the Data table refresh demo.
func (s *server) recipePublicFeedRefresh(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	isHX := strings.EqualFold(r.Header.Get("HX-Request"), "true")

	view := newRecipeFeedView("", "", feedDemoStore.takeToast())
	applyRequestTheme(r, view)
	view.Refreshed = true

	if isHX {
		trigger, err := toastTriggerJSON("success", "Feed refreshed.")
		if err != nil {
			http.Error(w, "refresh unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Trigger", trigger)
		s.renderRecipeTemplate(w, http.StatusOK, "recipe-public-feed-refresh-form", view)
		return
	}

	toast := newToast("success", "recipe-pf-refresh-result", "Feed refreshed.")
	view.RefreshToast = &toast
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-public-feed", view)
}

func (s *server) recipeFeedNotFound(w http.ResponseWriter, title, body string) {
	s.renderErrorPage(w, http.StatusNotFound, title, body, true, "/recipes/public-feed", "Back to the feed", "/recipes/public-feed")
}

// ----- view builders -----

// newRecipeFeedView validates the request against the closed view vocabulary,
// filters by view, orders reverse-chronologically (novelty is the value),
// paginates and builds the feed view model.
func newRecipeFeedView(viewParam, page string, flash *toastView) *recipeFeedView {
	viewValue := recipeClosedValue(viewParam, recipeFeedViews)
	if viewValue == "" {
		viewValue = "for-you"
	}
	pageNum := recipeParsePage(page)

	now := time.Now()
	snapshot := feedDemoStore.snapshot()
	items := make([]recipeFeedItem, 0, len(snapshot))
	for _, it := range snapshot {
		switch viewValue {
		case "following":
			if !recipeFeedFollowingAuthors[it.Author] {
				continue
			}
		case "new":
			if !it.New {
				continue
			}
		}
		items = append(items, it)
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Timestamp.After(items[j].Timestamp)
	})

	total := len(items)
	totalPages := (total + recipeFeedPageSize - 1) / recipeFeedPageSize
	if totalPages < 1 {
		totalPages = 1
	}
	if pageNum > totalPages {
		pageNum = totalPages
	}

	start := (pageNum - 1) * recipeFeedPageSize
	end := start + recipeFeedPageSize
	if end > total {
		end = total
	}

	rows := make([]recipeFeedItemView, 0, end-start)
	for _, it := range items[start:end] {
		rows = append(rows, recipeFeedRow(it, now))
	}

	var pagination *paginationView
	if totalPages > 1 {
		pagination = newPaginationView(pageNum, totalPages, func(n int) string {
			return recipeFeedHref(viewValue, n)
		})
	}

	views := make([]recipeFeedViewOption, 0, len(recipeFeedViews))
	for _, v := range recipeFeedViews {
		views = append(views, recipeFeedViewOption{
			Value:  v,
			Label:  recipeFeedViewLabels[v],
			Href:   recipeFeedHref(v, 0),
			Active: v == viewValue,
		})
	}

	meta := screenRecipeMeta(
		"Latest activity · Public/Social Feed recipe",
		"Scan the latest activity: a reverse-chronological feed with server-side views, reactions and documented loading and empty states.",
		"/recipes/public-feed",
	)

	return &recipeFeedView{
		Meta:         meta,
		ThemeClass:   themeClass(""),
		Title:        "Latest activity",
		Description:  "The Public/Social Feed screen recipe: a server-rendered activity feed composed from Card, Avatar, Badge, Tabs, Skeleton, Empty state, Button and Toast.",
		ViewValue:    viewValue,
		Views:        views,
		Items:        rows,
		Caption:      fmt.Sprintf("%d posts · page %d of %d", total, pageNum, totalPages),
		Pagination:   pagination,
		EmptyState:   recipeFeedEmptyState(viewValue),
		Skeleton:     &skeletonView{Avatar: true, Label: "Loading the feed", Lines: []int{1, 2, 3}},
		FlashToast:   flash,
	}
}

// recipeFeedRow renders one post: a decorative sm Avatar (paired with the
// visible author name), the kind + relative time meta, the body and the react
// form plus the comment count badge.
func recipeFeedRow(it recipeFeedItem, now time.Time) recipeFeedItemView {
	return recipeFeedItemView{
		ID:           it.ID,
		Author:       it.Author,
		Avatar:       avatarView{Initials: recipeInitials(it.Author), Decorative: true, Size: "sm"},
		KindLabel:    recipeFeedKindLabels[it.Kind],
		RelativeTime: recipeRelativeTime(it.Timestamp, now),
		Body:         it.Body,
		Likes:        it.Likes,
		Comments:     it.Comments,
		New:          it.New,
		ReactAction:  "/recipes/public-feed/" + it.ID + "/react",
	}
}

// recipeFeedEmptyState offers a real CTA per view: back to the full feed from a
// filtered empty view, or a refresh link when the feed itself is empty.
func recipeFeedEmptyState(viewValue string) emptyStateView {
	switch viewValue {
	case "new":
		return emptyStateView{
			Title:    "Nothing new yet",
			Body:     "You are up to date. Check back soon for fresh posts.",
			CTA:      true,
			CTAHref:  "/recipes/public-feed?view=for-you",
			CTALabel: "See all posts",
		}
	case "following":
		return emptyStateView{
			Title:    "No posts from people you follow",
			Body:     "Posts from the people you follow will appear here.",
			CTA:      true,
			CTAHref:  "/recipes/public-feed?view=for-you",
			CTALabel: "See all posts",
		}
	}
	return emptyStateView{
		Title:    "No posts yet",
		Body:     "Be the first to share something with the community.",
		CTA:      true,
		CTAHref:  "/recipes/public-feed",
		CTALabel: "Refresh",
	}
}

// recipeFeedHref builds a real GET link preserving the view and page (stable
// parameter order: view, page).
func recipeFeedHref(view string, page int) string {
	var b strings.Builder
	b.WriteByte('?')
	if view != "" {
		b.WriteString("view=")
		b.WriteString(url.QueryEscape(view))
	}
	if page >= 2 {
		if b.Len() > 1 {
			b.WriteByte('&')
		}
		b.WriteString("page=")
		b.WriteString(strconv.Itoa(page))
	}
	return b.String()
}

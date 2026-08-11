package app

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// WhatsApp demo: a functional chat application that dogfoods the Loom
// components with server-side mock data modeled on the Datafy WhatsApp Cloud
// API (https://app.datafyapi.com.br). The main experience is for non-technical
// users (conversation list, chat bubbles, message states, templates presented
// as "Mensajes prediseñados") with zero API jargon; a separate admin section
// behind the 🛡 lock holds the technical surface (numbers + quality rating, QR
// codes, webhooks + HMAC signatures, templates, tokens, rate limits).
//
// The 24-hour customer-service window is the key UX requirement: every
// conversation has a window that resets on the customer's last inbound message
// (free-form replies allowed inside; only pre-made messages after it expires).
// The UI shows the remaining window time in the conversation list AND in the
// chat header, color-coded (green ok / amber warning < 2h / expired), and the
// composer blocks free text when expired, offering only "Mensajes
// prediseñados".

// whatsAppWindowDuration is the WhatsApp 24h customer-service window.
const whatsAppWindowDuration = 24 * time.Hour

// whatsAppWindowWarningThreshold is when the chip turns amber.
const whatsAppWindowWarningThreshold = 2 * time.Hour

// ----- message kinds (closed vocabulary) -----

// whatsAppMessageKind enumerates the message content kinds.
type whatsAppMessageKind string

const (
	whatsAppKindText     whatsAppMessageKind = "text"
	whatsAppKindImage    whatsAppMessageKind = "image"
	whatsAppKindAudio    whatsAppMessageKind = "audio"
	whatsAppKindList     whatsAppMessageKind = "list"
	whatsAppKindTemplate whatsAppMessageKind = "template"
)

// whatsAppMessageStatus enumerates the delivery states.
type whatsAppMessageStatus string

const (
	whatsAppStatusSent      whatsAppMessageStatus = "sent"
	whatsAppStatusDelivered whatsAppMessageStatus = "delivered"
	whatsAppStatusRead      whatsAppMessageStatus = "read"
)

// ----- view models -----

// whatsAppConversationView is one row in the conversation sidebar.
type whatsAppConversationView struct {
	ID             string
	ContactName    string
	Phone          string
	Initial        string
	LastMessage    string
	LastTime       string
	Unread         int
	WindowLabel    string // e.g. "5h 23m restantes"
	WindowTone     string // "ok" | "warning" | "expired"
	WindowFraction int    // percent 0-100 of window used
	Active         bool
}

// whatsAppMessageView is one chat bubble.
type whatsAppMessageView struct {
	ID        string
	Inbound   bool
	Kind      whatsAppMessageKind
	Body      string
	Time      string
	Status    whatsAppMessageStatus // outbound only
	IsList    bool
	ListTitle string
	ListRows  []string
}

// whatsAppTemplateView is one "Mensaje prediseñado" card.
type whatsAppTemplateView struct {
	Name          string
	FriendlyLabel string
	Category      string
	Body          string
}

// whatsAppChatView is the full chat pane for one conversation.
type whatsAppChatView struct {
	Conversation   whatsAppConversationView
	Messages       []whatsAppMessageView
	WindowExpired  bool
	WindowLabel    string
	WindowTone     string
	WindowFraction int
	Templates      []whatsAppTemplateView
	Typing         bool
}

// whatsAppDemoView is the main app page.
type whatsAppDemoView struct {
	Conversations []whatsAppConversationView
	ActiveChat    *whatsAppChatView
	Search        string
	ThemeClass    string
	Meta          metaView
}

// whatsAppNumberView is one row of the admin numbers table.
type whatsAppNumberView struct {
	DisplayPhone string
	VerifiedName string
	Quality      string // GREEN | YELLOW | RED | UNKNOWN
	Throughput   string
}

// whatsAppQRView is one QR code row.
type whatsAppQRView struct {
	Code        string
	Prefilled   string
	DeepLinkURL string
}

// whatsAppWebhookDeliveryView is one webhook delivery log row.
type whatsAppWebhookDeliveryView struct {
	DeliveryID string
	Type       string
	Status     string // "✅ 200" | "❌ 500"
	HTTP       int
	Timestamp  string
}

// whatsAppAdminView is the admin/dev section.
type whatsAppAdminView struct {
	Numbers             []whatsAppNumberView
	QRCodes             []whatsAppQRView
	Deliveries          []whatsAppWebhookDeliveryView
	WebhookURL          string
	WebhookSecretMasked string
	SignatureOn         bool
	TokenMasked         string
	BaseURL             string
	RateLimitMsg        string
	RateUsed            int // percent
	ActiveTab           string
	ThemeClass          string
	Meta                metaView
}

// ----- mock store -----

type whatsAppMessage struct {
	ID        string
	Inbound   bool
	Kind      whatsAppMessageKind
	Body      string
	Timestamp time.Time
	Status    whatsAppMessageStatus
	ListTitle string
	ListRows  []string
}

type whatsAppConversation struct {
	ID             string
	ContactName    string
	Phone          string
	Unread         int
	WindowDeadline time.Time
	Messages       []whatsAppMessage
}

type whatsAppTemplate struct {
	Name          string
	FriendlyLabel string
	Category      string
	Body          string
}

type whatsAppStore struct {
	mu            sync.Mutex
	Numbers       []whatsAppNumber
	QRCodes       []whatsAppQR
	Deliveries    []whatsAppWebhookDelivery
	Conversations map[string]*whatsAppConversation
	Templates     []whatsAppTemplate
	WebhookURL    string
	Secret        string
	Token         string
	msgSeq        int
	convOrder     []string // ordered conversation IDs, newest first
}

type whatsAppNumber struct {
	DisplayPhone string
	VerifiedName string
	Quality      string
	Throughput   string
}

type whatsAppQR struct {
	Code        string
	Prefilled   string
	DeepLinkURL string
}

type whatsAppWebhookDelivery struct {
	DeliveryID string
	Type       string
	HTTP       int
	Timestamp  time.Time
}

var whatsAppDemoStore = newWhatsAppStore()

func newWhatsAppStore() *whatsAppStore {
	now := time.Now()
	s := &whatsAppStore{
		Numbers: []whatsAppNumber{
			{DisplayPhone: "+55 11 9xxxx-0001", VerifiedName: "LoomChat", Quality: "GREEN", Throughput: "STANDARD"},
			{DisplayPhone: "+55 21 8xxxx-0002", VerifiedName: "LoomChat", Quality: "YELLOW", Throughput: "LOW"},
		},
		QRCodes: []whatsAppQR{
			{Code: "A4O4YGZ", Prefilled: "Olá! Gostaria de saber mais.", DeepLinkURL: "https://wa.me/message/A4O4YGZ"},
			{Code: "B9X2KLM", Prefilled: "Promoção de fim de semana!", DeepLinkURL: "https://wa.me/message/B9X2KLM"},
		},
		Deliveries: []whatsAppWebhookDelivery{
			{DeliveryID: "a3f9...2c", Type: "messages (text)", HTTP: 200, Timestamp: now.Add(-90 * time.Second)},
			{DeliveryID: "7be1...90", Type: "history phase 1", HTTP: 200, Timestamp: now.Add(-35 * time.Minute)},
			{DeliveryID: "21dd...5e", Type: "messages (image)", HTTP: 500, Timestamp: now.Add(-2 * time.Hour)},
		},
		Templates: []whatsAppTemplate{
			{Name: "agendamento", FriendlyLabel: "Confirmación de cita", Category: "Utility", Body: "Hola {{1}}, tu cita es el {{2}} a las {{3}}. ¡Gracias!"},
			{Name: "boas_vindas", FriendlyLabel: "Mensaje de bienvenida", Category: "Utility", Body: "¡Hola {{1}}! Gracias por escribirnos. ¿En qué podemos ayudarte?"},
			{Name: "nota_fiscal", FriendlyLabel: "Boleto disponible", Category: "Utility", Body: "{{1}}, tu boleto ya está disponible: {{2}}"},
			{Name: "promocao", FriendlyLabel: "Promoción del fin de semana", Category: "Marketing", Body: "¡{{1}}! Este fin de semana tenés 20% de descuento. Usá el código {{2}}"},
		},
		WebhookURL: "https://miapp.com/api/whatsapp/webhook",
		Secret:     "whsec_••••••••••••••••••",
		Token:      "sk_live_••••••••••••••••",
	}
	s.Conversations = map[string]*whatsAppConversation{
		"ana": {
			ID: "ana", ContactName: "Ana Souza", Phone: "+55 11 9xxxx-8888",
			WindowDeadline: now.Add(6*time.Hour + 23*time.Minute),
			Messages: []whatsAppMessage{
				{ID: "m1", Inbound: true, Kind: whatsAppKindText, Body: "Olá! Quero agendar uma consulta.", Timestamp: now.Add(-50 * time.Minute), Status: whatsAppStatusRead},
				{ID: "m2", Inbound: false, Kind: whatsAppKindText, Body: "Oi Ana! Claro, temos horários amanhã.", Timestamp: now.Add(-48 * time.Minute), Status: whatsAppStatusRead},
				{ID: "m3", Inbound: false, Kind: whatsAppKindList, Body: "", Timestamp: now.Add(-46 * time.Minute), Status: whatsAppStatusRead, ListTitle: "Escolha um horário", ListRows: []string{"08:00", "09:00", "10:00"}},
				{ID: "m4", Inbound: true, Kind: whatsAppKindText, Body: "09:00, por favor!", Timestamp: now.Add(-40 * time.Minute), Status: whatsAppStatusRead},
				{ID: "m5", Inbound: false, Kind: whatsAppKindText, Body: "Perfecto! Te envío el QR de confirmación.", Timestamp: now.Add(-38 * time.Minute), Status: whatsAppStatusRead},
				{ID: "m6", Inbound: false, Kind: whatsAppKindImage, Body: "qr-consulta.png", Timestamp: now.Add(-37 * time.Minute), Status: whatsAppStatusDelivered},
			},
		},
		"carlos": {
			ID: "carlos", ContactName: "Carlos Lima", Phone: "+55 21 9xxxx-7777",
			WindowDeadline: now.Add(45 * time.Minute),
			Unread:         2,
			Messages: []whatsAppMessage{
				{ID: "m1", Inbound: true, Kind: whatsAppKindText, Body: "Recebi o boleto, obrigado!", Timestamp: now.Add(-2 * time.Hour), Status: whatsAppStatusRead},
				{ID: "m2", Inbound: true, Kind: whatsAppKindText, Body: "Pode me enviar novamente?", Timestamp: now.Add(-40 * time.Minute), Status: whatsAppStatusDelivered},
			},
		},
		"maria": {
			ID: "maria", ContactName: "María Fernanda", Phone: "+55 31 9xxxx-6666",
			WindowDeadline: now.Add(-3 * time.Hour), // expired
			Unread:         3,
			Messages: []whatsAppMessage{
				{ID: "m1", Inbound: true, Kind: whatsAppKindText, Body: "O pedido chegou!", Timestamp: now.Add(-28 * time.Hour), Status: whatsAppStatusRead},
				{ID: "m2", Inbound: false, Kind: whatsAppKindTemplate, Body: "Mensaje de bienvenida", Timestamp: now.Add(-27 * time.Hour), Status: whatsAppStatusRead},
			},
		},
		"grupo": {
			ID: "grupo", ContactName: "Grupo Familia", Phone: "+55 11 9xxxx-5555",
			WindowDeadline: now.Add(20 * time.Hour),
			Unread:         5,
			Messages: []whatsAppMessage{
				{ID: "m1", Inbound: true, Kind: whatsAppKindImage, Body: "foto-praia.jpg", Timestamp: now.Add(-4 * time.Hour), Status: whatsAppStatusRead},
				{ID: "m2", Inbound: true, Kind: whatsAppKindText, Body: "Que linda foto!", Timestamp: now.Add(-3*time.Hour + 50*time.Minute), Status: whatsAppStatusRead},
			},
		},
	}
	s.convOrder = []string{"ana", "carlos", "maria", "grupo"}
	return s
}

// ----- window helpers -----

// windowStatus computes the display label, tone, and used fraction for a
// window deadline. Tone: "ok" while >2h remain, "warning" under 2h, "expired"
// when past.
func windowStatus(deadline time.Time, now time.Time) (label string, tone string, fraction int) {
	remaining := time.Until(deadline)
	if remaining <= 0 {
		return "Ventana vencida", "expired", 100
	}
	used := whatsAppWindowDuration - remaining
	fraction = int(used * 100 / whatsAppWindowDuration)
	if fraction < 0 {
		fraction = 0
	}
	if remaining < whatsAppWindowWarningThreshold {
		return "⚠ " + formatDuration(remaining) + " restantes", "warning", fraction
	}
	return "✓ " + formatDuration(remaining) + " restantes", "ok", fraction
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh %dm", h, m)
}

func initial(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "?"
	}
	return strings.ToUpper(name[:1])
}

func lastMessageText(c *whatsAppConversation) string {
	if len(c.Messages) == 0 {
		return ""
	}
	m := c.Messages[len(c.Messages)-1]
	switch m.Kind {
	case whatsAppKindImage:
		return "📷 " + m.Body
	case whatsAppKindAudio:
		return "🎤 Mensaje de voz"
	case whatsAppKindList:
		return "🗓 Lista: " + m.ListTitle
	case whatsAppKindTemplate:
		return "📄 " + m.Body
	default:
		return m.Body
	}
}

func lastMessageTime(c *whatsAppConversation, now time.Time) string {
	if len(c.Messages) == 0 {
		return ""
	}
	t := c.Messages[len(c.Messages)-1].Timestamp
	if now.Sub(t) < 24*time.Hour && t.Day() == now.Day() {
		return t.Format("15:04")
	}
	return t.Format("02/01")
}

// ----- store operations -----

func (s *whatsAppStore) snapshot(now time.Time) []whatsAppConversationView {
	s.mu.Lock()
	defer s.mu.Unlock()
	views := make([]whatsAppConversationView, 0, len(s.convOrder))
	for _, id := range s.convOrder {
		c := s.Conversations[id]
		label, tone, fraction := windowStatus(c.WindowDeadline, now)
		views = append(views, whatsAppConversationView{
			ID:             c.ID,
			ContactName:    c.ContactName,
			Phone:          c.Phone,
			Initial:        initial(c.ContactName),
			LastMessage:    lastMessageText(c),
			LastTime:       lastMessageTime(c, now),
			Unread:         c.Unread,
			WindowLabel:    label,
			WindowTone:     tone,
			WindowFraction: fraction,
		})
	}
	return views
}

func (s *whatsAppStore) chat(id string, now time.Time) *whatsAppChatView {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.Conversations[id]
	if !ok {
		return nil
	}
	label, tone, fraction := windowStatus(c.WindowDeadline, now)
	messages := make([]whatsAppMessageView, 0, len(c.Messages))
	for _, m := range c.Messages {
		messages = append(messages, whatsAppMessageView{
			ID:        m.ID,
			Inbound:   m.Inbound,
			Kind:      m.Kind,
			Body:      m.Body,
			Time:      m.Timestamp.Format("15:04"),
			Status:    m.Status,
			IsList:    m.Kind == whatsAppKindList,
			ListTitle: m.ListTitle,
			ListRows:  m.ListRows,
		})
	}
	templates := make([]whatsAppTemplateView, 0, len(s.Templates))
	for _, t := range s.Templates {
		templates = append(templates, whatsAppTemplateView{Name: t.Name, FriendlyLabel: t.FriendlyLabel, Category: t.Category, Body: t.Body})
	}
	return &whatsAppChatView{
		Conversation: whatsAppConversationView{
			ID: c.ID, ContactName: c.ContactName, Phone: c.Phone, Initial: initial(c.ContactName),
		},
		Messages:       messages,
		WindowExpired:  tone == "expired",
		WindowLabel:    label,
		WindowTone:     tone,
		WindowFraction: fraction,
		Templates:      templates,
	}
}

func (s *whatsAppStore) addMessage(convID string, inbound bool, kind whatsAppMessageKind, body, listTitle string, listRows []string, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.Conversations[convID]
	if !ok {
		return
	}
	s.msgSeq++
	status := whatsAppStatusDelivered
	if inbound {
		status = whatsAppStatusRead
	}
	c.Messages = append(c.Messages, whatsAppMessage{
		ID: fmt.Sprintf("m%d", s.msgSeq), Inbound: inbound, Kind: kind, Body: body,
		Timestamp: now, Status: status, ListTitle: listTitle, ListRows: listRows,
	})
	if inbound {
		// The 24h window resets on the customer's last inbound message.
		c.WindowDeadline = now.Add(whatsAppWindowDuration)
		c.Unread++
	}
}

func (s *whatsAppStore) markRead(convID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.Conversations[convID]; ok {
		c.Unread = 0
	}
}

// updateWebhook persists the admin webhook settings. The masked secret is a
// demo placeholder, so no real credential ever leaves the page.
func (s *whatsAppStore) updateWebhook(url, secret string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.WebhookURL = url
	s.Secret = secret
}

func (s *whatsAppStore) adminView(now time.Time) whatsAppAdminView {
	s.mu.Lock()
	defer s.mu.Unlock()
	numbers := make([]whatsAppNumberView, 0, len(s.Numbers))
	for _, n := range s.Numbers {
		numbers = append(numbers, whatsAppNumberView{DisplayPhone: n.DisplayPhone, VerifiedName: n.VerifiedName, Quality: n.Quality, Throughput: n.Throughput})
	}
	qrs := make([]whatsAppQRView, 0, len(s.QRCodes))
	for _, q := range s.QRCodes {
		qrs = append(qrs, whatsAppQRView{Code: q.Code, Prefilled: q.Prefilled, DeepLinkURL: q.DeepLinkURL})
	}
	deliveries := make([]whatsAppWebhookDeliveryView, 0, len(s.Deliveries))
	for _, d := range s.Deliveries {
		status := "✅ 200"
		if d.HTTP != 200 {
			status = "❌ " + fmt.Sprint(d.HTTP)
		}
		deliveries = append(deliveries, whatsAppWebhookDeliveryView{
			DeliveryID: d.DeliveryID, Type: d.Type, Status: status, HTTP: d.HTTP, Timestamp: d.Timestamp.Format("15:04:05"),
		})
	}
	sort.Slice(deliveries, func(i, j int) bool { return false })
	return whatsAppAdminView{
		Numbers: numbers, QRCodes: qrs, Deliveries: deliveries,
		WebhookURL: s.WebhookURL, WebhookSecretMasked: s.Secret, SignatureOn: true,
		TokenMasked: s.Token, BaseURL: "https://cloud.datafyapi.com.br/v1/",
		RateLimitMsg: "500 msg/min · 60 media/min · 60 consultas/min",
		RateUsed:     42,
	}
}

// ----- handlers -----

func (s *server) whatsAppDemo(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	search := strings.TrimSpace(r.URL.Query().Get("q"))
	activeID := r.URL.Query().Get("c")
	if activeID == "" {
		activeID = "ana"
	}
	conversations := whatsAppDemoStore.snapshot(now)
	if search != "" {
		lower := strings.ToLower(search)
		filtered := conversations[:0]
		for _, c := range conversations {
			if strings.Contains(strings.ToLower(c.ContactName), lower) || strings.Contains(c.Phone, lower) {
				filtered = append(filtered, c)
			}
		}
		conversations = filtered
		// When the active chat does not match the filter, fall back to the
		// placeholder instead of showing a thread the user cannot see in the
		// list.
		if activeID != "" {
			matched := false
			for _, c := range conversations {
				if c.ID == activeID {
					matched = true
					break
				}
			}
			if !matched {
				activeID = ""
			}
		}
	}
	for i := range conversations {
		conversations[i].Active = conversations[i].ID == activeID
	}
	var chat *whatsAppChatView
	if activeID != "" {
		chat = whatsAppDemoStore.chat(activeID, now)
	}
	s.templates.ExecuteTemplate(w, "demo-whatsapp", whatsAppDemoView{
		Conversations: conversations,
		ActiveChat:    chat,
		Search:        search,
		ThemeClass:    themeClass(themeFromRequest(r)),
		Meta:          demoMetaES,
	})
}

func (s *server) whatsAppAdmin(w http.ResponseWriter, r *http.Request) {
	view := whatsAppDemoStore.adminView(time.Now())
	view.ActiveTab = "numbers"
	view.ThemeClass = themeClass(themeFromRequest(r))
	view.Meta = demoMetaES
	s.templates.ExecuteTemplate(w, "demo-whatsapp-admin", view)
}

// demoMetaES is the shared metadata for both WhatsApp demo screens: Spanish
// content (G2) and demo surfaces never indexed (SEO contract §4).
var demoMetaES = metaView{
	Lang:   "es",
	Robots: "noindex, nofollow",
}

// whatsAppWebhookSave persists the admin webhook form and redirects back with
// the POST+303 pattern, fixing the dead admin form that previously 405'd.
func (s *server) whatsAppWebhookSave(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	whatsAppDemoStore.updateWebhook(
		strings.TrimSpace(r.FormValue("webhook_url")),
		r.FormValue("webhook_secret"),
	)
	http.Redirect(w, r, "/demo/whatsapp/admin", http.StatusSeeOther)
}

func (s *server) whatsAppSend(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	convID := r.FormValue("conversation")
	body := strings.TrimSpace(r.FormValue("message"))
	now := time.Now()
	if convID != "" && body != "" {
		whatsAppDemoStore.addMessage(convID, false, whatsAppKindText, body, "", nil, now)
	}
	http.Redirect(w, r, "/demo/whatsapp?c="+convID, http.StatusSeeOther)
}

func (s *server) whatsAppSendTemplate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	convID := r.FormValue("conversation")
	tplName := r.FormValue("template")
	now := time.Now()
	if convID != "" && tplName != "" {
		whatsAppDemoStore.addMessage(convID, false, whatsAppKindTemplate, tplName, "", nil, now)
	}
	http.Redirect(w, r, "/demo/whatsapp?c="+convID, http.StatusSeeOther)
}

func (s *server) whatsAppTyping(w http.ResponseWriter, r *http.Request) {
	convID := r.URL.Query().Get("c")
	http.Redirect(w, r, "/demo/whatsapp?c="+convID, http.StatusSeeOther)
}

func (s *server) whatsAppRead(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	convID := r.FormValue("conversation")
	whatsAppDemoStore.markRead(convID)
	http.Redirect(w, r, "/demo/whatsapp?c="+convID, http.StatusSeeOther)
}

// Ensure template.HTML is referenced (demo templates may embed trusted SVG).
var _ = template.HTML("")

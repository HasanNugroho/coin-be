package bot

// internal/bot/handler_nlp.go
//
// Menggantikan handleText agar user bisa chat bebas.
// Bot pakai AI untuk mendeteksi intent, lalu eksekusi fungsi yang sesuai.
// State machine lama tetap dipakai untuk alur konfirmasi & pilih pocket/platform.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/HasanNugroho/coin-be/internal/bot/session"
	tele "gopkg.in/telebot.v4"
)

// Intent yang didukung — dibatasi supaya user nggak terlalu bebas
type Intent struct {
	Action    string  `json:"action"`     // "save_transaction" | "get_summary" | "help" | "unknown"
	TxType    string  `json:"tx_type"`    // "income" | "expense"
	Amount    float64 `json:"amount"`     // 0 jika tidak disebutkan
	Note      string  `json:"note"`       // deskripsi transaksi
	Date      string  `json:"date"`       // "yyyy-MM-dd", kosong = hari ini
	ReplyText string  `json:"reply_text"` // respons natural ke user
}

const intentSystemPrompt = `Kamu adalah asisten keuangan pribadi dalam Telegram bot.
Tugasmu adalah memahami pesan user dan mengekstrak intent serta data transaksi keuangan.

Respond ONLY with valid JSON, no markdown, no explanation.

Format JSON yang harus dikembalikan:
{
  "action": "save_transaction" | "get_summary" | "help" | "unknown",
  "tx_type": "income" | "expense",
  "amount": <number, 0 jika tidak ada>,
  "note": "<deskripsi singkat>",
  "date": "<yyyy-MM-dd, kosong string jika tidak disebutkan>",
  "reply_text": "<respons singkat dan natural dalam bahasa Indonesia>"
}

Rules:
- "save_transaction": user ingin mencatat pemasukan atau pengeluaran
- "get_summary": user ingin tahu ringkasan keuangan / saldo / laporan
- "help": user butuh bantuan atau tanya cara pakai bot
- "unknown": pesan tidak relevan dengan keuangan — tolak dengan sopan
- Jangan izinkan action lain di luar daftar di atas
- Jika amount disebutkan dengan kata seperti "ribu", kalikan 1000. "juta" kalikan 1000000.
- tx_type default "expense" jika tidak jelas
- date: isi dari konteks ("kemarin", "tadi pagi", tanggal eksplisit). Jika tidak ada, kosongkan.
- reply_text: selalu dalam bahasa Indonesia, singkat, dan ramah`

// handleTextNLP menggantikan handleText lama untuk pesan biasa (bukan state machine)
func (h *Handler) handleTextNLP(c tele.Context) error {
	sess := h.sessions.GetOrCreate(c.Sender().ID)
	ctx := context.Background()

	// State machine lama tetap berjalan untuk alur multi-step
	switch sess.State {
	case "awaiting_email":
		return h.handleEmailInput(ctx, c, sess)
	case "awaiting_otp":
		return h.handleOTPInput(ctx, c, sess)
	case "awaiting_tx_amount":
		return h.handleTXAmountInput(c, sess)
	case "awaiting_tx_note":
		return h.handleTXNoteInput(ctx, c, sess)
	}

	// Skip jika user belum login
	if sess.UserID.IsZero() {
		return h.handleStart(c)
	}

	// Parse intent pakai AI
	intent, err := h.parseIntent(ctx, c.Text(), sess)
	if err != nil {
		return c.Send("Maaf, aku lagi gangguan. Coba lagi atau pakai menu di bawah.")
	}

	switch intent.Action {
	case "save_transaction":
		return h.handleNLPTransaction(ctx, c, sess, intent)
	case "get_summary":
		c.Send(intent.ReplyText)
		return h.handleSummary(ctx, c, sess)
	case "help":
		return c.Send(intent.ReplyText + "\n\n" + h.helpText())
	default:
		return c.Send(intent.ReplyText)
	}
}

func (h *Handler) handleNLPTransaction(ctx context.Context, c tele.Context, sess *session.UserSession, intent *Intent) error {
	if intent.Amount <= 0 {
		// Amount belum ada, tanya dulu
		sess.State = "awaiting_tx_amount"
		sess.TempData["tx_type"] = intent.TxType
		if intent.Note != "" {
			sess.TempData["tx_note"] = intent.Note
		}
		if intent.Date != "" {
			sess.TempData["tx_date"] = toRFC3339(intent.Date)
		}
		label := "pengeluaran"
		if intent.TxType == "income" {
			label = "pemasukan"
		}
		return c.Send(fmt.Sprintf("%s\n\nBerapa jumlah %s-nya?", intent.ReplyText, label), tele.RemoveKeyboard)
	}

	// Semua data ada, langsung ke pilih pocket
	amtStr := fmt.Sprintf("%.0f", intent.Amount)
	sess.TempData["tx_type"] = intent.TxType
	sess.TempData["tx_amount"] = amtStr
	sess.TempData["tx_note"] = intent.Note
	sess.TempData["pocket_page"] = "0"
	if intent.Date != "" {
		sess.TempData["tx_date"] = toRFC3339(intent.Date)
	}
	sess.State = "awaiting_tx_pocket"

	c.Send(intent.ReplyText)
	return h.showPocketSelection(ctx, c, sess)
}

// func (h *Handler) parseIntent(ctx context.Context, userMsg string, sess *session.UserSession) (*Intent, error) {
// 	if h.svc.config.OpenAIKey == "" {
// 		// Fallback: tidak ada AI key, minta user pakai menu
// 		return &Intent{
// 			Action:    "unknown",
// 			ReplyText: "Silakan gunakan menu di bawah untuk mencatat transaksi.",
// 		}, nil
// 	}

// 	// Tambahkan konteks tanggal hari ini
// 	today := time.Now().Format("2006-01-02")
// 	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
// 	contextNote := fmt.Sprintf("Hari ini: %s. Kemarin: %s.", today, yesterday)

// 	messages := []map[string]string{
// 		{"role": "system", "content": intentSystemPrompt},
// 		{"role": "user", "content": contextNote + "\n\nPesan user: " + userMsg},
// 	}

// 	payload := map[string]interface{}{
// 		"model":       h.svc.config.AIModel,
// 		"messages":    messages,
// 		"max_tokens":  300,
// 		"temperature": 0,
// 	}

// 	jsonPayload, _ := json.Marshal(payload)
// 	req, err := http.NewRequestWithContext(ctx, "POST", h.svc.config.AIHost, bytes.NewBuffer(jsonPayload))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+h.svc.config.OpenAIKey)

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var result map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return nil, err
// 	}

// 	choices, ok := result["choices"].([]interface{})
// 	if !ok || len(choices) == 0 {
// 		return nil, fmt.Errorf("no choices")
// 	}

// 	content := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
// 	content = strings.TrimPrefix(content, "```json")
// 	content = strings.TrimSuffix(content, "```")
// 	content = strings.TrimSpace(content)

// 	var intent Intent
// 	if err := json.Unmarshal([]byte(content), &intent); err != nil {
// 		return nil, err
// 	}

// 	return &intent, nil
// }

func (h *Handler) parseIntent(ctx context.Context, userMsg string, sess *session.UserSession) (*Intent, error) {
	if h.svc.config.OpenAIKey == "" {
		// fallback tanpa AI
		return &Intent{
			Action:    "unknown",
			ReplyText: "Silakan gunakan menu di bawah untuk mencatat transaksi.",
		}, nil
	}

	// Context tanggal singkat
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	contextNote := fmt.Sprintf("Today: %s. Yesterday: %s.", today, yesterday)

	// Prompt bahasa Inggris, minimal
	prompt := fmt.Sprintf(`You are a personal finance assistant for a Telegram bot. 
Extract intent from the user's message below.
Respond ONLY with valid JSON in this format:
{
  "action":"save_transaction"|"get_summary"|"help"|"unknown",
  "tx_type":"income"|"expense",
  "amount":0,
  "note":"",
  "date":"",
  "reply_text":""
}

Rules:
- save_transaction: record income/expense
- get_summary: request report
- help: ask usage guidance
- unknown: not relevant
- amount: convert "ribu" *1000, "juta" *1000000
- tx_type default: expense
- date: use context (e.g., "kemarin") or leave empty
- reply_text: short natural Indonesian sentence

User message (Bahasa Indonesia): "%s"
Context: %s`, userMsg, contextNote)

	messages := []map[string]string{
		{"role": "system", "content": prompt},
	}

	payload := map[string]interface{}{
		"model":       h.svc.config.AIModel,
		"messages":    messages,
		"max_tokens":  150,
		"temperature": 0,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", h.svc.config.AIHost, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.svc.config.OpenAIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("no choices returned by AI")
	}

	content := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// Unmarshal JSON
	var intent Intent
	if err := json.Unmarshal([]byte(content), &intent); err != nil {
		return nil, err
	}

	// Post-process amount: handle "ribu"/"juta" jika AI lupa
	intent.Amount = normalizeAmount(intent.Amount, userMsg)

	// Default tx_type
	if intent.TxType == "" {
		intent.TxType = "expense"
	}

	// Normalize date → RFC3339
	if intent.Date != "" {
		intent.Date = toRFC3339(intent.Date)
	}

	return &intent, nil
}

// normalizeAmount mengubah angka berdasarkan kata "ribu"/"juta" di message
func normalizeAmount(amount float64, msg string) float64 {
	if amount <= 0 {
		return 0
	}
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "ribu") {
		amount *= 1000
	}
	if strings.Contains(lower, "juta") {
		amount *= 1000000
	}
	return amount
}

func (h *Handler) helpText() string {
	return `*Cara pakai bot ini:*
• Ketik bebas, contoh:
  _"tadi beli kopi 25rb"_
  _"kemarin dapat gaji 5 juta"_
  _"liat ringkasan keuangan bulan ini"_

• Atau pakai tombol menu:
  📊 Ringkasan — lihat laporan keuangan
  💰 Transaksi Baru — catat pengeluaran
  📥 Pemasukan — catat pemasukan
  📤 Pengeluaran — catat pengeluaran

• Kirim foto struk untuk scan otomatis 🧾`
}

// toRFC3339 mengubah "yyyy-MM-dd" ke RFC3339 (noon WIB)
func toRFC3339(dateStr string) string {
	t, err := time.ParseInLocation("2006-01-02", dateStr, time.FixedZone("WIB", 7*3600))
	if err != nil {
		return time.Now().Format(time.RFC3339)
	}
	t = t.Add(12 * time.Hour) // set ke tengah hari
	return t.Format(time.RFC3339)
}

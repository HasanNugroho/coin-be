package bot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/HasanNugroho/coin-be/internal/bot/session"
	"github.com/HasanNugroho/coin-be/internal/core/utils"
	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	svc      *TelegramService
	sessions *session.Store
}

func NewHandler(svc *TelegramService, sessions *session.Store) *Handler {
	return &Handler{
		svc:      svc,
		sessions: sessions,
	}
}

func (h *Handler) Register(b *tele.Bot) {
	// Commands
	b.Handle("/start", h.handleStart)
	b.Handle("/menu", h.handleMenu)
	b.Handle("/cancel", h.handleCancel)

	// Menu Buttons
	b.Handle("📊 Ringkasan", h.handleSummaryBtn)
	b.Handle("💰 Transaksi Baru", h.handleNewTXBtn)
	b.Handle("📥 Pemasukan", h.handleIncomeBtn)
	b.Handle("📤 Pengeluaran", h.handleExpenseBtn)

	// Photos
	b.Handle(tele.OnPhoto, h.handlePhoto)

	// Text (Generic State Handling)
	b.Handle(tele.OnText, h.handleText)

	// Callbacks
	b.Handle(tele.OnCallback, h.handleCallback)
}

func (h *Handler) handleStart(c tele.Context) error {
	ctx := context.Background()
	telegramID := c.Sender().ID
	sess := h.sessions.GetOrCreate(telegramID)

	user, err := h.svc.FindUserByTelegramID(ctx, fmt.Sprintf("%d", telegramID))
	if err == nil && user != nil {
		sess.UserID = user.ID
		c.Send(fmt.Sprintf("Selamat datang kembali, *%s*!", user.Name), tele.ModeMarkdown)
		return h.handleMenu(c)
	}

	sess.State = "awaiting_email"
	return c.Send("Halo! Akun kamu belum terhubung. Masukkan *alamat email* yang terdaftar untuk menghubungkan akun:", tele.ModeMarkdown)
}

func (h *Handler) handleMenu(c tele.Context) error {
	keyboard := &tele.ReplyMarkup{
		ReplyKeyboard: [][]tele.ReplyButton{
			{{Text: "📊 Ringkasan"}, {Text: "💰 Transaksi Baru"}},
			{{Text: "📥 Pemasukan"}, {Text: "📤 Pengeluaran"}},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	return c.Send("Mau ngapain?", keyboard)
}

func (h *Handler) handleCancel(c tele.Context) error {
	h.sessions.ClearState(c.Sender().ID)
	c.Send("❌ Aksi dibatalkan.")
	return h.handleMenu(c)
}

func (h *Handler) handleText(c tele.Context) error {
	sess := h.sessions.GetOrCreate(c.Sender().ID)
	ctx := context.Background()

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

	return h.handleMenu(c)
}

func (h *Handler) handleEmailInput(ctx context.Context, c tele.Context, sess *session.UserSession) error {
	email := strings.TrimSpace(strings.ToLower(c.Text()))
	user, err := h.svc.FindUserByEmail(ctx, email)
	if err != nil || user == nil {
		return c.Send("❌ Email tidak terdaftar. Masukkan email yang valid:")
	}

	if err := h.svc.SendOTP(ctx, email, sess.TelegramID); err != nil {
		fmt.Println("❌ Gagal mengirim OTP:", err)
		return c.Send("❌ Gagal mengirim OTP. Silakan coba lagi nanti.")
	}

	sess.TempData["email"] = email
	sess.State = "awaiting_otp"
	return c.Send(fmt.Sprintf("Kode verifikasi telah dikirim ke *%s*. Masukkan kode 6 digit:", email), tele.ModeMarkdown)
}

func (h *Handler) handleOTPInput(ctx context.Context, c tele.Context, sess *session.UserSession) error {
	otpCode := strings.TrimSpace(c.Text())
	email := sess.TempData["email"]

	user, err := h.svc.VerifyOTP(ctx, email, otpCode, sess.TelegramID)
	if err != nil {
		return c.Send("❌ OTP tidak valid, coba lagi:")
	}

	sess.UserID = user.ID
	sess.State = ""
	sess.TempData = make(map[string]string)

	c.Send("✅ Akun berhasil dihubungkan!")
	return h.handleMenu(c)
}

func (h *Handler) handleSummaryBtn(c tele.Context) error {
	sess := h.sessions.GetOrCreate(c.Sender().ID)
	if sess.UserID.IsZero() {
		return h.handleStart(c)
	}
	return h.handleSummary(context.Background(), c, sess)
}

func (h *Handler) handleSummary(ctx context.Context, c tele.Context, sess *session.UserSession) error {
	summary, err := h.svc.GetSummary(ctx, sess.UserID, "1m")
	if err != nil {
		return c.Send("❌ Gagal mengambil data ringkasan.")
	}

	msg := fmt.Sprintf("📊 *Ringkasan Keuangan*\n\n"+
		"💼 *Total Aset:* Rp %.0f\n"+
		"📈 *Pemasukan:* Rp %.0f\n"+
		"📉 *Pengeluaran:* Rp %.0f\n"+
		"💰 *Selisih:* Rp %.0f",
		summary.TotalNetWorth, summary.PeriodIncome, summary.PeriodExpense, summary.PeriodNet)

	return c.Send(msg, tele.ModeMarkdown)
}

func (h *Handler) handleNewTXBtn(c tele.Context) error {
	return h.handleNewTransaction(c, "expense")
}

func (h *Handler) handleIncomeBtn(c tele.Context) error {
	return h.handleNewTransaction(c, "income")
}

func (h *Handler) handleExpenseBtn(c tele.Context) error {
	return h.handleNewTransaction(c, "expense")
}

func (h *Handler) handleNewTransaction(c tele.Context, txType string) error {
	sess := h.sessions.GetOrCreate(c.Sender().ID)
	if sess.UserID.IsZero() {
		return h.handleStart(c)
	}

	typeLabel := "PENGELUARAN"
	if txType == "income" {
		typeLabel = "PEMASUKAN"
	}

	sess.State = "awaiting_tx_amount"
	sess.TempData["tx_type"] = txType
	return c.Send(fmt.Sprintf("Mencatat *%s*. Masukkan jumlahnya:", typeLabel), tele.ModeMarkdown, tele.RemoveKeyboard)
}

func (h *Handler) handleTXAmountInput(c tele.Context, sess *session.UserSession) error {
	amountStr := strings.ReplaceAll(c.Text(), ".", "")
	amountStr = strings.ReplaceAll(amountStr, ",", ".")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		return c.Send("❌ Jumlah tidak valid. Masukkan angka yang benar (contoh: 50000):")
	}

	sess.TempData["tx_amount"] = amountStr
	sess.State = "awaiting_tx_pocket"
	return h.showPocketSelection(context.Background(), c, sess)
}

const pageSize = 5

func (h *Handler) showPocketSelection(ctx context.Context, c tele.Context, sess *session.UserSession) error {
	pockets, err := h.svc.GetPockets(ctx, sess.UserID)
	if err != nil || len(pockets) == 0 {
		return h.showPlatformSelection(ctx, c, sess)
	}

	page, _ := strconv.Atoi(sess.TempData["pocket_page"])
	start := page * pageSize
	end := start + pageSize
	if end > len(pockets) {
		end = len(pockets)
	}

	selector := &tele.ReplyMarkup{}
	var rows []tele.Row

	for _, p := range pockets[start:end] {
		btn := selector.Data(
			fmt.Sprintf("📂 %s (Rp %.0f)", p.Name, utils.Decimal128ToFloat64(p.Balance)),
			"pocket", p.ID.Hex(),
		)
		rows = append(rows, selector.Row(btn))
	}

	// Tombol navigasi
	var navRow tele.Row
	if page > 0 {
		navRow = append(navRow, selector.Data("⬅️ Prev", "pocket_prev"))
	}
	if end < len(pockets) {
		navRow = append(navRow, selector.Data("Next ➡️", "pocket_next"))
	}
	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	selector.Inline(rows...)
	return c.Send(fmt.Sprintf("Pilih *Kantong* (halaman %d):", page+1), selector, tele.ModeMarkdown)
}

func (h *Handler) showPlatformSelection(ctx context.Context, c tele.Context, sess *session.UserSession) error {
	platforms, err := h.svc.GetPlatforms(ctx, sess.UserID)
	if err != nil || len(platforms) == 0 {
		sess.State = "awaiting_tx_note"
		return c.Send("Tambahkan catatan untuk transaksi ini (ketik /skip jika tidak ada):")
	}

	page, _ := strconv.Atoi(sess.TempData["platform_page"])
	start := page * pageSize
	end := start + pageSize
	if end > len(platforms) {
		end = len(platforms)
	}

	selector := &tele.ReplyMarkup{}
	var rows []tele.Row

	for _, p := range platforms[start:end] {
		name := "Platform"
		if p.AliasName != nil {
			name = *p.AliasName
		}
		btn := selector.Data(
			fmt.Sprintf("💳 %s (Rp %.0f)", name, utils.Decimal128ToFloat64(p.Balance)),
			"platform", p.ID.Hex(),
		)
		rows = append(rows, selector.Row(btn))
	}

	// Tombol navigasi
	var navRow tele.Row
	if page > 0 {
		navRow = append(navRow, selector.Data("⬅️ Prev", "platform_prev"))
	}
	if end < len(platforms) {
		navRow = append(navRow, selector.Data("Next ➡️", "platform_next"))
	}
	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	selector.Inline(rows...)
	return c.Send(fmt.Sprintf("Pilih *Platform* (halaman %d):", page+1), selector, tele.ModeMarkdown)
}

func (h *Handler) handleCallback(c tele.Context) error {
	sess := h.sessions.GetOrCreate(c.Sender().ID)
	data := c.Callback().Data
	ctx := context.Background()

	// Pocket dipilih (bukan navigasi)
	if strings.HasPrefix(data, "\fpocket|") {
		parts := strings.Split(data, "|")
		if len(parts) > 1 {
			sess.TempData["tx_pocket_id"] = parts[1]
			sess.TempData["platform_page"] = "0" // reset halaman platform
			sess.State = "awaiting_tx_platform"
			c.Respond()
			return h.showPlatformSelection(ctx, c, sess)
		}
	}

	// Platform dipilih (bukan navigasi)
	if strings.HasPrefix(data, "\fplatform|") {
		parts := strings.Split(data, "|")
		if len(parts) > 1 {
			sess.TempData["tx_platform_id"] = parts[1]
			sess.State = "awaiting_tx_note"
			c.Respond()
			return c.Send("Tambahkan catatan untuk transaksi ini (ketik /skip jika tidak ada):")
		}
	}

	switch data {
	case "\freceipt_save":
		sess.State = "awaiting_tx_pocket"
		sess.TempData["pocket_page"] = "0"
		c.Respond()
		return h.showPocketSelection(ctx, c, sess)

	case "\freceipt_cancel":
		c.Respond()
		c.Send("❌ Scan dibatalkan.")
		h.sessions.ClearState(sess.TelegramID)
		return h.handleMenu(c)

	// Paginasi pocket
	case "\fpocket_next":
		page, _ := strconv.Atoi(sess.TempData["pocket_page"])
		sess.TempData["pocket_page"] = strconv.Itoa(page + 1)
		c.Respond()
		return h.showPocketSelection(ctx, c, sess)
	case "\fpocket_prev":
		page, _ := strconv.Atoi(sess.TempData["pocket_page"])
		if page > 0 {
			sess.TempData["pocket_page"] = strconv.Itoa(page - 1)
		}
		c.Respond()
		return h.showPocketSelection(ctx, c, sess)

	// Paginasi platform
	case "\fplatform_next":
		page, _ := strconv.Atoi(sess.TempData["platform_page"])
		sess.TempData["platform_page"] = strconv.Itoa(page + 1)
		c.Respond()
		return h.showPlatformSelection(ctx, c, sess)
	case "\fplatform_prev":
		page, _ := strconv.Atoi(sess.TempData["platform_page"])
		if page > 0 {
			sess.TempData["platform_page"] = strconv.Itoa(page - 1)
		}
		c.Respond()
		return h.showPlatformSelection(ctx, c, sess)
	}

	return nil
}

func (h *Handler) handleTXNoteInput(ctx context.Context, c tele.Context, sess *session.UserSession) error {
	note := c.Text()
	if note == "/skip" {
		note = ""
	}

	sess.TempData["tx_note"] = note
	return h.submitTransaction(ctx, c, sess)
}

func (h *Handler) submitTransaction(ctx context.Context, c tele.Context, sess *session.UserSession) error {
	amount, _ := strconv.ParseFloat(sess.TempData["tx_amount"], 64)

	// Gunakan tx_date dari scan struk jika ada, fallback ke now
	date := sess.TempData["tx_date"]
	if date == "" {
		date = time.Now().Format(time.RFC3339)
	}

	err := h.svc.CreateTransaction(
		ctx,
		sess.UserID,
		sess.TempData["tx_type"],
		amount,
		sess.TempData["tx_pocket_id"],
		sess.TempData["tx_platform_id"],
		sess.TempData["tx_note"],
		date,
	)

	if err != nil {
		c.Send("❌ Gagal menyimpan transaksi: " + err.Error())
	} else {
		c.Send("✅ Transaksi berhasil disimpan!")
	}

	h.sessions.ClearState(sess.TelegramID)
	return h.handleMenu(c)
}
func (h *Handler) handlePhoto(c tele.Context) error {
	ctx := context.Background()
	sess := h.sessions.GetOrCreate(c.Sender().ID)
	if sess.UserID.IsZero() {
		return h.handleStart(c)
	}

	photo := c.Message().Photo
	c.Send("⏳ Sedang membaca struk...")

	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, photo.FileID+".jpg")

	file := &tele.File{FileID: photo.FileID}
	if err := c.Bot().Download(file, tmpFile); err != nil {
		return c.Send("❌ Gagal mengunduh gambar.")
	}
	defer os.Remove(tmpFile)

	imgData, err := os.ReadFile(tmpFile)
	if err != nil {
		return c.Send("❌ Gagal membaca file gambar.")
	}

	parsed, err := h.svc.ParseReceiptImage(ctx, imgData)
	if err != nil {
		return c.Send("❌ Tidak dapat membaca struk. Coba foto yang lebih jelas.")
	}

	sess.TempData["tx_amount"] = fmt.Sprintf("%.0f", parsed.Amount)
	sess.TempData["tx_description"] = parsed.Description
	sess.TempData["tx_type"] = parsed.Type
	sess.TempData["tx_date"] = parsed.Date
	sess.State = "awaiting_receipt_confirm"

	typeLabel := "Pengeluaran 📤"
	if parsed.Type == "income" {
		typeLabel = "Pemasukan 📥"
	}

	msg := fmt.Sprintf("✅ *Hasil Scan Struk*\n\n"+
		"📋 *Deskripsi:* %s\n"+
		"💵 *Jumlah:* Rp %.0f\n"+
		"🏷️ *Tipe:* %s\n"+
		"📅 *Tanggal:* %s\n\n"+
		"Apakah ingin disimpan?",
		parsed.Description, parsed.Amount, typeLabel, parsed.Date)

	selector := &tele.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("✅ Ya, simpan", "receipt_save"),
			selector.Data("❌ Batal", "receipt_cancel"),
		),
	)

	return c.Send(msg, tele.ModeMarkdown, selector)
}

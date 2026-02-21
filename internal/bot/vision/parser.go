package vision

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
)

// -------------------------
// Struct untuk hasil OCR / Receipt
// -------------------------
type ParsedReceipt struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Type        string  `json:"type"` // income / expense
	Date        string  `json:"date"` // RFC3339
}

// -------------------------
// Queue job
// -------------------------
type OCRJob struct {
	ImageData []byte
	ResultCh  chan *ParsedReceipt
}

// -------------------------
// Receipt Parser
// -------------------------
type ReceiptParser struct {
	OpenAIKey string
	Host      string
	Model     string
}

func NewReceiptParser(openaiKey, host, model string) *ReceiptParser {
	if host == "" {
		host = "https://api.openai.com/v1/chat/completions"
	}
	if model == "" {
		model = "gpt-4o"
	}
	return &ReceiptParser{
		OpenAIKey: openaiKey,
		Host:      host,
		Model:     model,
	}
}

// -------------------------
// Preprocessing: resize + grayscale
// -------------------------
func PreprocessImage(imageData []byte) ([]byte, error) {
	img, err := imaging.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	// Resize max width 1280, keep aspect ratio
	img = imaging.Resize(img, 1280, 0, imaging.Lanczos)
	// Grayscale
	img = imaging.Grayscale(img)

	// Encode kembali ke JPEG []byte
	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, img, imaging.JPEG); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// -------------------------
// OCR + parse langsung di Go
// -------------------------
func (p *ReceiptParser) Parse(ctx context.Context, imageData []byte) (*ParsedReceipt, error) {
	// Preprocess ringan
	preData, err := PreprocessImage(imageData)
	if err != nil {
		return nil, err
	}

	client := gosseract.NewClient()
	defer client.Close()
	client.SetLanguage("ind", "eng")
	client.SetPageSegMode(gosseract.PSM_AUTO) // PSM6 juga bisa dicoba

	if err := client.SetImageFromBytes(preData); err != nil {
		return nil, err
	}

	text, err := client.Text()
	if err != nil {
		return nil, err
	}

	// Jika teks cukup panjang → parse lokal
	if len(text) > 20 {
		res := ParseTextToReceipt(text)
		if res.Amount > 0 {
			return res, nil
		}
	}

	// Optional fallback ke OpenAI / Seed AI
	if p.OpenAIKey != "" {
		return p.parseWithOpenAI(ctx, preData)
	}

	return nil, fmt.Errorf("OCR gagal, teks terlalu pendek")
}

// -------------------------
// Parsing teks struk langsung di Go
// -------------------------
func ParseTextToReceipt(text string) *ParsedReceipt {
	loc := time.FixedZone("WIB", 7*3600)
	res := &ParsedReceipt{
		Type: "expense",
		Date: time.Now().In(loc).Format(time.RFC3339),
	}

	// Extract amount
	amountPatterns := []string{
		"(?i)total[:\\s]+(?:rp\\.?\\s*)?([\\d\\.,]+)",
		"(?i)grand\\s*total[:\\s]+(?:rp\\.?\\s*)?([\\d\\.,]+)",
		"(?i)amount[:\\s]+(?:rp\\.?\\s*)?([\\d\\.,]+)",
		"(?i)subtotal[:\\s]+(?:rp\\.?\\s*)?([\\d\\.,]+)",
		"(?i)rp\\.?\\s*([\\d\\.,]+)",
	}
	for _, pattern := range amountPatterns {
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(text)
		if len(match) > 1 {
			amtStr := match[1]
			amtStr = strings.ReplaceAll(amtStr, ".", "")
			amtStr = strings.ReplaceAll(amtStr, ",", ".")
			amt, err := strconv.ParseFloat(amtStr, 64)
			if err == nil && amt > 0 {
				res.Amount = amt
				break
			}
		}
	}

	// Extract date
	datePatterns := []struct {
		re     string
		layout string
	}{
		{`(\d{4})[/-](\d{2})[/-](\d{2})`, "2006-01-02"},
		{`(\d{2})[/-](\d{2})[/-](\d{4})`, "02-01-2006"},
		{`(\d{2})[/-](\d{2})[/-](\d{2})`, "02-01-06"},
	}
	for _, dp := range datePatterns {
		re := regexp.MustCompile(dp.re)
		match := re.FindString(text)
		if match != "" {
			normalized := strings.ReplaceAll(match, "/", "-")
			t, err := time.ParseInLocation(dp.layout, normalized, loc)
			if err == nil {
				res.Date = t.Format(time.RFC3339)
				break
			}
		}
	}

	// Description: first meaningful line
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 5 && !regexp.MustCompile(`^\d`).MatchString(line) {
			res.Description = line
			break
		}
	}

	return res
}

// -------------------------
// Optional: fallback ke OpenAI/Seed AI
// -------------------------
func (p *ReceiptParser) parseWithOpenAI(ctx context.Context, imageData []byte) (*ParsedReceipt, error) {
	loc := time.FixedZone("WIB", 7*3600)
	now := time.Now().In(loc).Format(time.RFC3339)

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	payload := map[string]interface{}{
		"model": p.Model,
		"messages": []interface{}{
			map[string]interface{}{
				"role": "user",
				"content": []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": fmt.Sprintf(`Analyze this receipt image and extract transaction data.
Current time (RFC3339): %s
Respond ONLY with valid JSON, no explanation:
{"amount": 0, "description": "", "type": "expense", "date": "<RFC3339 format>"}`, now),
					},
					map[string]interface{}{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url": fmt.Sprintf("data:image/jpeg;base64,%s", base64Image),
						},
					},
				},
			},
		},
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", p.Host, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.OpenAIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	choices := result["choices"].([]interface{})
	if len(choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	content := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var parsed ParsedReceipt
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, err
	}

	if _, err := time.Parse(time.RFC3339, parsed.Date); err != nil {
		parsed.Date = time.Now().In(loc).Format(time.RFC3339)
	}

	return &parsed, nil
}

func (p *ReceiptParser) IdentifyCategories(ctx context.Context, note string) ([]string, error) {
	if p.OpenAIKey == "" {
		return nil, fmt.Errorf("OpenAI key is not set")
	}

	payload := map[string]interface{}{
		"model": p.Model,
		"messages": []interface{}{
			map[string]interface{}{
				"role":    "system",
				"content": "You are a finance assistant. Given a transaction note, suggest top 3 most relevant category names in Indonesian. Respond ONLY with a JSON array of strings, for example: [\"Makanan\", \"Transportasi\", \"Belanja\"]",
			},
			map[string]interface{}{
				"role":    "user",
				"content": fmt.Sprintf("Note: %s", note),
			},
		},
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", p.Host, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.OpenAIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	content := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var categories []string
	if err := json.Unmarshal([]byte(content), &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// -------------------------
// Worker Pool
// -------------------------
func StartWorkerPool(numWorkers int, jobs <-chan OCRJob) {
	for i := 0; i < numWorkers; i++ {
		go func(id int) {
			parser := NewReceiptParser("", "", "") // Set OpenAI key jika mau fallback
			for job := range jobs {
				fmt.Printf("Worker %d processing image\n", id)
				ctx := context.Background()
				res, err := parser.Parse(ctx, job.ImageData)
				if err != nil {
					log.Printf("Worker %d error: %v\n", id, err)
					job.ResultCh <- nil
					continue
				}
				job.ResultCh <- res
			}
		}(i + 1)
	}
}

// -------------------------
// Main example
// -------------------------
func main() {
	// Simulasi gambar
	imageBytes, _ := os.ReadFile("input.jpg")

	jobs := make(chan OCRJob, 10)
	StartWorkerPool(2, jobs) // 2 core

	resultCh := make(chan *ParsedReceipt)
	jobs <- OCRJob{ImageData: imageBytes, ResultCh: resultCh}

	parsed := <-resultCh
	if parsed != nil {
		fmt.Printf("Parsed Receipt: %+v\n", parsed)
	} else {
		fmt.Println("Failed to parse receipt")
	}
}

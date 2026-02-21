package vision

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/otiai10/gosseract/v2"
)

type ParsedReceipt struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Type        string  `json:"type"` // "income" or "expense"
	Date        string  `json:"date"` // format: RFC3339
}

type ReceiptParser struct {
	openaiKey string
	host      string
	model     string
}

func NewReceiptParser(openaiKey, host, model string) *ReceiptParser {
	if host == "" {
		host = "https://api.openai.com/v1/chat/completions"
	}
	if model == "" {
		model = "gpt-4o"
	}
	return &ReceiptParser{
		openaiKey: openaiKey,
		host:      host,
		model:     model,
	}
}

func (p *ReceiptParser) Parse(ctx context.Context, imageData []byte) (*ParsedReceipt, error) {
	client := gosseract.NewClient()
	defer client.Close()

	if err := client.SetLanguage("ind", "eng"); err != nil {
		return nil, err
	}
	if err := client.SetImageFromBytes(imageData); err != nil {
		return nil, err
	}

	text, err := client.Text()
	if err != nil {
		return nil, err
	}

	if len(text) > 20 {
		res := p.parseTextToReceipt(text)
		if res.Amount > 0 {
			return res, nil
		}
	}

	return p.parseWithOpenAI(ctx, imageData)
}

func (p *ReceiptParser) parseTextToReceipt(text string) *ParsedReceipt {
	loc := time.FixedZone("WIB", 7*3600)
	res := &ParsedReceipt{
		Type: "expense",
		Date: time.Now().In(loc).Format(time.RFC3339),
	}

	// Amount extraction
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

	// Date extraction — normalize to RFC3339
	type datePattern struct {
		re     string
		layout string
	}
	datePatterns := []datePattern{
		{`(\d{4})[/-](\d{2})[/-](\d{2})`, "2006-01-02"}, // yyyy-MM-dd or yyyy/MM/dd
		{`(\d{2})[/-](\d{2})[/-](\d{4})`, "02-01-2006"}, // dd-MM-yyyy or dd/MM/yyyy
		{`(\d{2})[/-](\d{2})[/-](\d{2})`, "02-01-06"},   // dd-MM-yy
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

func (p *ReceiptParser) parseWithOpenAI(ctx context.Context, imageData []byte) (*ParsedReceipt, error) {
	if p.openaiKey == "" {
		return nil, fmt.Errorf("openai key not configured")
	}

	loc := time.FixedZone("WIB", 7*3600)
	now := time.Now().In(loc).Format(time.RFC3339)

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	payload := map[string]interface{}{
		"model": p.model,
		"messages": []interface{}{
			map[string]interface{}{
				"role": "user",
				"content": []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": fmt.Sprintf(`Analyze this receipt image and extract transaction data.
Current time (RFC3339): %s
Respond ONLY with valid JSON, no explanation, no markdown:
{"amount": 0, "description": "", "type": "expense", "date": "<RFC3339 format, e.g. 2006-01-02T15:04:05+07:00>"}`, now),
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
	req, _ := http.NewRequestWithContext(ctx, "POST", p.host, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.openaiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai api error: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	choices := result["choices"].([]interface{})
	if len(choices) == 0 {
		return nil, fmt.Errorf("no response from openai")
	}

	content := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	// Clean markdown block if present
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var parsed ParsedReceipt
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, err
	}

	// Validate & normalize date — fallback ke now jika AI return format salah
	if _, err := time.Parse(time.RFC3339, parsed.Date); err != nil {
		parsed.Date = time.Now().In(loc).Format(time.RFC3339)
	}

	return &parsed, nil
}

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
	Date        string  `json:"date"` // format: yyyy-MM-dd
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
	res := &ParsedReceipt{
		Type: "expense",
		Date: time.Now().Format("2006-01-02"),
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
			// Clean separators: remove thousand dots, replace decimal comma
			// Simple logic: if there's both . and , , replace . with "" and , with .
			// If only . and it looks like a thousand separator (e.g. 10.000), remove it.
			// This is tricky and varies by locale.
			amtStr = strings.ReplaceAll(amtStr, ".", "")
			amtStr = strings.ReplaceAll(amtStr, ",", ".")
			amt, err := strconv.ParseFloat(amtStr, 64)
			if err == nil {
				res.Amount = amt
				break
			}
		}
	}

	// Date extraction
	datePatterns := []string{
		"(\\d{2})[/-](\\d{2})[/-](\\d{4})",
		"(\\d{4})[/-](\\d{2})[/-](\\d{2})",
	}
	for _, pattern := range datePatterns {
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(text)
		if len(match) > 0 {
			// Just use the first one found for now
			res.Date = match[0]
			break
		}
	}

	// Description: first line with length > 5 chars that doesn't start with a digit
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

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	payload := map[string]interface{}{
		"model": p.model,
		"messages": []interface{}{
			map[string]interface{}{
				"role": "user",
				"content": []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": `Analyze this receipt image and extract transaction data.
Respond ONLY with valid JSON, no explanation, no markdown:
{"amount": 0, "description": "", "type": "expense", "date": "yyyy-MM-dd"}`,
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

	return &parsed, nil
}

package cmd

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/llm-net/adb-claw/pkg/input"
	"github.com/llm-net/adb-claw/pkg/monitor"
	"github.com/spf13/cobra"
)

var (
	liveCartCount int
)

// Product represents a parsed product from the shopping cart.
type Product struct {
	Number int      `json:"number"`
	Title  string   `json:"title"`
	Price  string   `json:"price"`
	Sold   string   `json:"sold,omitempty"`
	Shop   string   `json:"shop,omitempty"`
	Tags   []string `json:"tags,omitempty"`
}

var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "Live stream related commands (Douyin)",
}

var liveCartCmd = &cobra.Command{
	Use:   "cart",
	Short: "Grab products from Douyin live stream shopping cart (小黄车)",
	Long: `Captures product information from a Douyin (抖音) live stream shopping cart.

1. Reads the currently explaining product from the floating card (no cart open needed)
2. Opens the shopping cart (小黄车), scrolls to capture the first N products
3. Closes the cart and outputs structured JSON

This command is designed specifically for Douyin live streams. It uses the Android
accessibility framework to read product data — no screenshots or OCR involved.

Must be used while viewing a Douyin live stream with a shopping cart.

Examples:
  adb-claw live cart              # top 10 products + explaining product
  adb-claw live cart --count 5    # top 5 products`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLiveCart()
	},
}

func init() {
	liveCartCmd.Flags().IntVar(&liveCartCount, "count", 10, "Number of products to capture from the top of the list")

	liveCmd.AddCommand(liveCartCmd)
	rootCmd.AddCommand(liveCmd)
}

func runLiveCart() error {
	start := time.Now()

	// 1. Get screen size
	screenW, screenH, err := input.GetScreenSize(client)
	if err != nil {
		writer.Fail("live cart", "SCREEN_SIZE_FAILED", err.Error(),
			"Ensure device is connected", start)
		return nil
	}
	writer.Verbose("live cart: screen %dx%d, target count=%d", screenW, screenH, liveCartCount)

	// 2. Push monitor DEX
	if err := monitor.EnsureDEX(client); err != nil {
		writer.Fail("live cart", "DEX_PUSH_FAILED", err.Error(),
			"Check device connection", start)
		return nil
	}

	// 3. Capture explaining product (before opening cart)
	writer.Verbose("live cart: capturing explaining product...")
	explaining := captureExplaining()

	// 4. Open shopping cart — tap the "商品列表" button in Douyin's bottom toolbar
	writer.Verbose("live cart: opening shopping cart...")
	cartX := screenW * 59 / 100
	cartY := screenH * 925 / 1000
	if err := input.Tap(client, cartX, cartY); err != nil {
		writer.Fail("live cart", "TAP_FAILED", err.Error(),
			"Could not tap shopping cart button", start)
		return nil
	}
	time.Sleep(2 * time.Second)

	// 5. Scroll & capture until we have 1..N consecutively
	writer.Verbose("live cart: capturing products...")
	allTexts := collectCartTexts(screenW, screenH)

	// 6. Close cart
	writer.Verbose("live cart: closing cart...")
	input.KeyEvent(client, "BACK")

	// 7. Parse products from accessibility text
	products := parseCartProducts(allTexts)

	// 8. Output
	data := map[string]interface{}{
		"products": products,
		"total":    len(products),
	}
	if explaining != nil {
		data["explaining"] = explaining
	}

	writer.Success("live cart", data, start)
	return nil
}

// ---------------------------------------------------------------------------
// Explaining product (floating card, no cart open)
// ---------------------------------------------------------------------------

// captureExplaining polls accessibility to find the "讲解中" product card.
func captureExplaining() map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	proc, err := monitor.Start(ctx, client, 1500, 2)
	if err != nil {
		writer.Verbose("live cart: explaining capture failed: %v", err)
		return nil
	}
	defer proc.Stop()

	var texts []string
	for line := range proc.Lines() {
		if entry, err := monitor.ParseLine(line); err == nil {
			texts = append(texts, entry.Text)
		}
	}
	proc.Wait()

	// Pattern: "热卖NNN主播讲解中【title】...券后价XX.X元购买"
	explainRe := regexp.MustCompile(`热卖(\d+).*?讲解中(.+?)(?:券后价|$)`)
	priceRe := regexp.MustCompile(`券后价\s*(?:¥\s*)?(\d+\.?\d*)`)

	for _, t := range texts {
		if m := explainRe.FindStringSubmatch(t); m != nil {
			titlePart := strings.TrimSpace(m[2])
			titlePart = regexp.MustCompile(`\d+\.?\d*元.*$`).ReplaceAllString(titlePart, "")
			result := map[string]interface{}{
				"hot_sales": m[1],
				"title":     strings.TrimSpace(titlePart),
			}
			if pm := priceRe.FindStringSubmatch(t); pm != nil {
				result["price"] = "¥ " + pm[1]
			}
			return result
		}
	}

	// Fallback: separate text nodes
	for i, t := range texts {
		if strings.Contains(t, "讲解中") {
			result := map[string]interface{}{"status": "讲解中"}
			for j := max(0, i-3); j < min(len(texts), i+5); j++ {
				if strings.HasPrefix(texts[j], "【") {
					result["title"] = texts[j]
				}
				if strings.Contains(texts[j], "¥") {
					result["price"] = texts[j]
				}
			}
			return result
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Cart scrolling & collection
// ---------------------------------------------------------------------------

// collectCartTexts streams accessibility data while slowly scrolling the cart.
// Uses small scroll steps and long waits for high accuracy on the first N products.
func collectCartTexts(screenW, screenH int) []monitor.TextEntry {
	// Small scroll step (~15% of screen) to avoid skipping products
	scrollStartY := screenH * 68 / 100
	scrollEndY := screenH * 52 / 100
	scrollX := screenW / 2

	maxScrolls := liveCartCount * 2 // generous upper bound
	if maxScrolls < 8 {
		maxScrolls = 8
	}

	timeout := time.Duration(maxScrolls*1500+20000) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	proc, err := monitor.Start(ctx, client, 300, 0)
	if err != nil {
		writer.Verbose("live cart: monitor start failed: %v", err)
		return nil
	}
	defer proc.Stop()

	var allTexts []monitor.TextEntry
	seenProducts := make(map[int]bool)

	// Initial read — wait for cart panel to be fully rendered
	allTexts = append(allTexts, drainLines(proc, 2*time.Second)...)
	updateSeen(allTexts, seenProducts)

	for scroll := 0; scroll < maxScrolls; scroll++ {
		writer.Verbose("live cart: scroll %d, products: %d/%d", scroll, len(seenProducts), liveCartCount)

		// Check if we have consecutive 1..N
		if hasConsecutive(seenProducts, liveCartCount) {
			break
		}

		// Gentle scroll
		input.Swipe(client, scrollX, scrollStartY, scrollX, scrollEndY, 300)

		// Generous wait for accessibility to reflect new nodes
		newTexts := drainLines(proc, 1200*time.Millisecond)
		allTexts = append(allTexts, newTexts...)
		updateSeen(allTexts, seenProducts)
	}

	proc.Stop()
	return allTexts
}

func updateSeen(texts []monitor.TextEntry, seen map[int]bool) {
	for _, t := range texts {
		if n := extractProductNumber(t.Text); n > 0 {
			seen[n] = true
		}
	}
}

// hasConsecutive returns true if seen contains all numbers 1..n.
func hasConsecutive(seen map[int]bool, n int) bool {
	for i := 1; i <= n; i++ {
		if !seen[i] {
			return false
		}
	}
	return true
}

// drainLines reads entries from the monitor for the given duration.
func drainLines(proc *monitor.Process, d time.Duration) []monitor.TextEntry {
	var entries []monitor.TextEntry
	timer := time.NewTimer(d)
	defer timer.Stop()
	for {
		select {
		case line, ok := <-proc.Lines():
			if !ok {
				return entries
			}
			if entry, err := monitor.ParseLine(line); err == nil {
				entries = append(entries, *entry)
			}
		case <-timer.C:
			return entries
		}
	}
}

// ---------------------------------------------------------------------------
// Product parsing
// ---------------------------------------------------------------------------

var productHeaderRe = regexp.MustCompile(`^(\d+)号商品`)

func extractProductNumber(text string) int {
	m := productHeaderRe.FindStringSubmatch(text)
	if m == nil {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return n
}

// parseCartProducts builds structured products from collected accessibility entries.
// Uses individual TextView nodes (title, price, sold) which are cleaner than the
// concatenated ViewGroup description.
func parseCartProducts(texts []monitor.TextEntry) []Product {
	productMap := make(map[int]*Product)

	priceRe := regexp.MustCompile(`^¥\s*(\d+\.?\d*)\s*(起)?$`)
	soldRe := regexp.MustCompile(`^已售(.+)$`)

	var currentNum int

	for _, t := range texts {
		text := t.Text
		cls := t.Class

		// Detect product header (ViewGroup content_desc "N号商品...")
		if n := extractProductNumber(text); n > 0 {
			currentNum = n
			if _, exists := productMap[n]; !exists {
				productMap[n] = &Product{Number: n}
			}
			continue
		}

		if currentNum == 0 {
			continue
		}
		p, ok := productMap[currentNum]
		if !ok {
			continue
		}

		// Title — first substantial TextView starting with 【 or a long text
		if p.Title == "" && strings.Contains(cls, "TextView") {
			if strings.HasPrefix(text, "【") || strings.HasPrefix(text, " 【") {
				p.Title = strings.TrimSpace(text)
				continue
			}
			// Some titles don't start with 【 (e.g. brand name first)
			if len([]rune(text)) > 10 && !strings.HasPrefix(text, "¥") &&
				!strings.HasPrefix(text, "已售") && !strings.HasPrefix(text, "来自") &&
				!strings.HasPrefix(text, "平台") && !strings.HasPrefix(text, "x ") {
				p.Title = strings.TrimSpace(text)
				continue
			}
		}

		// Price
		if m := priceRe.FindStringSubmatch(text); m != nil {
			p.Price = "¥ " + m[1]
			if m[2] == "起" {
				p.Price += " 起"
			}
			continue
		}

		// Sold count
		if m := soldRe.FindStringSubmatch(text); m != nil {
			p.Sold = m[1]
			continue
		}

		// Shop
		if strings.HasPrefix(text, "来自") {
			p.Shop = strings.TrimPrefix(text, "来自")
			continue
		}

		// Tags
		switch text {
		case "品牌低价", "运费险", "7天无理由退货", "晚发即赔", "今日发", "优质评价多":
			p.Tags = append(p.Tags, text)
		}
	}

	// Collect 1..count, fill any missing titles from ViewGroup fallback
	fillMissingTitles(productMap, texts)

	var products []Product
	for i := 1; i <= liveCartCount; i++ {
		if p, ok := productMap[i]; ok {
			products = append(products, *p)
		}
	}
	return products
}

// fillMissingTitles fills products that have no title from the ViewGroup description.
func fillMissingTitles(productMap map[int]*Product, texts []monitor.TextEntry) {
	titleFromHeaderRe := regexp.MustCompile(`^\d+号商品(.+?)(?:\d+\.?\d*元(?:起)?)?$`)

	for _, t := range texts {
		n := extractProductNumber(t.Text)
		if n == 0 {
			continue
		}
		p, ok := productMap[n]
		if !ok || p.Title != "" {
			continue
		}
		// Extract title from concatenated ViewGroup text as fallback
		if m := titleFromHeaderRe.FindStringSubmatch(t.Text); m != nil {
			raw := m[1]
			// Strip trailing marketing noise
			p.Title = cleanTitle(raw)
		}
	}
}

// cleanTitle removes common trailing marketing text patterns from a product title.
func cleanTitle(s string) string {
	s = strings.TrimSpace(s)
	// Patterns that mark the start of marketing suffix
	noise := []string{
		"万+个", "万+人", "万人", "千+人", "人评价", "人收藏", "人加购", "人正在看",
		"近期店铺", "近7天", "本店同款", "商品回头客", "热搜度超",
		"带图评价", "优质评价", "%速食", "%同类",
	}
	for _, n := range noise {
		// Find the last occurrence of the noise pattern after at least 10 chars of real title
		idx := strings.Index(s, n)
		if idx > 10 {
			s = strings.TrimSpace(s[:idx])
			break
		}
	}
	// Remove trailing price pattern
	s = regexp.MustCompile(`\d+\.?\d*元(?:起)?$`).ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}

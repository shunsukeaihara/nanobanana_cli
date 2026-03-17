package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/genai"
)

var (
	resolution string
	aspect     string
	outputDir  string
	model      string
	references refFlags
)

type refFlags []string

func (r *refFlags) String() string { return strings.Join(*r, ",") }
func (r *refFlags) Set(v string) error {
	*r = append(*r, v)
	return nil
}

func main() {
	flag.StringVar(&resolution, "resolution", "2K", "Output resolution (1K, 2K, 4K)")
	flag.StringVar(&aspect, "aspect", "16:9", "Aspect ratio (e.g. 16:9, 1:1, 9:16)")
	flag.StringVar(&outputDir, "output", "./generated_images", "Output directory")
	flag.StringVar(&model, "model", "gemini-3-pro-image-preview", "Gemini model name")
	flag.Var(&references, "reference", "Reference image path (can be specified multiple times, max 14)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: nanobanana \"prompt\" [flags]")
		flag.PrintDefaults()
		os.Exit(1)
	}
	prompt := flag.Arg(0)

	if len(references) > 14 {
		fmt.Fprintln(os.Stderr, "[Error] Reference images are limited to 14 max")
		os.Exit(1)
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "[Error] GEMINI_API_KEY environment variable is not set")
		os.Exit(1)
	}

	if err := run(prompt, apiKey); err != nil {
		fmt.Fprintf(os.Stderr, "[Error] %v\n", err)
		os.Exit(1)
	}
}

func run(prompt, apiKey string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Build contents
	var parts []*genai.Part

	// Add reference images
	for _, refPath := range references {
		data, err := os.ReadFile(refPath)
		if err != nil {
			return fmt.Errorf("failed to read reference image %s: %w", refPath, err)
		}
		mimeType := detectMIME(refPath, data)
		parts = append(parts, genai.NewPartFromBytes(data, mimeType))
	}

	// Add prompt text
	parts = append(parts, genai.NewPartFromText(prompt))

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, "user"),
	}

	// Generate
	resp, err := client.Models.GenerateContent(ctx, model, contents, &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
		ImageConfig: &genai.ImageConfig{
			AspectRatio: aspect,
			ImageSize:   resolution,
		},
	})
	if err != nil {
		return fmt.Errorf("image generation failed: %w", err)
	}

	// Process response
	saved := false
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.Text != "" {
				fmt.Printf("[Info] %s\n", part.Text)
			}
			if part.InlineData != nil {
				outPath, err := savePart(part)
				if err != nil {
					return fmt.Errorf("failed to save image: %w", err)
				}
				fmt.Printf("[Success] Image saved: %s\n", outPath)
				saved = true
			}
		}
	}

	if !saved {
		fmt.Println("[Warning] No image was included in the response")
	}
	return nil
}

func savePart(part *genai.Part) (string, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", err
	}

	mimeType := part.InlineData.MIMEType
	ext := mimeToExt(mimeType)
	timestamp := time.Now().Format("20060102_150405")
	filename := timestamp + ext
	outPath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(outPath, part.InlineData.Data, 0o644); err != nil {
		return "", err
	}
	return outPath, nil
}

func mimeToExt(mime string) string {
	switch mime {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/webp":
		return ".webp"
	default:
		return ".png"
	}
}

func detectMIME(path string, data []byte) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	default:
		ct := http.DetectContentType(data)
		if strings.HasPrefix(ct, "image/") {
			return ct
		}
		return "image/png"
	}
}

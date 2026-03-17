---
name: nano-banana-pro
description: |
  Generate and edit images using Gemini's image generation model via the `nanobanana` CLI.
  Use this skill whenever the user asks to create, generate, edit, or modify images, illustrations,
  icons, logos, thumbnails, posters, or any visual content.
---

# Nano Banana Pro Image Generation

Generate and edit images via `nanobanana` CLI, which calls the Gemini API.

## Prerequisites

- `nanobanana` command in PATH
- `GEMINI_API_KEY` environment variable set

## Command

```
nanobanana "prompt" [flags]
```

| Flag | Values | Default | Description |
|---|---|---|---|
| `--resolution` | `1K`, `2K`, `4K` | `2K` | Output resolution |
| `--aspect` | `16:9`, `1:1`, `9:16`, `4:3`, etc. | `16:9` | Aspect ratio |
| `--output` | directory path | `./generated_images` | Output directory |
| `--reference` | image file path | — | Reference image (repeatable, max 14) |
| `--model` | model name | `gemini-3-pro-image-preview` | Gemini model |

Output files are saved with timestamp filenames (e.g. `20251130_153045.png`). Extension is determined by the API response MIME type.

## Prompt Writing

Write prompts in English as natural sentences, not keyword lists. Include: subject, composition, action, location, style, and lighting as relevant.

**Example:**
- Bad: `"cool car, neon, city, night, 8k"`
- Good: `"A cinematic wide shot of a futuristic sports car speeding through rainy Tokyo streets at night, neon reflections on wet asphalt"`

## Reference Images

Use `--reference` to pass existing images for editing, style transfer, or character consistency.

```bash
nanobanana "Change the background to a sunset beach" --reference ./photo.png
nanobanana "Redraw in watercolor style" --reference ./original.jpg
```

For character consistency, explicitly instruct identity preservation in the prompt (e.g. "keep the exact facial features from the reference image").

## Workflow

1. Compose an English prompt from the user's request
2. Choose appropriate `--aspect` and `--resolution` for the use case
3. Add `--reference` if editing or using existing images
4. Run the command
5. Show the generated image to the user using the Read tool

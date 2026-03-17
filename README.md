# nanobanana

A single-binary CLI tool for image generation using Google Gemini API. Go reimplementation of [ccskill-nanobanana](https://github.com/feedtailor/ccskill-nanobanana) — no Python, no venv, just one binary.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/shunsukeaihara/nanobanana_cli/main/install.sh | bash
```

This downloads the binary to `~/.local/bin/nanobanana` and installs the Claude Code skill to `~/.claude/skills/nano-banana-pro/`.

### From source

```bash
go install github.com/shunsukeaihara/nanobanana_cli@latest
```

## Setup

Set your Gemini API key:

```bash
export GEMINI_API_KEY="your-api-key"
```

## Usage

```bash
nanobanana "A cat playing piano in watercolor style"
```

### Options

```
--resolution  1K|2K|4K          Output resolution (default: 2K)
--aspect      16:9|1:1|9:16|... Aspect ratio (default: 16:9)
--output      dir               Output directory (default: ./generated_images)
--reference   file              Reference image (repeatable, max 14)
--model       name              Gemini model (default: gemini-3-pro-image-preview)
```

### Examples

```bash
# High resolution wide image
nanobanana "Sunset coastline" --resolution 4K --aspect 16:9

# Edit an existing image
nanobanana "Change the background to a beach" --reference ./photo.png

# Use multiple reference images
nanobanana "Combine these characters in one scene" \
  --reference ./char_a.png --reference ./char_b.png
```

## Claude Code Skill

This repo includes a Claude Code skill at `skill/SKILL.md`. The install script places it at `~/.claude/skills/nano-banana-pro/SKILL.md` so Claude Code can automatically use it for image generation tasks.

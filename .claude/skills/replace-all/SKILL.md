---
name: bash-replace
description: Text replacement using bash tools (sed, awk, perl). Use when replacing text in files: single file replacement, recursive directory replacement, regex patterns, special character handling, or bulk find-and-replace operations.
compatibility:
  target: claude-code
  shell: bash
  tools: [sed, awk, perl, grep, find]
---

# Bash Text Replacement

## Quick Reference

### Single File Replacement

```bash
sed -i 's/old/new/g' file.txt
```

### Recursive Directory Replacement

```bash
find . -type f -name "*.txt" -exec sed -i 's/old/new/g' {} +
```

### With Backup

```bash
sed -i.bak 's/old/new/g' file.txt
```

## Special Characters

Escape these in patterns: `/ \ & . * [ ] ^ $`

Use alternate delimiter for paths:
```bash
sed -i 's|/old/path|/new/path|g' file.txt
```

## Common Patterns

| Task | Command |
|------|---------|
| Case insensitive | `sed -i 's/old/new/gi' file.txt` |
| First occurrence only | `sed -i 's/old/new/' file.txt` |
| Specific line | `sed -i '5s/old/new/g' file.txt` |
| Line range | `sed -i '5,10s/old/new/g' file.txt` |
| Delete lines matching | `sed -i '/pattern/d' file.txt` |
| Replace whole line | `sed -i 's/.*pattern.*/replacement/' file.txt` |

## Regex Groups

Capture and reuse:
```bash
sed -i 's/\(foo\)\(bar\)/\2\1/g' file.txt  # swaps foobar → barfoo
```

Extended regex (cleaner syntax):
```bash
sed -i -E 's/(foo)(bar)/\2\1/g' file.txt
```

## Multi-file with Confirmation

Preview changes first:
```bash
grep -rl "old" . --include="*.txt" | xargs sed 's/old/new/g'
```

Then apply:
```bash
grep -rl "old" . --include="*.txt" | xargs sed -i 's/old/new/g'
```

## Complex Replacements (Perl)

For advanced regex or multi-line:
```bash
perl -i -pe 's/old/new/g' file.txt
perl -i -0pe 's/old\ntext/new\ntext/g' file.txt  # multiline
```

## Platform Notes

- **macOS**: `sed -i ''` (empty string required)
- **Linux**: `sed -i` (no argument needed)
- **Portable**: Use `sed -i.bak` then remove backups

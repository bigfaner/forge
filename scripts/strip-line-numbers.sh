#!/usr/bin/env bash
# 去除文件开头的行号（支持 N\tN\t内容 或 N\t内容 格式）
# 用法:
#   ./strip-line-numbers.sh file1.md file2.md          # 处理指定文件
#   ./strip-line-numbers.sh plugins/zcode/skills/       # 处理目录下所有 .md
#   find . -name "*.md" | xargs ./strip-line-numbers.sh # 配合 find 使用

set -euo pipefail

strip_file() {
    local file="$1"
    if ! grep -qP '^\d+\t' "$file" 2>/dev/null; then
        echo "skip: $file (no line numbers found)"
        return
    fi
    sed -i 's/^[0-9]*\t//' "$file"
    # 第二遍：处理双重行号（如 1\t2\t内容 -> 去掉第一遍后还剩 2\t内容）
    if grep -qP '^\d+\t' "$file" 2>/dev/null; then
        sed -i 's/^[0-9]*\t//' "$file"
    fi
    echo "done: $file"
}

for target in "$@"; do
    if [ -d "$target" ]; then
        for f in "$target"/**/*.md "$target"/*.md; do
            [ -f "$f" ] && strip_file "$f"
        done
    elif [ -f "$target" ]; then
        strip_file "$target"
    else
        echo "warn: $target not found, skipping"
    fi
done

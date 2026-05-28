#!/usr/bin/env python3
"""Generate baseline artifacts for slim-task-prompt-templates feature.

Creates three baseline files:
1. eval/baseline-token-counts.json - per-file token and line counts
2. eval/frontmatter-baseline.json - frontmatter structure for all 41 templates
3. eval/functional-snapshots/ - per-template functional snapshot checklists
"""

from __future__ import annotations

import json
import os
import re
import sys
from pathlib import Path
from typing import Dict, List, Optional, Tuple, Any

import tiktoken

REPO_ROOT = Path(__file__).resolve().parent.parent

# Directories containing templates
PROMPT_TEMPLATES_DIR = REPO_ROOT / "forge-cli" / "pkg" / "prompt" / "templates"
TASK_TEMPLATES_DIR = REPO_ROOT / "forge-cli" / "pkg" / "task" / "templates"
RECORD_TEMPLATES_DIR = REPO_ROOT / "forge-cli" / "pkg" / "task" / "records"
SKILLS_DIR = REPO_ROOT / "plugins" / "forge" / "skills"
AGENTS_DIR = REPO_ROOT / "plugins" / "forge" / "agents"

# Classification dictionary
CLASSIFICATION = {
    "instruction": "正面指令",
    "constraint": "负面约束",
    "example": "行为示范",
    "format": "格式约定",
    "metadata": "元数据声明",
}

# Tokenizer: use cl100k_base (Claude Sonnet compatible)
ENCODER = tiktoken.get_encoding("cl100k_base")


def count_tokens(text: str) -> int:
    """Count tokens using cl100k_base encoding."""
    return len(ENCODER.encode(text))


def count_lines(text: str) -> int:
    """Count lines in text."""
    return len(text.splitlines())


def parse_frontmatter(content: str) -> Tuple[Optional[Dict], str]:
    """Parse YAML frontmatter from content. Returns (frontmatter_dict, body)."""
    if not content.startswith("---"):
        return None, content

    parts = content.split("---", 2)
    if len(parts) < 3:
        return None, content

    fm_text = parts[1].strip()
    body = parts[2]

    # Simple YAML-like parsing for our frontmatter format
    result = {}
    current_key = None
    current_list = None
    current_map = None

    for line in fm_text.split("\n"):
        stripped = line.strip()
        if not stripped:
            continue

        # Top-level key: value
        m = re.match(r"^(\w+):\s*(.*)", stripped)
        if m and not line.startswith(" "):
            key, value = m.group(1), m.group(2).strip()
            if value:
                if value.startswith("[") and value.endswith("]"):
                    # Inline list
                    items = [v.strip().strip('"').strip("'") for v in value[1:-1].split(",")]
                    result[key] = items
                else:
                    result[key] = value
            else:
                # Could be a map or list
                current_key = key
                current_list = None
                current_map = None
                result[key] = None  # placeholder
            continue

        # List item (  - value)
        if stripped.startswith("- ") and current_key:
            item_val = stripped[2:].strip().strip('"').strip("'")
            if result[current_key] is None:
                result[current_key] = [item_val]
            elif isinstance(result[current_key], list):
                result[current_key].append(item_val)
            continue

        # Map entry (  key: value)
        if line.startswith("  ") and not stripped.startswith("- ") and current_key:
            m2 = re.match(r"^\s*(\w+):\s*(.*)", stripped)
            if m2:
                mk, mv = m2.group(1), m2.group(2).strip()
                if result[current_key] is None:
                    result[current_key] = {}
                if isinstance(result[current_key], dict):
                    result[current_key][mk] = mv if mv else True
                continue

    return result, body


def get_all_prompt_templates() -> list[Path]:
    """Get all 21 prompt template files."""
    return sorted(PROMPT_TEMPLATES_DIR.glob("*.md"))


def get_all_task_templates() -> List[Path]:
    """Get all 14 task template files from forge-cli/pkg/task/templates/."""
    return sorted(TASK_TEMPLATES_DIR.glob("*.md"))


def get_all_record_templates() -> List[Path]:
    """Get all 6 record template files from forge-cli/pkg/task/records/."""
    return sorted(RECORD_TEMPLATES_DIR.glob("*.md"))


def generate_token_counts() -> dict:
    """Generate baseline-token-counts.json."""
    # Content slimming scope: all 21 prompt templates + task-executor.md
    files_to_measure = []

    for p in get_all_prompt_templates():
        files_to_measure.append(("prompt", p))

    task_executor = AGENTS_DIR / "task-executor.md"
    if task_executor.exists():
        files_to_measure.append(("agent", task_executor))

    results = {}
    total_tokens = 0
    total_lines = 0

    for category, fpath in files_to_measure:
        rel_path = str(fpath.relative_to(REPO_ROOT))
        content = fpath.read_text()
        tokens = count_tokens(content)
        lines = count_lines(content)
        total_tokens += tokens
        total_lines += lines

        results[rel_path] = {
            "tokens": tokens,
            "lines": lines,
            "category": category,
        }

    results["_summary"] = {
        "totalTokens": total_tokens,
        "totalLines": total_lines,
        "fileCount": len(files_to_measure),
        "tokenizer": "cl100k_base",
    }

    return results


def generate_frontmatter_baseline() -> dict:
    """Generate frontmatter-baseline.json for all 41 templates."""
    all_templates = []

    # 21 prompt templates
    for p in get_all_prompt_templates():
        all_templates.append(("prompt", p))

    # Task templates
    for p in get_all_task_templates():
        all_templates.append(("task", p))

    # Record templates
    for p in get_all_record_templates():
        all_templates.append(("record", p))

    results = {}

    for template_type, fpath in all_templates:
        rel_path = str(fpath.relative_to(REPO_ROOT))
        content = fpath.read_text()

        fm, body = parse_frontmatter(content)

        if fm is None:
            results[rel_path] = {
                "type": template_type,
                "hasMetadataFrontmatter": False,
                "fields": {},
                "variablesCount": 0,
                "fieldCount": 0,
            }
            continue

        # Analyze frontmatter structure
        fields = {}
        var_count = 0
        for key, value in fm.items():
            if key.startswith("_"):
                continue
            if isinstance(value, list):
                fields[key] = {"type": "list", "count": len(value)}
                if key == "variables":
                    var_count = len(value)
            elif isinstance(value, dict):
                fields[key] = {"type": "map", "count": len(value), "keys": list(value.keys())}
            else:
                fields[key] = {"type": "scalar", "value": str(value)}

        results[rel_path] = {
            "type": template_type,
            "hasMetadataFrontmatter": True,
            "fields": fields,
            "variablesCount": var_count,
            "fieldCount": len([k for k in fm.keys() if not k.startswith("_")]),
        }

    # Summary
    prompt_count = sum(1 for v in results.values() if v["type"] == "prompt")
    task_count = sum(1 for v in results.values() if v["type"] == "task")
    record_count = sum(1 for v in results.values() if v["type"] == "record")

    results["_summary"] = {
        "promptTemplates": prompt_count,
        "taskTemplates": task_count,
        "recordTemplates": record_count,
        "totalTemplates": prompt_count + task_count + record_count,
    }

    return results


def classify_node(text: str) -> str:
    """Classify a semantic node into one of the classification categories."""
    text_lower = text.lower().strip()

    # Negative constraints (否定约束)
    if re.search(r"\b(do not|don't|must not|never|forbidden|avoid|no \|)\b", text_lower):
        return "constraint"

    # Check for constraint patterns like "Failed step | Action |"
    if re.search(r"\|.*fix.*retry", text_lower):
        return "constraint"

    # Format conventions (格式约定)
    if re.match(r"^```", text.strip()):
        return "format"
    if re.match(r"^(Step \d|Output:|`Step)", text.strip()):
        return "format"
    if re.match(r"^\|.*\|.*\|", text.strip()):
        return "format"
    if re.match(r"^#{1,4}\s", text.strip()):
        return "format"

    # Metadata declarations (元数据声明)
    if re.match(r"^(TASK_ID|TASK_FILE|TASK_CATEGORY|SURFACE_KEY):", text.strip()):
        return "metadata"
    if re.match(r"^- \*\*\w+", text.strip()):
        return "metadata"

    # Examples / demonstrations (行为示范)
    if re.search(r"\b(example|e\.g\.|for example|good:|bad:)\b", text_lower):
        return "example"
    if re.search(r"\(.*\.\.\..*\)", text):
        return "example"

    # Positive instructions (正面指令) - default for imperative sentences
    if re.search(r"\b(must|should|need to|ensure|check|read|run|write|create|apply|extract|validate|follow|load|execute|implement|use|populate|confirm)\b", text_lower):
        return "instruction"

    # If it's a list item with content, treat as instruction
    if re.match(r"^[-*]\s", text.strip()):
        return "instruction"

    # Default to instruction for non-empty lines
    if text.strip():
        return "instruction"

    return "format"


def extract_semantic_nodes(content: str) -> list[dict]:
    """Extract semantic nodes from template content."""
    lines = content.split("\n")
    nodes = []
    node_id = 0

    i = 0
    while i < len(lines):
        line = lines[i]
        line_num = i + 1  # 1-based

        # Skip empty lines
        if not line.strip():
            i += 1
            continue

        # Skip frontmatter
        if i == 0 and line.strip() == "---":
            # Find closing ---
            i += 1
            while i < len(lines) and lines[i].strip() != "---":
                i += 1
            i += 1  # skip closing ---
            continue

        # Handle code blocks
        if line.strip().startswith("```"):
            block_lines = [line]
            i += 1
            while i < len(lines) and not lines[i].strip().startswith("```"):
                block_lines.append(lines[i])
                i += 1
            if i < len(lines):
                block_lines.append(lines[i])
                i += 1
            node_id += 1
            block_text = "\n".join(block_lines)
            nodes.append({
                "id": node_id,
                "category": "format",
                "summary": block_text.strip()[:120],
                "sourceLine": line_num,
            })
            continue

        # Handle table blocks (lines starting with |)
        if line.strip().startswith("|"):
            table_lines = [line]
            i += 1
            while i < len(lines) and lines[i].strip().startswith("|"):
                table_lines.append(lines[i])
                i += 1
            node_id += 1
            table_text = "\n".join(table_lines)
            category = classify_node(table_text)
            nodes.append({
                "id": node_id,
                "category": category,
                "summary": table_text.strip()[:120],
                "sourceLine": line_num,
            })
            continue

        # Handle section headers
        if re.match(r"^#{1,4}\s", line):
            node_id += 1
            nodes.append({
                "id": node_id,
                "category": "format",
                "summary": line.strip(),
                "sourceLine": line_num,
            })
            i += 1
            continue

        # Handle CRITICAL/IMPORTANT/EXTREMELY-IMPORTANT blocks as single nodes
        if re.match(r"^<(CRITICAL|IMPORTANT|EXTREMELY-IMPORTANT)>", line.strip()):
            tag = re.match(r"^<(\w[-\w]*)>", line.strip()).group(1)
            block_lines = [line]
            i += 1
            while i < len(lines) and not re.match(rf"^</{tag}>", lines[i].strip()):
                block_lines.append(lines[i])
                i += 1
            if i < len(lines):
                block_lines.append(lines[i])
                i += 1
            node_id += 1
            block_text = "\n".join(block_lines)
            # Classify the block as a whole
            category = classify_node(block_text)
            nodes.append({
                "id": node_id,
                "category": category,
                "summary": f"<{tag}> block ({len(block_lines)} lines)",
                "sourceLine": line_num,
            })
            continue

        # Handle list items (potentially multi-line)
        if re.match(r"^[-*]\s", line):
            item_lines = [line]
            i += 1
            # Collect continuation lines (indented under the list item)
            while i < len(lines) and (lines[i].startswith("  ") or lines[i].startswith("\t")) and lines[i].strip():
                item_lines.append(lines[i])
                i += 1
            node_id += 1
            item_text = "\n".join(item_lines)
            category = classify_node(item_text)
            nodes.append({
                "id": node_id,
                "category": category,
                "summary": item_text.strip()[:120],
                "sourceLine": line_num,
            })
            continue

        # Handle numbered list items
        if re.match(r"^\d+\.\s", line):
            item_lines = [line]
            i += 1
            while i < len(lines) and (lines[i].startswith("  ") or lines[i].startswith("\t")) and lines[i].strip():
                item_lines.append(lines[i])
                i += 1
            node_id += 1
            item_text = "\n".join(item_lines)
            category = classify_node(item_text)
            nodes.append({
                "id": node_id,
                "category": category,
                "summary": item_text.strip()[:120],
                "sourceLine": line_num,
            })
            continue

        # Regular text line
        node_id += 1
        category = classify_node(line)
        nodes.append({
            "id": node_id,
            "category": category,
            "summary": line.strip()[:120],
            "sourceLine": line_num,
        })
        i += 1

    return nodes


def generate_functional_snapshots() -> dict:
    """Generate functional snapshot checklists for content slimming scope files."""
    files_to_snapshot = []

    for p in get_all_prompt_templates():
        files_to_snapshot.append(("prompt", p))

    task_executor = AGENTS_DIR / "task-executor.md"
    if task_executor.exists():
        files_to_snapshot.append(("agent", task_executor))

    results = {}

    for category, fpath in files_to_snapshot:
        rel_path = str(fpath.relative_to(REPO_ROOT))
        content = fpath.read_text()
        nodes = extract_semantic_nodes(content)

        results[rel_path] = {
            "type": category,
            "nodeCount": len(nodes),
            "nodes": nodes,
            "categoryBreakdown": {
                cat: sum(1 for n in nodes if n["category"] == cat)
                for cat in CLASSIFICATION.keys()
            },
        }

    return results


def main():
    eval_dir = REPO_ROOT / "eval"
    eval_dir.mkdir(exist_ok=True)

    # 1. Token counts baseline
    print("Generating baseline-token-counts.json...")
    token_counts = generate_token_counts()
    with open(eval_dir / "baseline-token-counts.json", "w") as f:
        json.dump(token_counts, f, indent=2, ensure_ascii=False)

    # 2. Frontmatter baseline
    print("Generating frontmatter-baseline.json...")
    fm_baseline = generate_frontmatter_baseline()
    with open(eval_dir / "frontmatter-baseline.json", "w") as f:
        json.dump(fm_baseline, f, indent=2, ensure_ascii=False)

    # 3. Functional snapshots
    print("Generating functional snapshots...")
    snapshots_dir = eval_dir / "functional-snapshots"
    snapshots_dir.mkdir(exist_ok=True)

    snapshots = generate_functional_snapshots()
    for rel_path, data in snapshots.items():
        if rel_path.startswith("_"):
            continue
        # Create safe filename
        safe_name = rel_path.replace("/", "_").replace(".", "_")
        with open(snapshots_dir / f"{safe_name}.json", "w") as f:
            json.dump(data, f, indent=2, ensure_ascii=False)

    # Summary
    print("\n=== Baseline Generation Complete ===")
    summary = token_counts.get("_summary", {})
    fm_summary = fm_baseline.get("_summary", {})
    print(f"Token counts: {summary.get('fileCount', 0)} files, {summary.get('totalTokens', 0)} total tokens, {summary.get('totalLines', 0)} total lines")
    print(f"Frontmatter: {fm_summary.get('totalTemplates', 0)} templates ({fm_summary.get('promptTemplates', 0)} prompt + {fm_summary.get('taskTemplates', 0)} task + {fm_summary.get('recordTemplates', 0)} record)")
    print(f"Functional snapshots: {len(snapshots)} files")


if __name__ == "__main__":
    main()

#!/usr/bin/env node

/**
 * validate-specs.mjs — ts-morph-based AST validation for generated Playwright spec files.
 *
 * Checks spec files against 8 structural rules (4 ERROR + 4 WARNING):
 *   E1: waitForTimeout / setTimeout usage
 *   E2: TC ID full coverage (all TC-\d+ from test-cases.md must appear in specs)
 *   E3: Every test() has a Traceability comment
 *   E4: No DOM parent traversal locator('..')
 *   W1: Serial suite with >15 test() calls
 *   W2: Serial suite without afterAll cleanup
 *   W3: beforeEach containing login/loginViaUI calls
 *   W4: CSS class selectors (.xxx pattern in locator strings)
 *
 * Usage:
 *   node validate-specs.mjs <spec-directory> [--test-cases <path>]
 *
 * Output: JSON to stdout with { errors: [...], warnings: [...] }
 * Exit code: 0 if no errors (warnings OK), 1 if any errors present
 */

import { Project, SyntaxKind } from 'ts-morph';
import { readdirSync, readFileSync, existsSync, statSync } from 'node:fs';
import { join, extname, resolve } from 'node:path';

// ── CLI Argument Parsing ─────────────────────────────────────────────

function parseArgs(argv) {
  const args = argv.slice(2);
  let specDir = null;
  let testCasesPath = null;

  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--test-cases' && i + 1 < args.length) {
      testCasesPath = resolve(args[++i]);
    } else if (!args[i].startsWith('-')) {
      specDir = resolve(args[i]);
    }
  }

  if (!specDir) {
    console.error('Usage: node validate-specs.mjs <spec-directory> [--test-cases <path>]');
    process.exit(2);
  }

  return { specDir, testCasesPath };
}

// ── File Discovery ───────────────────────────────────────────────────

function findSpecFiles(dir) {
  if (!existsSync(dir)) {
    return { files: [], error: `Spec directory not found: ${dir}` };
  }

  const files = [];
  const entries = readdirSync(dir);
  for (const entry of entries) {
    const fullPath = join(dir, entry);
    const stat = statSync(fullPath);
    if (stat.isFile() && extname(entry) === '.ts') {
      files.push(fullPath);
    }
  }
  return { files, error: null };
}

// ── E2: TC ID Coverage ──────────────────────────────────────────────

function extractTcIdsFromTestCases(testCasesPath) {
  if (!testCasesPath || !existsSync(testCasesPath)) {
    return { ids: [], error: testCasesPath ? `test-cases.md not found: ${testCasesPath}` : null };
  }

  const content = readFileSync(testCasesPath, 'utf-8');
  const regex = /\bTC-\d+\b/g;
  const ids = [...new Set(content.match(regex) || [])];
  return { ids, error: null };
}

function extractTcIdsFromSpecs(specFiles) {
  const ids = new Set();
  for (const file of specFiles) {
    const content = readFileSync(file, 'utf-8');
    const regex = /\bTC-\d+\b/g;
    const matches = content.match(regex) || [];
    for (const m of matches) {
      ids.add(m);
    }
  }
  return [...ids];
}

// ── AST-based Rule Checking ──────────────────────────────────────────

function validateSpecs(specFiles, tcIdsFromTestCases) {
  const errors = [];
  const warnings = [];

  for (const filePath of specFiles) {
    const project = new Project({
      useInMemoryFileSystem: true,
      compilerOptions: {
        allowJs: true,
        strict: false,
      },
    });

    let sourceFile;
    try {
      const content = readFileSync(filePath, 'utf-8');
      sourceFile = project.createSourceFile(filePath, content, { overwrite: true });
    } catch (e) {
      warnings.push({
        rule: 'PARSE',
        file: filePath,
        line: 0,
        message: `Failed to parse file (reported as WARNING per fallback policy): ${e.message}`,
      });
      continue;
    }

    checkE1_WaitForTimeout(sourceFile, filePath, errors);
    checkE3_Traceability(sourceFile, filePath, errors);
    checkE4_DomTraversal(sourceFile, filePath, errors);
    checkW1_SerialSuiteSize(sourceFile, filePath, warnings);
    checkW2_SerialSuiteNoAfterAll(sourceFile, filePath, warnings);
    checkW3_BeforeEachLogin(sourceFile, filePath, warnings);
    checkW4_CssClassSelectors(sourceFile, filePath, warnings);
  }

  // E2: TC ID coverage check
  checkE2_TcIdCoverage(specFiles, tcIdsFromTestCases, errors);

  return { errors, warnings };
}

// ── E1: waitForTimeout / setTimeout ──────────────────────────────────

function checkE1_WaitForTimeout(sourceFile, filePath, errors) {
  sourceFile.forEachDescendant((node) => {
    if (node.getKind() === SyntaxKind.CallExpression) {
      const expr = node.getExpression();
      const text = expr.getText();

      // Check for page.waitForTimeout() or waitForTimeout() standalone
      if (/waitForTimeout/.test(text)) {
        errors.push({
          rule: 'E1',
          file: filePath,
          line: node.getStartLineNumber(),
          message: `Forbidden waitForTimeout call: ${text}`,
        });
      }

      // Check for setTimeout
      if (text === 'setTimeout') {
        errors.push({
          rule: 'E1',
          file: filePath,
          line: node.getStartLineNumber(),
          message: 'Forbidden setTimeout call — use waitForApiAction / withRetry instead',
        });
      }
    }
  });
}

// ── E2: TC ID Coverage ──────────────────────────────────────────────

function checkE2_TcIdCoverage(specFiles, tcIdsFromTestCases, errors) {
  if (!tcIdsFromTestCases || tcIdsFromTestCases.length === 0) {
    return; // No test-cases.md provided, skip this check
  }

  const specIds = new Set(extractTcIdsFromSpecs(specFiles));
  const missingIds = tcIdsFromTestCases.filter((id) => !specIds.has(id));

  if (missingIds.length > 0) {
    errors.push({
      rule: 'E2',
      file: '(coverage)',
      line: 0,
      message: `Missing TC IDs in spec files: ${missingIds.join(', ')}`,
    });
  }
}

// ── E3: Traceability Comment ────────────────────────────────────────

function checkE3_Traceability(sourceFile, filePath, errors) {
  const content = sourceFile.getFullText();

  sourceFile.forEachDescendant((node) => {
    if (node.getKind() !== SyntaxKind.CallExpression) return;

    const expr = node.getExpression();
    const exprText = expr.getText();

    // Only check test() calls, not test.describe, test.skip, etc.
    if (exprText !== 'test') return;

    // Check for Traceability in the node's leading comment ranges
    const start = node.getStart();
    const hasLeadingTraceability = hasTraceabilityCommentBefore(content, start);

    // Check for Traceability string inside the test callback
    const hasInlineTraceability = hasTraceabilityInNode(node);

    if (!hasLeadingTraceability && !hasInlineTraceability) {
      errors.push({
        rule: 'E3',
        file: filePath,
        line: node.getStartLineNumber(),
        message: `test() call missing Traceability comment: // Traceability:`,
      });
    }
  });
}

function hasTraceabilityCommentBefore(content, position) {
  // Look backwards from the node position for a Traceability comment
  const before = content.substring(Math.max(0, position - 500), position);
  return /\/\/\s*Traceability:/i.test(before);
}

function hasTraceabilityInNode(node) {
  let found = false;
  node.forEachDescendant((child) => {
    if (child.getKind() === SyntaxKind.StringLiteral) {
      if (/Traceability:/i.test(child.getText())) {
        found = true;
      }
    }
  });
  return found;
}

// ── E4: DOM Parent Traversal ────────────────────────────────────────

function checkE4_DomTraversal(sourceFile, filePath, errors) {
  sourceFile.forEachDescendant((node) => {
    if (node.getKind() !== SyntaxKind.CallExpression) return;

    const expr = node.getExpression();
    const exprText = expr.getText();

    // Check for .locator('..') calls
    if (!/locator$/i.test(exprText)) return;

    const args = node.getArguments();
    if (args.length === 0) return;

    const firstArg = args[0];
    const argText = firstArg.getText().replace(/^['"`]|['"`]$/g, '');

    if (argText === '..' || argText.includes('/..') || argText.startsWith('../')) {
      errors.push({
        rule: 'E4',
        file: filePath,
        line: node.getStartLineNumber(),
        message: `Forbidden DOM parent traversal: locator('${argText}')`,
      });
    }
  });
}

// ── W1: Serial Suite Size ───────────────────────────────────────────

function checkW1_SerialSuiteSize(sourceFile, filePath, warnings) {
  const serialSuites = findSerialDescribeCalls(sourceFile);

  for (const suite of serialSuites) {
    const testCount = countTestCalls(suite.node);
    if (testCount > 15) {
      warnings.push({
        rule: 'W1',
        file: filePath,
        line: suite.node.getStartLineNumber(),
        message: `Serial suite "${suite.name}" has ${testCount} test() calls (max 15)`,
      });
    }
  }
}

// ── W2: Serial Suite Without afterAll ───────────────────────────────

function checkW2_SerialSuiteNoAfterAll(sourceFile, filePath, warnings) {
  const serialSuites = findSerialDescribeCalls(sourceFile);

  for (const suite of serialSuites) {
    const hasAfterAll = suite.node.forEachDescendant((node) => {
      if (node.getKind() === SyntaxKind.CallExpression) {
        const expr = node.getExpression();
        const text = expr.getText();
        if (text === 'test.afterAll' || text === 'afterAll') {
          return true; // Found
        }
      }
      return undefined;
    });

    if (!hasAfterAll) {
      warnings.push({
        rule: 'W2',
        file: filePath,
        line: suite.node.getStartLineNumber(),
        message: `Serial suite "${suite.name}" missing afterAll cleanup`,
      });
    }
  }
}

// ── W3: beforeEach with login ───────────────────────────────────────

function checkW3_BeforeEachLogin(sourceFile, filePath, warnings) {
  sourceFile.forEachDescendant((node) => {
    if (node.getKind() !== SyntaxKind.CallExpression) return;

    const expr = node.getExpression();
    const exprText = expr.getText();

    if (exprText !== 'test.beforeEach' && exprText !== 'beforeEach') return;

    // Check inside the beforeEach callback for login calls
    const args = node.getArguments();
    const callback = args.find((a) =>
      a.getKind() === SyntaxKind.ArrowFunction ||
      a.getKind() === SyntaxKind.FunctionExpression
    );

    if (!callback) return;

    let hasLogin = false;
    callback.forEachDescendant((child) => {
      if (child.getKind() === SyntaxKind.CallExpression) {
        const childExpr = child.getExpression().getText();
        if (/login/i.test(childExpr)) {
          hasLogin = true;
        }
      }
    });

    if (hasLogin) {
      warnings.push({
        rule: 'W3',
        file: filePath,
        line: node.getStartLineNumber(),
        message: 'beforeEach contains login call — use beforeAll or storageState instead',
      });
    }
  });
}

// ── W4: CSS Class Selectors ─────────────────────────────────────────

function checkW4_CssClassSelectors(sourceFile, filePath, warnings) {
  sourceFile.forEachDescendant((node) => {
    if (node.getKind() !== SyntaxKind.CallExpression) return;

    const expr = node.getExpression();
    const exprText = expr.getText();

    // Check for locator/getByRole/etc. calls with CSS class selectors
    if (!/locator|getByText|getByRole|getByLabel|getByTestId|getByPlaceholder|getByAltText|getByTitle$/i.test(exprText)) return;
    // For locator specifically, check for .xxx pattern
    if (!/locator$/i.test(exprText)) return;

    const args = node.getArguments();
    if (args.length === 0) return;

    const firstArg = args[0];
    if (firstArg.getKind() !== SyntaxKind.StringLiteral) return;

    const argValue = firstArg
      .getText()
      .replace(/^['"`]|['"`]$/g, '');

    // Match CSS class selector: starts with '.' followed by a class name
    // But not '..' (which is E4) and not valid XPath/CSS combinators
    if (/^\.[a-zA-Z_][\w-]*/.test(argValue) || /\s\.[a-zA-Z_][\w-]*/.test(argValue)) {
      warnings.push({
        rule: 'W4',
        file: filePath,
        line: node.getStartLineNumber(),
        message: `CSS class selector detected: '${argValue}' — use role/name/testid selectors instead`,
      });
    }
  });
}

// ── Helpers ──────────────────────────────────────────────────────────

function findSerialDescribeCalls(sourceFile) {
  const suites = [];

  sourceFile.forEachDescendant((node) => {
    if (node.getKind() !== SyntaxKind.CallExpression) return;

    const expr = node.getExpression();
    const exprText = expr.getText();

    // Match test.describe.serial
    if (
      exprText === 'test.describe.serial' ||
      exprText === 'describe.serial'
    ) {
      const args = node.getArguments();
      const name = args.length > 0 && args[0].getKind() === SyntaxKind.StringLiteral
        ? args[0].getText().replace(/^['"`]|['"`]$/g, '')
        : '(unnamed)';

      suites.push({ node, name });
    }
  });

  return suites;
}

function countTestCalls(node) {
  let count = 0;
  // Only count direct test() calls, not test.describe/test.skip/etc.
  node.forEachDescendant((child) => {
    if (child.getKind() !== SyntaxKind.CallExpression) return;
    const expr = child.getExpression();
    if (expr.getText() === 'test') {
      count++;
    }
  });
  return count;
}

// ── Main ─────────────────────────────────────────────────────────────

function main() {
  const { specDir, testCasesPath } = parseArgs(process.argv);

  const { files: specFiles, error: dirError } = findSpecFiles(specDir);
  if (dirError) {
    const result = {
      errors: [{ rule: 'IO', file: specDir, line: 0, message: dirError }],
      warnings: [],
    };
    console.log(JSON.stringify(result, null, 2));
    process.exit(1);
  }

  if (specFiles.length === 0) {
    const result = {
      errors: [{ rule: 'IO', file: specDir, line: 0, message: 'No .ts spec files found in directory' }],
      warnings: [],
    };
    console.log(JSON.stringify(result, null, 2));
    process.exit(1);
  }

  // Extract TC IDs from test-cases.md if provided
  const { ids: tcIdsFromTestCases, error: tcError } = extractTcIdsFromTestCases(testCasesPath);
  if (tcError) {
    console.error(`Warning: ${tcError} (E2 check skipped)`);
  }

  const { errors, warnings } = validateSpecs(specFiles, tcIdsFromTestCases);

  const result = { errors, warnings };
  console.log(JSON.stringify(result, null, 2));

  process.exit(errors.length > 0 ? 1 : 0);
}

main();

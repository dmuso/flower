import { existsSync, readdirSync } from "node:fs";
import { dirname, join, relative } from "node:path";

export const testFileNameLintIssue =
  "unit test files must be named <code_file>.test.ts or <code_file>.test.tsx with a same-directory <code_file>.ts or <code_file>.tsx, or use <domain>.integration.test.ts / <domain>.integration.test.tsx.";

const skippedDirectories = new Set([".git", "tmp", "node_modules", "dist"]);

export function lintTestFileNames(root: string): string[] {
  const issues: string[] = [];
  walk(root, root, issues);
  issues.sort();
  return issues;
}

export function writeIssues(issues: string[], write: (text: string) => void): void {
  for (const path of issues) {
    write(`${path}: ${testFileNameLintIssue}\n`);
  }
  if (issues.length > 0) {
    write(`Test file name lint failed with ${issues.length} issue(s).\n`);
  }
}

function walk(dir: string, root: string, issues: string[]): void {
  const entries = readdirSync(dir, { withFileTypes: true });
  for (const entry of entries) {
    const fullPath = join(dir, entry.name);
    if (entry.isDirectory()) {
      if (skippedDirectories.has(entry.name)) {
        continue;
      }
      walk(fullPath, root, issues);
      continue;
    }
    if (!isTestFileName(entry.name)) {
      continue;
    }
    if (allowedTestFileName(fullPath, entry.name)) {
      continue;
    }
    const rel = relative(root, fullPath);
    issues.push(rel === "" ? entry.name : rel);
  }
}

function isTestFileName(name: string): boolean {
  return name.endsWith(".test.ts") || name.endsWith(".test.tsx");
}

function allowedTestFileName(path: string, name: string): boolean {
  if (name.endsWith(".integration.test.ts") || name.endsWith(".integration.test.tsx")) {
    return true;
  }
  const stem = name.endsWith(".test.tsx") ? name.slice(0, -".test.tsx".length) : name.slice(0, -".test.ts".length);
  const dir = dirname(path);
  return existsSync(join(dir, `${stem}.ts`)) || existsSync(join(dir, `${stem}.tsx`));
}

function main(): void {
  const root = process.argv[2] && process.argv[2] !== "" ? process.argv[2] : ".";
  const issues = lintTestFileNames(root);
  writeIssues(issues, (text) => {
    process.stderr.write(text);
  });
  if (issues.length > 0) {
    process.exit(1);
  }
}

if (import.meta.main) {
  main();
}

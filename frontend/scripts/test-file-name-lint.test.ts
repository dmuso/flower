import { describe, expect, test } from "bun:test";
import { mkdir, mkdtemp, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { dirname, join } from "node:path";

import { lintTestFileNames, testFileNameLintIssue, writeIssues } from "./test-file-name-lint";

async function writeTree(files: string[]): Promise<string> {
  const root = await mkdtemp(join(tmpdir(), "test-file-name-lint-"));
  for (const rel of files) {
    const path = join(root, rel);
    await mkdir(dirname(path), { recursive: true });
    await writeFile(path, "export {}\n");
  }
  return root;
}

describe("lintTestFileNames", () => {
  test("pairing allowed", async () => {
    const root = await writeTree(["Board.tsx", "Board.test.tsx", "cn.ts", "cn.test.ts"]);
    expect(lintTestFileNames(root)).toEqual([]);
  });

  test("pairing tsx source with ts test allowed", async () => {
    const root = await writeTree(["Board.tsx", "Board.test.ts"]);
    expect(lintTestFileNames(root)).toEqual([]);
  });

  test("integration allowed without sibling", async () => {
    const root = await writeTree(["auth.integration.test.ts", "board.integration.test.tsx"]);
    expect(lintTestFileNames(root)).toEqual([]);
  });

  test("extra suffix rejected", async () => {
    const root = await writeTree(["slice0.test.ts", "handler_helpers.test.ts"]);
    expect(lintTestFileNames(root)).toEqual(["handler_helpers.test.ts", "slice0.test.ts"]);
  });

  test("skips node_modules dist git and tmp", async () => {
    const root = await writeTree([
      "node_modules/pkg/unpaired.test.ts",
      "dist/unpaired.test.ts",
      ".git/unpaired.test.ts",
      "tmp/unpaired.test.ts",
      "ok.ts",
      "ok.test.ts",
    ]);
    expect(lintTestFileNames(root)).toEqual([]);
  });
});

describe("writeIssues", () => {
  test("prints path then summary on stderr format", () => {
    const chunks: string[] = [];
    writeIssues(["slice0.test.ts"], (text) => {
      chunks.push(text);
    });
    expect(chunks.join("")).toBe(
      `slice0.test.ts: ${testFileNameLintIssue}\nTest file name lint failed with 1 issue(s).\n`,
    );
  });
});

import { describe, expect, test } from "bun:test";
import { mkdtemp, writeFile } from "node:fs/promises";
import { join } from "node:path";
import { tmpdir } from "node:os";

import { normalisePathname, pickStaticFile, resolveStaticPath, shouldServeHtml } from "./serve-helpers";

describe("serve helpers", () => {
  test("normalisePathname maps root to index.html", () => {
    expect(normalisePathname("/")).toBe("/index.html");
    expect(normalisePathname("")).toBe("/index.html");
  });

  test("resolveStaticPath blocks path traversal", async () => {
    const distDir = await mkdtemp(join(tmpdir(), "flower-dist-"));
    const resolved = resolveStaticPath(distDir, "/../secrets.txt");

    expect(resolved).toBe(join(distDir, "index.html"));
  });

  test("pickStaticFile reports existence", async () => {
    const distDir = await mkdtemp(join(tmpdir(), "flower-dist-"));
    await writeFile(join(distDir, "index.html"), "<html></html>");
    await writeFile(join(distDir, "app.js"), "console.log('ok')");

    const existing = await pickStaticFile(distDir, "/app.js");
    expect(existing.exists).toBe(true);
    expect(existing.path).toBe(join(distDir, "app.js"));

    const missing = await pickStaticFile(distDir, "/missing.css");
    expect(missing.exists).toBe(false);
    expect(missing.path).toBe(join(distDir, "missing.css"));
  });

  test("shouldServeHtml accepts GET HTML requests", () => {
    expect(shouldServeHtml("GET", "text/html")).toBe(true);
    expect(shouldServeHtml("GET", "*/*")).toBe(true);
    expect(shouldServeHtml("GET", null)).toBe(true);
    expect(shouldServeHtml("POST", "text/html")).toBe(false);
    expect(shouldServeHtml("GET", "application/json")).toBe(false);
  });
});

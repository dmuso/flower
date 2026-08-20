import { cp, mkdir, readFile, rm, writeFile } from "node:fs/promises";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { solidPlugin } from "./plugins/solid-plugin";
import { rewriteIndexHTMLContents } from "./build-helpers";

const here = dirname(fileURLToPath(import.meta.url));
const projectRoot = resolve(here, "..");
const distDir = join(projectRoot, "dist");
const publicDir = join(projectRoot, "public");
const tailwindInput = join(projectRoot, "src/styles/app.css");
const tailwindOutput = join(projectRoot, "src/styles/tailwind.css");

async function main() {
  await rm(distDir, { recursive: true, force: true });
  await mkdir(distDir, { recursive: true });

  await buildTailwind();
  await copyPublicAssets();
  await bundleApplication();
  await rewriteIndexHTML();
}

async function buildTailwind() {
  const command = ["bunx", "tailwindcss", "-i", tailwindInput, "-o", tailwindOutput];

  const task = Bun.spawn(command, {
    cwd: projectRoot,
    stdout: "inherit",
    stderr: "inherit",
  });

  const exitCode = await task.exited;
  if (exitCode !== 0) {
    throw new Error("Tailwind build failed");
  }
}

async function copyPublicAssets() {
  await cp(publicDir, distDir, { recursive: true });
}

export function createBundleOptions(): Bun.BuildConfig {
  const nodeEnv = process.env.NODE_ENV ?? "production";

  return {
    entrypoints: [join(projectRoot, "src/index.tsx")],
    outdir: distDir,
    target: "browser",
    publicPath: "/",
    splitting: true,
    sourcemap: "linked",
    minify: true,
    define: {
      "process.env.NODE_ENV": JSON.stringify(nodeEnv),
      "process.env.ENVIRONMENT": JSON.stringify(process.env.ENVIRONMENT ?? "prod"),
    },
    env: "FRONTEND_*",
    plugins: [solidPlugin],
  };
}

async function bundleApplication() {
  const apiUrl = process.env.FRONTEND_API_URL;
  if (!apiUrl) {
    throw new Error("FRONTEND_API_URL must be set before building the frontend");
  }

  process.env.NODE_ENV = process.env.NODE_ENV ?? "production";

  const result = await Bun.build(createBundleOptions());

  if (!result.success) {
    result.logs.forEach((log) => console.error(log));
    throw new Error("Bun build failed");
  }

  console.log(`Built frontend to ${distDir}`);
}

async function rewriteIndexHTML() {
  const indexPath = join(distDir, "index.html");
  const original = await readFile(indexPath, { encoding: "utf8" });
  const updated = rewriteIndexHTMLContents(original);

  if (updated !== original) {
    await writeFile(indexPath, updated, { encoding: "utf8" });
  }
}

export { rewriteIndexHTMLContents } from "./build-helpers";

if (import.meta.main) {
  main().catch((error) => {
    console.error(error);
    process.exit(1);
  });
}

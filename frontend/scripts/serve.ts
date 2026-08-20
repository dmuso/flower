import { join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { dirname } from "node:path";

import { pickStaticFile, shouldServeHtml } from "./serve-helpers";

export { normalisePathname, resolveStaticPath, pickStaticFile } from "./serve-helpers";

const here = dirname(fileURLToPath(import.meta.url));
const distDir = resolve(here, "..", "dist");
const indexPath = join(distDir, "index.html");

const port = Number(process.env.FRONTEND_PORT ?? 4273);

Bun.serve({
  port,
  async fetch(request) {
    const url = new URL(request.url);
    const { path, exists } = await pickStaticFile(distDir, url.pathname);

    if (exists) {
      return new Response(Bun.file(path));
    }

    if (shouldServeHtml(request.method, request.headers.get("accept"))) {
      return new Response(Bun.file(indexPath), {
        headers: {
          "content-type": "text/html; charset=utf-8",
        },
      });
    }

    return new Response("Not Found", { status: 404 });
  },
});

console.log(`Frontend static server listening on ${port}`);

import { access } from "node:fs/promises";
import { join, resolve, sep } from "node:path";

export function normalisePathname(pathname: string): string {
  if (!pathname || pathname === "/") {
    return "/index.html";
  }

  return pathname.startsWith("/") ? pathname : `/${pathname}`;
}

export function resolveStaticPath(distRoot: string, pathname: string): string {
  const normalised = normalisePathname(pathname);
  const resolved = resolve(distRoot, `.${normalised}`);

  const prefix = distRoot.endsWith(sep) ? distRoot : `${distRoot}${sep}`;
  if (resolved === distRoot || resolved.startsWith(prefix)) {
    return resolved;
  }

  return join(distRoot, "index.html");
}

async function fileExists(path: string): Promise<boolean> {
  try {
    await access(path);
    return true;
  } catch {
    return false;
  }
}

export async function pickStaticFile(distRoot: string, pathname: string) {
  const resolved = resolveStaticPath(distRoot, pathname);
  return { path: resolved, exists: await fileExists(resolved) };
}

export function shouldServeHtml(method: string, acceptHeader: string | null): boolean {
  if (method !== "GET") {
    return false;
  }

  if (!acceptHeader) {
    return true;
  }

  return acceptHeader.includes("text/html") || acceptHeader.includes("*/*");
}

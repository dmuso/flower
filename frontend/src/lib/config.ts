export function resolveApiBaseUrl(): string {
  const configuredUrl = process.env.FRONTEND_API_URL;

  if (!configuredUrl) {
    throw new Error("FRONTEND_API_URL is required but was not provided");
  }
  return configuredUrl.replace(/\/$/, "");
}

export const API_BASE_URL = resolveApiBaseUrl();

export const ROUTES = {
  home: "/",
} as const;

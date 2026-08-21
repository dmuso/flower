import { API_BASE_URL } from "../config";

export class ApiRequestError extends Error {
  status: number;
  code?: string;

  constructor(message: string, status: number, code?: string) {
    super(message);
    this.name = "ApiRequestError";
    this.status = status;
    this.code = code;
  }
}

type ErrorEnvelope = {
  error?: { code?: string; message?: string };
};

export async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  if (init.body !== undefined && !headers.has("content-type")) {
    headers.set("content-type", "application/json");
  }
  let response: Response;
  try {
    response = await fetch(`${API_BASE_URL}${path}`, {
      ...init,
      headers,
      credentials: "include",
    });
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    throw new ApiRequestError(`request failed: ${message}`, 0);
  }
  if (response.status === 204) {
    return undefined as T;
  }
  const text = await response.text();
  let parsed: unknown = {};
  if (text) {
    parsed = JSON.parse(text) as unknown;
  }
  if (!response.ok) {
    const envelope = parsed as ErrorEnvelope;
    const message = envelope.error?.message ?? `request failed with status ${response.status}`;
    throw new ApiRequestError(message, response.status, envelope.error?.code);
  }
  return parsed as T;
}

export { API_BASE_URL };

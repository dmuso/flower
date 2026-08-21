import { API_BASE_URL, ApiRequestError } from "./core";

export type HealthResponse = {
  status: string;
  version: string;
};

export async function fetchHealth(signal?: AbortSignal): Promise<HealthResponse> {
  let response: Response;
  try {
    response = await fetch(`${API_BASE_URL}/health`, { signal });
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    throw new ApiRequestError(`health request failed: ${message}`, 0);
  }
  if (!response.ok) {
    throw new ApiRequestError(`health request failed with status ${response.status}`, response.status);
  }
  const body = (await response.json()) as HealthResponse;
  if (!body.status || !body.version) {
    throw new Error("health response is missing status or version");
  }
  return body;
}

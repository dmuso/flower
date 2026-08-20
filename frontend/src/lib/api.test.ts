import { afterEach, describe, expect, it, mock } from "bun:test";
import { ApiRequestError, fetchHealth } from "./api";

const originalFetch = globalThis.fetch;

afterEach(() => {
  globalThis.fetch = originalFetch;
});

describe("fetchHealth", () => {
  it("returns status and version from the API", async () => {
    globalThis.fetch = mock(async () => {
      return new Response(JSON.stringify({ status: "ok", version: "dev" }), {
        status: 200,
        headers: { "content-type": "application/json" },
      });
    }) as unknown as typeof fetch;

    const health = await fetchHealth();
    expect(health.status).toBe("ok");
    expect(health.version).toBe("dev");
  });

  it("throws ApiRequestError when the response is not ok", async () => {
    globalThis.fetch = mock(async () => {
      return new Response("nope", { status: 503 });
    }) as unknown as typeof fetch;

    try {
      await fetchHealth();
      throw new Error("expected fetchHealth to throw");
    } catch (error) {
      expect(error).toBeInstanceOf(ApiRequestError);
      expect((error as ApiRequestError).status).toBe(503);
    }
  });

  it("throws ApiRequestError when fetch itself fails", async () => {
    globalThis.fetch = mock(() => Promise.reject(new Error("network down"))) as unknown as typeof fetch;

    try {
      await fetchHealth();
      throw new Error("expected fetchHealth to throw");
    } catch (error) {
      expect(error).toBeInstanceOf(ApiRequestError);
      expect((error as ApiRequestError).status).toBe(0);
      expect((error as ApiRequestError).message).toContain("network down");
    }
  });

  it("throws when the payload is missing required fields", async () => {
    globalThis.fetch = mock(async () => {
      return new Response(JSON.stringify({ status: "ok" }), {
        status: 200,
        headers: { "content-type": "application/json" },
      });
    }) as unknown as typeof fetch;

    await expect(fetchHealth()).rejects.toThrow("health response is missing status or version");
  });
});

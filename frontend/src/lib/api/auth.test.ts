import { afterEach, describe, expect, it, mock } from "bun:test";
import { signUp } from "./auth";
import { ApiRequestError } from "./core";

const originalFetch = globalThis.fetch;

afterEach(() => {
  globalThis.fetch = originalFetch;
});

describe("auth api", () => {
  it("posts signup with credentials", async () => {
    globalThis.fetch = mock(async (input: RequestInfo | URL, init?: RequestInit) => {
      expect(String(input)).toContain("/api/v1/auth/signup");
      expect(init?.credentials).toBe("include");
      return new Response(JSON.stringify({ email: "maya@example.com" }), {
        status: 201,
        headers: { "content-type": "application/json" },
      });
    }) as unknown as typeof fetch;
    const res = await signUp("maya@example.com", "secret12");
    expect(res.email).toBe("maya@example.com");
  });

  it("surfaces API errors", async () => {
    globalThis.fetch = mock(async () => {
      return new Response(JSON.stringify({ error: { code: "email_taken", message: "That email already belongs to a user." } }), {
        status: 409,
      });
    }) as unknown as typeof fetch;
    try {
      await signUp("maya@example.com", "secret12");
      throw new Error("expected failure");
    } catch (error) {
      expect(error).toBeInstanceOf(ApiRequestError);
      expect((error as ApiRequestError).code).toBe("email_taken");
    }
  });
});

import { describe, expect, it } from "bun:test";
import { ROUTES, resolveApiBaseUrl } from "./config";

describe("config", () => {
  it("exposes stable paths", () => {
    expect(ROUTES.home).toBe("/");
    expect(ROUTES.signIn).toBe("/signin");
  });

  it("reads FRONTEND_API_URL from the environment", () => {
    const configured = process.env.FRONTEND_API_URL;
    if (!configured) {
      throw new Error("FRONTEND_API_URL must be set in tests");
    }
    expect(resolveApiBaseUrl()).toBe(configured.replace(/\/$/, ""));
  });
});

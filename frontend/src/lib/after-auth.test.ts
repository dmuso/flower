import { describe, expect, it } from "bun:test";
import { afterAuthPath } from "./after-auth";
import type { Me } from "./api/auth";

function me(partial: Partial<Me>): Me {
  return {
    id: "u1",
    username: "maya",
    email: "maya@example.com",
    display_name: "maya",
    email_verified_at: "2026-08-21T00:00:00Z",
    organisations: [],
    last_project: null,
    ...partial,
  };
}

describe("afterAuthPath", () => {
  it("lands on last project", () => {
    expect(
      afterAuthPath(
        me({
          last_project: { id: "p1", organisation_id: "o1", name: "Trail", slug: "trail" },
        }),
      ),
    ).toBe("/projects/p1");
  });

  it("names a project when an organisation exists", () => {
    expect(afterAuthPath(me({ organisations: [{ id: "o1", name: "Acme", role: "owner" }] }))).toBe(
      "/organisations/o1/projects/new",
    );
  });

  it("names an organisation when none exist", () => {
    expect(afterAuthPath(me({}))).toBe("/organisations/new");
  });
});

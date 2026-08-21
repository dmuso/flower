import { request } from "./core";

export type OrganisationSummary = {
  id: string;
  name: string;
  role: string;
};

export type ProjectSummary = {
  id: string;
  organisation_id: string;
  name: string;
  slug: string;
};

export type Me = {
  id: string;
  username: string;
  email: string;
  display_name: string;
  email_verified_at: string | null;
  organisations: OrganisationSummary[];
  last_project: ProjectSummary | null;
};

export async function signUp(email: string, password: string): Promise<{ email: string }> {
  return request("/api/v1/auth/signup", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export async function signIn(email: string, password: string): Promise<Me> {
  return request("/api/v1/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export async function requestMagicLink(email: string): Promise<{ email: string }> {
  return request("/api/v1/auth/magic-link", {
    method: "POST",
    body: JSON.stringify({ email }),
  });
}

export async function consumeMagicLink(token: string): Promise<Me> {
  return request("/api/v1/auth/magic-link/consume", {
    method: "POST",
    body: JSON.stringify({ token }),
  });
}

export async function consumeVerifyEmail(token: string): Promise<Me> {
  return request("/api/v1/auth/verify-email/consume", {
    method: "POST",
    body: JSON.stringify({ token }),
  });
}

export async function signOut(): Promise<void> {
  await request("/api/v1/auth/logout", { method: "POST" });
}

export async function fetchMe(): Promise<Me> {
  return request("/api/v1/me");
}

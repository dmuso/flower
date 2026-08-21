import { request } from "./core";

export type Organisation = {
  id: string;
  name: string;
};

export async function createOrganisation(name: string): Promise<Organisation> {
  return request("/api/v1/organisations", {
    method: "POST",
    body: JSON.stringify({ name }),
  });
}

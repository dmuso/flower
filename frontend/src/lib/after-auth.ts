import type { Me } from "./api/auth";

export function afterAuthPath(me: Me): string {
  if (me.last_project) {
    return `/projects/${me.last_project.id}`;
  }
  if (me.organisations.length > 0) {
    return `/organisations/${me.organisations[0].id}/projects/new`;
  }
  return "/organisations/new";
}

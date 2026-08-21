import { request } from "./core";

export type Project = {
  id: string;
  organisation_id: string;
  name: string;
  slug: string;
  point_scale: string;
  timezone: string;
  velocity_strategy: number;
  initial_velocity: number;
  iteration_start_weekday: number;
  iteration_length_days: number;
  created_at: string;
};

export async function createProject(organisationId: string, name: string): Promise<Project> {
  return request(`/api/v1/organisations/${organisationId}/projects`, {
    method: "POST",
    body: JSON.stringify({ name }),
  });
}

export async function fetchProject(projectId: string): Promise<Project> {
  return request(`/api/v1/projects/${projectId}`);
}

export async function fetchOrganisationProjects(organisationId: string): Promise<{ projects: Project[] }> {
  return request(`/api/v1/organisations/${organisationId}/projects`);
}

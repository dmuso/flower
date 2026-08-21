import { request } from "./core";

export type Story = {
  id: string;
  project_id: string;
  title: string;
  state: string;
  story_type: string;
  estimate: number | null;
  rank: string;
};

export type Pack = {
  current_points: number;
  denominator: number;
  velocity_source: string;
  current_window_ends_at: string;
};

export type StoryList = {
  stories: Story[];
  pack: Pack;
};

export async function fetchStories(projectId: string): Promise<StoryList> {
  return request(`/api/v1/projects/${projectId}/stories`);
}

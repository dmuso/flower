/** @jsxImportSource solid-js */
import { useNavigate, useParams } from "@solidjs/router";
import { Show, createResource } from "solid-js";
import { EmptyBoard } from "../features/board";
import { fetchMe, signOut } from "../lib/api/auth";
import { fetchProject } from "../lib/api/projects";
import { fetchStories } from "../lib/api/stories";

async function loadBoard(projectId: string) {
  const [project, stories, me] = await Promise.all([fetchProject(projectId), fetchStories(projectId), fetchMe()]);
  const org = me.organisations.find((item) => item.id === project.organisation_id);
  if (!org) {
    throw new Error("organisation not in session");
  }
  return { project, stories, organisationName: org.name };
}

export function BoardPage() {
  const params = useParams();
  const navigate = useNavigate();
  const [board] = createResource(() => params.projectId, loadBoard);

  async function onSignOut() {
    await signOut();
    navigate("/signin");
  }

  return (
    <Show when={board()} fallback={<p class="p-4 text-sm text-ink-700">Loading…</p>}>
      {(state) => (
        <EmptyBoard
          organisation={state().organisationName}
          project={state().project.name}
          pack={state().stories.pack}
          onSignOut={onSignOut}
        />
      )}
    </Show>
  );
}

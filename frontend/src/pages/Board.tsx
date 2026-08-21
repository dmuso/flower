/** @jsxImportSource solid-js */
import { A, useNavigate, useParams } from "@solidjs/router";
import { Match, Switch, createResource } from "solid-js";
import { Board } from "../features/board";
import { primaryButtonClass } from "../components/AuthCard";
import { ApiRequestError } from "../lib/api/core";
import { fetchMe, signOut, type Me } from "../lib/api/auth";
import { fetchProject } from "../lib/api/projects";
import { fetchStories } from "../lib/api/stories";
import { afterAuthPath } from "../lib/after-auth";

type BoardFailure = ApiRequestError & { me?: Me };

function withMe(err: unknown, me: Me): unknown {
  if (err instanceof ApiRequestError) {
    (err as BoardFailure).me = me;
  }
  return err;
}

function missingBoardHref(err: unknown): string {
  if (err instanceof ApiRequestError) {
    const me = (err as BoardFailure).me;
    if (me) {
      return afterAuthPath(me);
    }
  }
  return "/organisations/new";
}

async function loadBoard(projectId: string) {
  const [projectResult, storiesResult, meResult] = await Promise.allSettled([
    fetchProject(projectId),
    fetchStories(projectId),
    fetchMe(),
  ]);

  if (meResult.status === "rejected") {
    throw meResult.reason;
  }
  const me = meResult.value;

  if (projectResult.status === "rejected") {
    throw withMe(projectResult.reason, me);
  }
  if (storiesResult.status === "rejected") {
    throw withMe(storiesResult.reason, me);
  }

  const project = projectResult.value;
  const stories = storiesResult.value;
  const org = me.organisations.find((item) => item.id === project.organisation_id);
  if (!org) {
    throw withMe(new ApiRequestError("We can’t find that.", 404), me);
  }
  return { project, stories, organisationName: org.name };
}

function boardError(err: unknown): { signedOut: boolean; copy: string } {
  if (err instanceof ApiRequestError && err.status === 401) {
    return { signedOut: true, copy: "You’re signed out." };
  }
  return { signedOut: false, copy: "We can’t find that." };
}

function BoardSkeletons() {
  return (
    <div class="flex h-screen flex-col bg-paper" data-testid="board-skeletons">
      <div class="h-14 border-b border-paper-200 bg-paper" />
      <div class="flex min-h-0 flex-1 overflow-x-auto">
        <section class="flex h-full min-w-[22rem] flex-col border-r border-paper-300 bg-paper" data-testid="skeleton-icebox" />
        <section class="flex h-full min-w-[22rem] flex-col border-r border-paper-300 bg-paper" data-testid="skeleton-backlog" />
        <section class="flex h-full min-w-[22rem] flex-col border-r border-paper-300 bg-paper" data-testid="skeleton-current" />
        <section class="flex h-full min-w-[22rem] flex-col border-r border-paper-300 bg-paper" data-testid="skeleton-done" />
      </div>
    </div>
  );
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
    <Switch fallback={<BoardSkeletons />}>
      <Match when={board.error}>
        {(() => {
          const state = boardError(board.error);
          return (
            <main class="mx-auto flex min-h-screen max-w-md flex-col justify-center gap-6 bg-paper px-6 py-16">
              <p class="text-sm text-ink-700" data-testid="board-error">
                {state.copy}
              </p>
              {state.signedOut ? (
                <A class={primaryButtonClass} href="/signin">
                  Sign in
                </A>
              ) : (
                <A class={primaryButtonClass} href={missingBoardHref(board.error)}>
                  Back to your projects
                </A>
              )}
            </main>
          );
        })()}
      </Match>
      <Match when={board()}>
        {(state) => (
          <Board
            organisation={state().organisationName}
            project={state().project.name}
            pack={state().stories.pack}
            onSignOut={onSignOut}
          />
        )}
      </Match>
    </Switch>
  );
}

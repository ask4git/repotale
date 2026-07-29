// Fetches recent merged PRs (with diff + body/title) for a repo via the GitHub REST API.
import { githubFetch } from "../github-fetch.ts";

const DEFAULT_LIMIT = 20;
// ponytail: GitHub's PR list isn't filterable by merged-only, so we over-fetch
// closed PRs and filter client-side. 50 is a simple heuristic buffer, not a
// guarantee — a repo with many closed-but-unmerged PRs could yield fewer than
// `limit` results. Real pagination is future work if that turns out to matter.
const LIST_PAGE_SIZE = 50;

export type MergedPR = {
  number: number;
  title: string;
  body: string | null;
  htmlUrl: string;
  mergedAt: string;
  diff: string;
};

type GithubPullRequestListItem = {
  number: number;
  title: string;
  body: string | null;
  html_url: string;
  merged_at: string | null;
};

export async function fetchRecentMergedPRs(
  owner: string,
  repo: string,
  accessToken: string,
  limit = DEFAULT_LIMIT
): Promise<MergedPR[]> {
  const listRes = await githubFetch(
    `https://api.github.com/repos/${owner}/${repo}/pulls?state=closed&sort=updated&direction=desc&per_page=${LIST_PAGE_SIZE}`,
    accessToken,
    { accept: "application/vnd.github+json" }
  );
  const list = (await listRes.json()) as GithubPullRequestListItem[];
  const merged = list.filter((pr) => pr.merged_at !== null).slice(0, limit);

  // allSettled, not all: one failed diff fetch (rate limit, transient 5xx)
  // shouldn't discard every other diff already fetched successfully.
  const results = await Promise.allSettled(
    merged.map(async (pr): Promise<MergedPR> => {
      const diffRes = await githubFetch(
        `https://api.github.com/repos/${owner}/${repo}/pulls/${pr.number}`,
        accessToken,
        { accept: "application/vnd.github.v3.diff" }
      );
      const diff = await diffRes.text();
      return {
        number: pr.number,
        title: pr.title,
        body: pr.body,
        htmlUrl: pr.html_url,
        // safe: filtered to merged_at !== null above
        mergedAt: pr.merged_at as string,
        diff,
      };
    })
  );
  for (const result of results) {
    if (result.status === "rejected") {
      console.error("[pipeline] PR diff fetch failed:", result.reason);
    }
  }
  return results
    .filter((r): r is PromiseFulfilledResult<MergedPR> => r.status === "fulfilled")
    .map((r) => r.value);
}

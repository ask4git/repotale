// Fetches recent merged PRs (with diff + body/title) for a repo via the GitHub REST API.

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

async function githubFetch(url: string, accessToken: string, accept: string): Promise<Response> {
  const res = await fetch(url, {
    headers: {
      Authorization: `Bearer ${accessToken}`,
      Accept: accept,
      "X-GitHub-Api-Version": "2022-11-28",
    },
  });
  if (!res.ok) {
    throw new Error(`GitHub API request to ${url} failed: ${res.status} ${res.statusText}`);
  }
  return res;
}

export async function fetchRecentMergedPRs(
  owner: string,
  repo: string,
  accessToken: string,
  limit = DEFAULT_LIMIT
): Promise<MergedPR[]> {
  const listRes = await githubFetch(
    `https://api.github.com/repos/${owner}/${repo}/pulls?state=closed&sort=updated&direction=desc&per_page=${LIST_PAGE_SIZE}`,
    accessToken,
    "application/vnd.github+json"
  );
  const list = (await listRes.json()) as GithubPullRequestListItem[];
  const merged = list.filter((pr) => pr.merged_at !== null).slice(0, limit);

  return Promise.all(
    merged.map(async (pr) => {
      const diffRes = await githubFetch(
        `https://api.github.com/repos/${owner}/${repo}/pulls/${pr.number}`,
        accessToken,
        "application/vnd.github.v3.diff"
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
}

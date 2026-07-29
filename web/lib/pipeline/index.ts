import type { Post } from "../types.ts";
import { fetchRecentMergedPRs } from "./fetch-history.ts";
import { clusterPullRequests } from "./cluster.ts";
import { generatePost } from "./generate-post.ts";

export async function analyzeRepo(owner: string, repo: string, accessToken: string): Promise<Post[]> {
  const prs = await fetchRecentMergedPRs(owner, repo, accessToken);
  const clusters = clusterPullRequests(prs);
  // allSettled, not all: one failed generation (LLM error, bad response)
  // shouldn't discard every other already-succeeded post in the batch.
  const results = await Promise.allSettled(clusters.map(generatePost));
  for (const result of results) {
    if (result.status === "rejected") {
      console.error("[pipeline] generatePost failed:", result.reason);
    }
  }
  return results
    .filter((r): r is PromiseFulfilledResult<Post> => r.status === "fulfilled")
    .map((r) => r.value);
}

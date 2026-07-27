import type { Post } from "../types.ts";
import { fetchRecentMergedPRs } from "./fetch-history.ts";
import { clusterPullRequests } from "./cluster.ts";
import { generatePost } from "./generate-post.ts";

export async function analyzeRepo(owner: string, repo: string, accessToken: string): Promise<Post[]> {
  const prs = await fetchRecentMergedPRs(owner, repo, accessToken);
  const clusters = clusterPullRequests(prs);
  return Promise.all(clusters.map(generatePost));
}

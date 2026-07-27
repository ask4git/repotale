import type { MergedPR } from "./fetch-history.ts";

// ponytail: MVP clustering rule is "one merged PR = one post" — this matches
// what was manually validated in experiments/prompt-test/. Clustering commits
// that have no associated PR into a shared post is future work, out of scope
// here.
export type Cluster = MergedPR;

export function clusterPullRequests(prs: MergedPR[]): Cluster[] {
  return prs;
}

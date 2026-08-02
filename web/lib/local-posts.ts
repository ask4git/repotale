import { readdir, readFile } from "fs/promises";
import { homedir } from "os";
import path from "path";

export type LocalPost = {
  slug: string;
  title: string;
  excerpt: string;
  tags: string[];
  commitSha: string;
  publishedAt: string;
};

// mirrors what `repotale` (the CLI) writes to ~/.repotale/<repoId>/posts/*.json
export async function getLocalPosts(repoId: string): Promise<LocalPost[]> {
  const dir = path.join(homedir(), ".repotale", repoId, "posts");

  let files: string[];
  try {
    files = await readdir(dir);
  } catch {
    return [];
  }

  const posts = await Promise.all(
    files
      .filter((f) => f.endsWith(".json"))
      .map(async (f) => {
        const raw = await readFile(path.join(dir, f), "utf-8");
        return JSON.parse(raw) as LocalPost;
      }),
  );

  return posts.sort((a, b) => b.publishedAt.localeCompare(a.publishedAt));
}

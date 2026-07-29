// Shared GitHub REST API fetch helper: Bearer auth, ok-check, throw with
// response body on failure. Used wherever we call GitHub's API with a token
// (OAuth user token, GitHub App JWT, or installation token).
export async function githubFetch(
  url: string,
  token: string,
  opts: { method?: string; accept?: string } = {}
): Promise<Response> {
  const res = await fetch(url, {
    method: opts.method ?? "GET",
    headers: {
      Authorization: `Bearer ${token}`,
      Accept: opts.accept ?? "application/vnd.github+json",
      "X-GitHub-Api-Version": "2022-11-28",
    },
  });
  if (!res.ok) {
    throw new Error(`GitHub API request to ${url} failed: ${res.status} ${await res.text()}`);
  }
  return res;
}

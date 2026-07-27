import { test } from "node:test";
import assert from "node:assert/strict";
import { analyzeRepo } from "./index.ts";

// Live end-to-end check against a small real public repo. Requires a real
// GITHUB_TOKEN and LLM_API_KEY in the environment to actually execute.
test("analyzeRepo returns well-formed posts for a real public repo", async () => {
  const githubToken = process.env.GITHUB_TOKEN;
  const llmApiKey = process.env.LLM_API_KEY;
  assert.ok(githubToken, "GITHUB_TOKEN env var is required to run this test");
  assert.ok(llmApiKey, "LLM_API_KEY env var is required to run this test");

  const posts = await analyzeRepo("vercel", "ai", githubToken);

  assert.ok(posts.length > 0, "expected at least one generated post");
  for (const post of posts) {
    assert.equal(typeof post.slug, "string");
    assert.ok(post.slug.length > 0);
    assert.equal(typeof post.title, "string");
    assert.equal(typeof post.excerpt, "string");
    assert.ok(Array.isArray(post.tags));
    assert.equal(typeof post.prNumber, "number");
    assert.equal(typeof post.prUrl, "string");
    assert.equal(typeof post.publishedAt, "string");
  }
});

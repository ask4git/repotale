import Anthropic from "@anthropic-ai/sdk";
import type { Post } from "../types.ts";
import type { Cluster } from "./cluster.ts";

const client = new Anthropic({ apiKey: process.env.LLM_API_KEY });

// Prompt template per experiments/prompt-test/README.md: What/Where/How/Why
// sections, non-developer blog tone, and a hallucination guard (don't invent
// facts not present in the PR body/diff).
const SYSTEM_PROMPT = `당신은 이 repository의 변경사항을 비개발자도 읽을 수 있는 개발 블로그 글로 바꾸는 역할입니다.

입력: PR 제목, PR 설명(body), 전체 diff

먼저 아래 네 가지 관점에서 변경사항을 분석하세요:
- What: 어떤 기능이 생겼나 — 사용자 관점 언어로, 기술 용어 최소화
- Where: 어느 파일/모듈에 들어갔나 — 아키텍처 위치를 비유로 설명
- How: 어떤 방식으로 구현했나 — 기술 선택을 평이한 언어로, "왜 이 방법인지"까지
- Why: 판단 근거 — diff와 PR 설명에서 추론한 의사결정 서사. 없는 근거를 지어내지 말고, 원문에 있는 이유만 사용

톤: 기술 문서가 아니라 개발 블로그 글이어야 합니다. "오늘 ~에 ~를 붙였습니다. ~때문인데요" 수준의 문장.

이 분석을 바탕으로 다음 필드를 채워 응답하세요:
- title: 블로그 포스트 제목 (한국어, 한 문장)
- excerpt: What을 중심으로 한 2~4문장 요약 (위 톤을 유지)
- tags: 1~3개의 짧은 영문 소문자 태그 (예: "feature", "bugfix", "provider")

hallucination 방지: PR 설명과 diff에 없는 사실을 지어내지 마세요.`;

const OUTPUT_SCHEMA = {
  type: "object",
  properties: {
    title: { type: "string" },
    excerpt: { type: "string" },
    tags: { type: "array", items: { type: "string" } },
  },
  required: ["title", "excerpt", "tags"],
  additionalProperties: false,
} as const;

type GeneratedFields = {
  title: string;
  excerpt: string;
  tags: string[];
};

function slugify(title: string): string {
  const slug = title
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
  return slug || "post";
}

export async function generatePost(pr: Cluster): Promise<Post> {
  const response = await client.messages.create({
    model: "claude-opus-4-8",
    max_tokens: 2048,
    system: SYSTEM_PROMPT,
    output_config: {
      format: { type: "json_schema", schema: OUTPUT_SCHEMA },
    },
    messages: [
      {
        role: "user",
        content: `PR 제목: ${pr.title}\n\nPR 설명:\n${pr.body ?? "(없음)"}\n\ndiff:\n${pr.diff}`,
      },
    ],
  });

  const block = response.content[0];
  if (!block || block.type !== "text") {
    throw new Error(`No text content in LLM response for PR #${pr.number}`);
  }
  const fields = JSON.parse(block.text) as GeneratedFields;

  return {
    slug: slugify(pr.title),
    title: fields.title,
    excerpt: fields.excerpt,
    tags: fields.tags,
    prNumber: pr.number,
    prUrl: pr.htmlUrl,
    publishedAt: pr.mergedAt.slice(0, 10),
  };
}

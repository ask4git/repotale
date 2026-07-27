export type Post = {
  slug: string;
  title: string;
  excerpt: string;
  tags: string[];
  prNumber: number;
  prUrl: string;
  publishedAt: string;
};

// ponytail: hardcoded demo content from the manual prompt experiment
// (experiments/prompt-test/). Swap for the real analysis pipeline's output
// once it exists.
export const demoPosts: Record<string, Post[]> = {
  "vercel/ai": [
    {
      slug: "chat-transport-tui",
      title: "터미널 채팅창이 원격 AI 서버와도 대화할 수 있게 됐습니다",
      excerpt:
        "지금까지 터미널 채팅 화면은 내 컴퓨터 안에서 도는 AI하고만 연결할 수 있었습니다. 화면은 그대로 두고, 뒤에 원격 서버를 연결할 수 있는 입구를 하나 더 열었습니다.",
      tags: ["tui", "feature"],
      prNumber: 17246,
      prUrl: "https://github.com/vercel/ai/pull/17246",
      publishedAt: "2026-07-15",
    },
    {
      slug: "vertex-tuned-models",
      title: "구글 클라우드에서 직접 학습시킨 AI 모델도 쓸 수 있게 됐습니다",
      excerpt:
        "구글 Vertex의 맞춤 학습 모델은 기본 모델과 다른 주소로 서비스됩니다. 모델 이름이 endpoints/로 시작하면 자동으로 올바른 주소로 요청을 보내도록 분기했습니다.",
      tags: ["provider", "feature"],
      prNumber: 16498,
      prUrl: "https://github.com/vercel/ai/pull/16498",
      publishedAt: "2026-06-30",
    },
    {
      slug: "use-object-stable",
      title: "1년 넘게 실험적이라 적혀 있던 기능이 정식 기능이 됐습니다",
      excerpt:
        "useObject·StructuredObject는 이미 널리 쓰이던 기능이지만 이름 앞에 experimental_ 표시가 남아 있었습니다. 기존 이름은 그대로 두고, 정식 이름을 새로 내보내는 방식으로 승격했습니다.",
      tags: ["ui", "graduation"],
      prNumber: 16888,
      prUrl: "https://github.com/vercel/ai/pull/16888",
      publishedAt: "2026-07-08",
    },
  ],
};

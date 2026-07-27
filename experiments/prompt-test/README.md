# 프롬프트 실험 — 1차

> repo-story-idea.md 8절 "다음 액션 1" 실행: 실제 커밋/PR 히스토리로 포스트 3개 수동 생성 →
> 비개발자에게 읽혀보고 "이해됐다"는 반응이 나오는지 검증.

## 실험 대상

- repo: [vercel/ai](https://github.com/vercel/ai) (Vercel AI SDK, public)
- 선정 기준: feat 타입 merged PR 중 파일 수 20개 미만, 추가 라인 500줄 미만
  (너무 크지 않아 한 편의 글로 요약 가능한 규모)
- 입력 데이터: PR 제목, PR body(Background/Summary/Manual Verification), diff 전문
  (`gh pr view`, `gh pr diff`로 수집, 별도 코드 인터뷰 없음 — commit/diff/PR만으로
  충분한지가 이 실험의 핵심 가설)

## 사용한 프롬프트 (뼈대)

```
너는 이 repo의 변경사항을 비개발자도 읽을 수 있는 개발 블로그 글로 바꾸는 역할이다.

입력: PR 제목, PR 설명(Background/Summary), 전체 diff

출력 구조 (repo-story-idea.md 4절 기준):
- What: 어떤 기능이 생겼나 — 사용자 관점 언어로, 기술 용어 최소화
- Where: 어느 파일/모듈에 들어갔나 — 아키텍처 위치를 비유로 설명
- How: 어떤 방식으로 구현했나 — 기술 선택을 평이한 언어로, "왜 이 방법인지"까지
- Why: 판단 근거 — diff와 PR 설명에서 추론한 의사결정 서사. 없는 근거를
  지어내지 말고, 원문에 있는 이유만 사용

톤: 기술 문서 아님. "오늘 ~에 ~를 붙였습니다. ~때문인데요" 수준의 블로그 문장.
근거: 사용한 PR 번호를 글 끝에 항상 병기 (hallucination 검증 가능하게).
```

## 결과물

- [`01-chat-transport-tui.md`](./01-chat-transport-tui.md) — PR [#17246](https://github.com/vercel/ai/pull/17246)
- [`02-vertex-tuned-models.md`](./02-vertex-tuned-models.md) — PR [#16498](https://github.com/vercel/ai/pull/16498)
- [`03-use-object-stable.md`](./03-use-object-stable.md) — PR [#16888](https://github.com/vercel/ai/pull/16888)

## 다음 확인할 것

- [ ] 비개발자 1명에게 3편 읽혀보기 — "이해됐다"는 반응 나오는지
- [x] Why 섹션이 원문(PR body)에 없는 내용을 지어내진 않았는지 재검증 →
      [`why-verification.md`](./why-verification.md)
- [x] diff만 주고 PR body 없이도 같은 품질이 나오는지 (실제 서비스는 PR
      description이 없는 커밋도 다뤄야 함 — 더 어려운 케이스) →
      [`04-diff-only-mcp-tool-drift.md`](./04-diff-only-mcp-tool-drift.md),
      [`05-diff-only-fireworks-cache-affinity.md`](./05-diff-only-fireworks-cache-affinity.md)

## 결과

**Why 검증 (why-verification.md)**: 12개 Why 문장 중 1개가 명백한 hallucination으로
확인됐다 — 01편에서 "원격 AI 서버를 붙인 여러 팀이 요청해서 응답했다"고 썼지만, 실제
이슈(#17245)를 열어보니 작성자는 PR 작성자 본인 1명이었다. PR body에 이슈 번호만
언급되고 내용은 없을 때, 그 이슈를 열어보지 않고 "누가 왜 요청했는지"를 그럴듯하게
지어내는 패턴이 실제로 나왔다는 뜻이다. 근거 링크를 병기하는 것만으로는 부족하고,
링크된 이슈까지 열어서 대조해야 이 종류의 hallucination을 막을 수 있다. 그 외에는
경계선 사례(근거 없는 독자 심리 묘사) 1건을 빼면 나머지 문장은 전부 직접 근거이거나
안전한 범위의 추론이었다.

**diff-only 실험 (04, 05편)**: 하필 고른 두 PR 모두 diff 안에 changeset 요약과 문서
(MDX) 변경이 포함돼 있어서, PR 제목/본문이 없어도 Why를 뒷받침할 문장이 diff 자체에
남아 있는 운 좋은 케이스였다. 그래서 "기술적으로 왜 필요한지"는 diff만으로도 꽤 자신
있게 쓸 수 있었다. 다만 "왜 하필 이 시점인지", "관련 없어 보이는 변경이 왜 같이
묶였는지" 같은 맥락(=PR body의 Background/Related Issues가 주는 정보)은 diff만으로는
전혀 복원되지 않았고, 그 부분은 추측하지 않고 비워뒀다. **더 중요한 시사점**: 이번
실험은 diff에 changeset·문서가 딸려온 "우호적인" 경우만 봤다. changeset이나 문서
갱신이 없는 순수 코드 diff(실제로는 이런 커밋이 더 흔할 수 있다)였다면 Why 품질은
이번보다 훨씬 나빴을 가능성이 크다 — 이건 실제 제품 리스크로 남겨둔다. PR/커밋
description이 부실한 repo를 다뤄야 하는 게 이 서비스의 기본 전제인 만큼, "diff에
문서/changeset이 없는 최악의 케이스"로 후속 실험을 한 번 더 해볼 가치가 있다.

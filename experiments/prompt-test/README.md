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
- [ ] Why 섹션이 원문(PR body)에 없는 내용을 지어내진 않았는지 재검증
- [ ] diff만 주고 PR body 없이도 같은 품질이 나오는지 (실제 서비스는 PR
      description이 없는 커밋도 다뤄야 함 — 더 어려운 케이스)

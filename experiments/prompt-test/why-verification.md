# Why 섹션 검증 — 원문 대조

> `repo-story-idea.md` 9절 "다음 확인할 것" 두 번째 항목 실행: 3편의 Why 섹션이 PR body에 없는 내용을
> 지어내진 않았는지, PR body를 다시 받아와서 문장 단위로 대조한다.
>
> 분류 기준:
> - **직접 근거** — PR body에 문장/표현이 그대로 있거나 거의 그대로 옮긴 수준
> - **합리적 추론** — PR body에 명시되진 않았지만 diff·body 맥락에서 자연스럽게 나오는 해석
> - **근거 없음/지어냄** — PR body·diff 어디에도 없는 내용을 사실인 것처럼 서술함 (문제)

## 01-chat-transport-tui.md (PR #17246)

Why 섹션 원문 문장별 대조:

1. "터미널 화면이 'AI가 보내는 메시지 조각(UIMessageChunk)'을 그리는 방식으로 만들어져 있다는 걸
   알고 있었습니다"
   → **직접 근거**. PR body Background: "its terminal renderer already consumes `UIMessageChunk`"

2. "원격 서버와 통신하는 창구(ChatTransport)도 똑같은 형식의 메시지를 주고받고 있었죠"
   → **합리적 추론**. body에 "Remote agents exposed through an AI SDK UI message endpoint" 및
   "send the runner's existing UI-message history through the transport"라는 문장은 있지만,
   "똑같은 형식"이라고 명시적으로 대조하진 않았음. diff·AI SDK의 `ChatTransport` 설계상 타당한 해석.

3. "화면을 다시 만들 필요 없이 '입구만 넓히면' 재사용이 가능한 상황이었던 겁니다"
   → **합리적 추론/사실상 직접 근거**. Summary의 "accept either an agent or a ChatTransport... keep
   the existing Agent execution path unchanged"를 풀어쓴 수준.

4. **"실제로 원격 AI 서버를 붙인 팀들이 '이 터미널 화면을 우리도 쓰고 싶다'는 요청(이슈 #17245)을
   했고, 이번 변경은 그 요청에 대한 응답입니다"**
   → **근거 없음/지어냄 (문제)**. PR body는 "Fixes #17245"라고만 적혀 있을 뿐, 요청자가 누구인지,
   복수의 "팀들"이었는지는 전혀 언급하지 않는다. 실제로 이슈 #17245를 확인해보니 제목은 "Enable
   runAgentTUI for any ToolLoopAgent with ChatTransport"이고, 작성자는 `Albert-Gao` — 이 PR의
   작성자 본인과 동일 인물이다. 즉 "여러 팀이 요청해서 응답했다"는 서사는 사실이 아니라, **PR
   작성자가 본인이 필요해서 이슈를 먼저 올리고 본인이 직접 구현한 것**이다. "팀들이 요청했다"는
   복수·외부 수요 프레이밍은 원문에 없는 내용을 지어낸 것으로 분류한다. 실제 서비스라면 이런
   과장은 신뢰도를 깎는 명백한 hallucination 사례.

5. "작성자는 실제로 빌드한 뒤 터미널에서 'hello'라고 쳐서 원격 창구를 통해 응답이 돌아오는 것까지
   직접 확인했다고 밝혔습니다"
   → **직접 근거**. Manual Verification 섹션 그대로.

**요약**: 5문장 중 1문장이 원문에 없는 "복수 팀의 요청"이라는 디테일을 지어냈다. 나머지는 직접
근거 또는 합리적 추론 범위 안. 이슈 번호(#17245)까지 정확히 병기해서 "검증 가능하게" 만들어둔 건
잘한 부분이지만, 정작 그 이슈 내용은 확인하지 않고 "요청했다"는 서사를 얹은 것이 이번 실험에서
나온 유일한 실질적 hallucination이다.

## 02-vertex-tuned-models.md (PR #16498)

Why 섹션 원문 문장별 대조:

1. "구글 클라우드 쪽에서 맞춤 모델 기능 자체는 이미 존재했지만, 이 SDK가 그 주소 체계를 몰라서
   못 쓰는 상태였습니다 (연결된 이슈 #6084)"
   → **직접 근거**. PR body Background 문단 그대로: "Google Vertex has tuned models but AI SDK
   didn't support them because they're served from a different endpoint" + "Fixes
   https://github.com/vercel/ai/issues/6084".

2. "작성자는 실제로 gemini-2.5-flash를 학습시켜 엔드포인트에 배포한 뒤, 이 기능으로 직접 텍스트
   생성을 시켜보고 정상 동작을 확인했다고 밝혔습니다"
   → **직접 근거**. Manual Verification: "Created a tuned model of gemini-2.5-flash on an endpoint
   in Google Cloud, and tested it using a new example."

3. "즉, '이미 구글이 제공하는 기능인데 SDK가 못 따라가고 있던 격차'를 메운 변경입니다"
   → **합리적 추론**. body가 명시적으로 "격차(gap)"라는 단어를 쓰진 않지만, Background 문단의
   내용을 한 문장으로 요약한 것으로 과장이나 왜곡 없음.

**요약**: 이 포스트는 문제없음. 세 문장 모두 직접 근거이거나 안전한 요약. 지어낸 디테일 없음.

## 03-use-object-stable.md (PR #16888)

Why 섹션 원문 문장별 대조:

1. "PR 설명에 따르면 이 기능은 이미 '성숙했고, 문서화도 잘 되어 있고, 널리 쓰이고 있는'
   상태였습니다"
   → **직접 근거**. Background: "The API is mature, widely documented, and heavily used."

2. "즉 실력은 정식 기능인데 이름표만 '실험적'으로 남아 있어서 사용자들이 '이거 써도 되는 거
   맞아?'하고 불안해할 수 있는 상황이었던 겁니다"
   → **합리적 추론(경계선)**. "사용자들이 불안해했다"는 감정 묘사는 PR body 어디에도 없다. 다만
   "experimental_ 접두사를 1년 넘게 달고 있었던 성숙한 API"라는 사실관계로부터 자연스럽게
   나올 수 있는 흔한 서술적 장치(narrative device)로 판단해 완전한 지어냄까지는 아니라고 본다.
   다만 이런 식의 "독자 심리 대변" 문장은 다음 실험에서 조금 더 보수적으로 쓸 필요가 있다 —
   본문에 실제로 그런 불만/혼란이 언급됐는지(예: 링크된 이슈)까지 확인했어야 더 안전했다.

3. "이 변경은 더 큰 계획(이슈 #16562의 '졸업 전략')의 일부로, 다음 정식 버전(v8)에서 옛날
   이름을 완전히 제거하기 전 중간 단계입니다"
   → **직접 근거**. Background: "Per the graduation strategy in #16562, this promotes the stable...
   names while keeping the experimental names as @deprecated aliases until they are removed in v8".

4. "작성자는 실제로 예제 앱을 띄워서 새 이름으로 임포트한 코드가 알림 카드 3개를 스트리밍으로
   정상 렌더링하는 것까지 확인했습니다"
   → **직접 근거**. Manual Verification 섹션과 정확히 일치("the three notification cards... rendered
   progressively").

**요약**: 지어낸 사실은 없음. 다만 2번 문장처럼 "사용자 심리"를 대변하는 표현은 근거 없이
분위기를 얹은 것이라 다음 실험에서는 주의가 필요한 패턴으로 기록해둔다.

## 종합

- 검증 대상 Why 문장 총 12개 중 **1개(01편, "팀들이 요청했다")가 명백한 hallucination**,
  1개(03편, "사용자들이 불안해했다")가 근거 없는 감정 서술로 경계선 사례.
  나머지 10개는 직접 근거 또는 안전한 범위의 합리적 추론.
- hallucination이 발생한 패턴: PR body에 "이슈 번호"만 언급된 경우, 그 이슈의 **내용을
  확인하지 않고** "누가, 왜 요청했는지"를 그럴듯하게 지어내는 경향. 이슈/PR 링크를 근거로
  병기하는 것만으로는 부족하고, **링크된 이슈 자체도 열어서 대조해야** 이런 종류의 hallucination을
  막을 수 있다는 게 실질적 시사점.
- 근거 없는 심리 서술("불안해할 수 있는", "궁금해할 만한")도 완전한 사실 왜곡은 아니지만 정확히는
  "원문에 없는 내용"이므로, 프롬프트에 "독자/사용자의 감정을 추측해서 서술하지 말 것" 같은 가드를
  추가하는 게 좋겠다.

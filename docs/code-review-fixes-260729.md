# 코드리뷰 수정 기록 (PR #1)

> 대상: `integration/parallel-tracks` → `main` PR #1 (CLI/GitHub App/파이프라인/기획 4개
> 병렬 트랙 통합)
> 리뷰 방식: `/code-review` 스킬, high effort — 8개 앵글(정확성 3, 정리 3, altitude 1,
> CLAUDE.md 컨벤션 1) 병렬 탐색 → 후보 중복제거 → 후보 10개 각각 독립 검증 에이전트로
> 1표 검증 → 10개 전부 CONFIRMED

## 1. GitHub App private key 형식 불일치 (최고 심각도)

- **파일**: `web/lib/github-app.ts`
- **문제**: `jose`의 `importPKCS8()`는 PKCS#8 PEM(`BEGIN PRIVATE KEY`)만 받는데,
  GitHub App이 실제로 발급하는 키는 PKCS#1(`BEGIN RSA PRIVATE KEY`) 형식. 이 PR
  자체가 작성한 `docs/github-app-setup.md`를 그대로 따라 키를 발급받아 붙여넣으면
  `signAppJwt()`가 매번 실패.
- **판단근거**: 검증 에이전트가 `node_modules/jose` 소스를 직접 열어
  `importPKCS8`이 PEM 헤더 문자열을 검사해서 PKCS#8이 아니면 즉시 `TypeError`를
  던지는 걸 확인. GitHub App 키가 PKCS#1로 발급되는 건 잘 알려진 함정
  (`openssl pkcs8 -topk8` 변환이 필요하다고 흔히 문서화됨).
- **수정**: `importPKCS8` 대신 Node 내장 `crypto.createPrivateKey()` 사용.
  PKCS#1/PKCS#8을 자동 판별해서 파싱하고, 결과 `KeyObject`를 `jose`의
  `SignJWT.sign()`에 그대로 넘김 (`sign()`은 `KeyObject`를 공식 지원하는 것을
  타입 정의로 확인). 변환 로직을 따로 안 짜도 되는 가장 단순한 해법.

## 2. 포스트 slug 충돌

- **파일**: `web/lib/pipeline/generate-post.ts`
- **문제**: `slug`가 PR 제목만으로 생성돼서, 제목이 같거나 비슷한 PR(Dependabot
  버전업 PR 등)이 한 배치(`analyzeRepo`는 PR 최대 20개 처리) 안에 여러 개 있으면
  slug가 겹침.
- **판단근거**: `Post.slug`는 `web/app/repo/[owner]/[repo]/page.tsx`에서 React 리스트
  `key`로 쓰이므로, 겹치면 단순 URL 중복이 아니라 실제 렌더링 버그(중복 key)로
  이어짐. 검증 에이전트가 소비처까지 추적해서 확인.
- **수정**: `slugify(title, prNumber)`로 시그니처 변경, PR 번호를 slug 끝에
  붙여서 배치 내 유일성 보장.

## 3. CLI 로그인 시 개행 없는 입력에서 토큰 유실

- **파일**: `cli/main.go` (`cmdLogin`)
- **문제**: `reader.ReadString('\n')`이 trailing newline 없는 입력(예:
  `echo -n "$TOKEN" | repotale login`)에서 데이터를 다 읽고도 `io.EOF`를
  같이 반환하는데, 기존 코드는 `err != nil`만 보고 바로 실패 처리 → 정상 읽은
  토큰을 버림.
- **판단근거**: 검증 에이전트가 이 경로가 `GITHUB_TOKEN` env var 미설정 시의
  일반적인 stdin 경로임을 확인, 개행 없는 파이프 입력은 실제로 흔한 시나리오.
- **수정**: `err != nil && err != io.EOF` 조건으로 변경 — EOF는 "trailing
  newline 없음"으로 취급하고 읽은 데이터를 그대로 사용.

## 4. 파이프라인 배치 생성이 fail-fast

- **파일**: `web/lib/pipeline/index.ts`
- **문제**: `Promise.all(clusters.map(generatePost))`는 하나라도 실패하면
  전체가 reject됨 — PR 20개 중 1개의 LLM 호출이 실패하면 이미 성공(=이미 비용
  지불)한 나머지 19개 결과까지 통째로 버려짐.
- **판단근거**: 아직 `web/app/api/webhooks/github/route.ts`가 `analyzeRepo`를
  호출하지 않아 지금 당장 터지는 버그는 아니지만, 연결되는 순간 바로
  재현되는 구조적 결함.
- **수정**: `Promise.allSettled`로 변경. 실패한 항목은 `console.error`로
  로그만 남기고, 성공한 `Post`만 필터링해서 반환.

## 5. PR diff 조회도 동일한 fail-fast + 무제한 동시요청

- **파일**: `web/lib/pipeline/fetch-history.ts`
- **문제**: PR별 diff를 `Promise.all`로 동시에 최대 20개까지 요청. 하나라도
  실패(일시적 5xx, secondary rate limit)하면 이미 받아온 나머지 diff까지 버려짐.
  GitHub는 한 토큰에서 동시 요청이 많으면 secondary rate limit에 걸릴 수 있다고
  문서화돼 있어 실제 리스크.
- **수정**: 4번과 동일하게 `Promise.allSettled` + 실패 로깅 + 성공분만 반환으로
  변경.

## 6. `jose`가 package.json에 미선언

- **파일**: `web/lib/github-app.ts`, `web/package.json`
- **문제**: `jose`를 직접 import하는데 `dependencies`엔 없고, `next-auth`가
  끌어오는 transitive dependency로 npm이 우연히 최상위에 hoist해준 덕에
  동작하던 상태.
- **판단근거**: npm의 hoisting은 보장이 아니라서, `next-auth`가 `jose` 버전을
  바꾸거나 pnpm처럼 isolated node_modules를 쓰는 환경으로 옮기면 즉시 깨짐.
- **수정**: `web/package.json`에 `jose` 직접 의존성으로 추가(현재 resolve된
  `^6.2.4`로 고정), lockfile 갱신.

## 7. GitHub fetch 패턴 3곳 중복

- **파일**: `web/lib/pipeline/fetch-history.ts`, `web/lib/github-app.ts`
- **문제**: "Bearer 인증 fetch → `res.ok` 체크 → 실패 시 throw" 패턴이
  `fetch-history.ts`(자체 로컬 헬퍼), `github-app.ts`(인라인), 그리고 기존
  `web/app/connect/page.tsx`(이번 PR 범위 밖, 안 건드림)까지 3곳에 흩어져 있었음.
- **판단근거**: 검증 에이전트가 세 곳의 실제 차이(GET/POST, Accept 헤더, 에러
  바디 처리 방식)를 비교해서 "억지로 합친 게 아니라 자연스럽게 파라미터화되는
  중복"이라고 확인. CLAUDE.md의 "단일 사용처엔 추상화 금지" 규칙에 걸리지 않음
  (3곳 이상 사용).
- **수정**: `web/lib/github-fetch.ts`에 공용 `githubFetch(url, token, opts)`
  헬퍼 추출, `fetch-history.ts`와 `github-app.ts` 둘 다 이걸 쓰도록 변경.
  기존에 동작 중인 `web/app/connect/page.tsx`는 이번 PR 범위 밖이라 안 건드림
  (surgical changes 원칙).

## 8. Webhook 핸들러의 `JSON.parse` 예외 미처리

- **파일**: `web/app/api/webhooks/github/route.ts`
- **문제**: 서명 검증 통과 후 `JSON.parse(rawBody)`에 try/catch가 없어서,
  형식이 깨진 바디가 오면 처리되지 않은 예외로 던져짐.
- **판단근거**: Next.js가 프레임워크 레벨에서 잡아서 500으로 응답하기 때문에
  서버가 죽는 건 아니고, GitHub가 서명한 페이로드라 실제로 깨진 JSON이 올
  가능성은 낮음 — 그래도 500은 GitHub의 재시도를 유발하고 로그에 처리되지 않은
  예외로 잡혀 노이즈가 됨. 두 줄이면 막을 수 있어서 수정.
- **수정**: `JSON.parse`를 try/catch로 감싸고 실패 시 400 응답. 페이로드
  타입을 `Record<string, unknown>` 대신 실제 쓰는 필드만 담은 `PushPayload`
  타입으로 좁힘.

## 9. CLI 에러 처리 8회 중복

- **파일**: `cli/main.go`
- **문제**: `fmt.Fprintln(os.Stderr, "error:", err); os.Exit(1)` 패턴이
  `cmdLogin`/`cmdConnect`/`cmdOpen`에 걸쳐 8번 그대로 복붙돼 있었음.
- **판단근거**: 검증 에이전트가 8곳 전부 바이트 단위로 동일한 걸 확인 —
  단일 사용처가 아니라 8회 반복이라 CLAUDE.md의 "단일 사용처 추상화 금지"
  규칙에 해당하지 않음.
- **수정**: `die(err error)` 헬퍼 하나 추가, 8곳 전부 `die(err)` 한 줄로 축약.
  (3번 수정에서 io.EOF를 걸러내는 분기가 새로 생겼는데, 그 분기도 `die()`를
  그대로 재사용.)

## 10. `cluster.ts`의 단일 PR 전제가 향후 진짜 클러스터링과 충돌 — 코드 미수정, 기록만

- **파일**: `web/lib/pipeline/cluster.ts`, `web/lib/types.ts`
- **문제**: 지금 `cluster.ts`는 "PR 1개 = 포스트 1개" identity passthrough인데,
  `Post` 타입은 `prNumber`/`prUrl`을 단수로 고정. `repo-story-idea.md`가 정의한
  진짜 제품 단위는 "여러 커밋을 하나의 기능 이야기로 클러스터링"(PR 여러 개가
  포스트 하나가 될 수 있음)이라, 나중에 진짜 클러스터링을 구현하면 `Post` 타입을
  단수→복수로 바꿔야 하고 그게 `web/lib/types.ts`부터 UI 소비처
  (`page.tsx`, `posts.ts`)까지 파급됨. `cluster.ts`의 `ponytail:` 주석은 "PR 없는
  커밋 클러스터링"만 후속 과제로 적어놔서, 이 파급 범위를 과소평가하고 있었음.
- **왜 코드를 안 고쳤나**: 지금 시점에 `Post`를 복수형으로 미리 바꾸는 건
  아직 존재하지 않는 요구사항에 대비한 speculative 설계 변경 — CLAUDE.md
  "simplicity first: 요청되지 않은 유연성 금지" 원칙과 정면으로 충돌. 실제
  멀티 PR 클러스터링을 구현하는 시점에 타입을 같이 바꾸는 게 맞음.
  **대신 이 문서에 알려진 설계 부채로 기록**해서, 나중에 그 작업을 시작할
  사람이 "포함된 범위"를 `cluster.ts` 주석보다 정확히 파악할 수 있게 남김.

## 검증

- `cd web && npx tsc --noEmit` — clean
- `cd web && npm run build` — clean, 라우트 8개 정상 생성
- `cd cli && go build ./... && go vet ./... && go test ./... && gofmt -l .` — clean
- `node --experimental-strip-types --test lib/pipeline/index.test.ts` — 모듈 체인
  정상 로드 확인 (실제 실행은 여전히 `GITHUB_TOKEN`/`LLM_API_KEY` 없어서
  가드에서 멈춤 — 이전과 동일, 이번 수정으로 새로 깨진 것 없음)

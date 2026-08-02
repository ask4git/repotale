<p align="center">
  <strong>repotale</strong> — 커밋을 하나씩 읽는 대신, 이야기로 읽는다
</p>

<p align="center">
  vibe coding으로 코드는 빨리 쌓이는데, 그걸 이해하는 속도는 못 따라간다 — "comprehension debt".<br>
  repotale은 git 히스토리를 <strong>비개발자도 읽을 수 있는 블로그</strong>로 바꿔서 그 빚을 갚는다.
</p>

---

## 30초 안에 보기

repo 안에서 실행하면, 최근 커밋을 읽고 로컬 웹 화면에 블로그로 띄운다.

```bash
go install github.com/ask4git/repotale/repotale@v0.1.0
export PATH="$PATH:$(go env GOPATH)/bin"  # 이미 PATH에 있으면 스킵됨, `repotale`이 바로 실행되면 필요 없음
repotale
```

버전 태그(`@v0.1.0`) 대신 `@latest` 쓰면 방금 push된 커밋을 Go 모듈 프록시가 아직 못 봐서 옛날 버전이 깔릴 수 있음 - 태그 박힌 버전이 제일 확실함. 이후 업데이트는 `repotale update`.

`command not found: repotale`이 뜨면 Go의 설치 경로(`~/go/bin`)가 PATH에 없는 것 - 위 `export` 한 줄이면 그 세션에서 바로 해결되고, 매번 새 터미널에서도 되게 하려면:

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc && source ~/.zshrc
```

```
Select language / 언어를 선택하세요:
  [1] English
  [2] 한국어
Choice / 선택 [1]: 2
분석할 repo 경로 [/Users/you/my-project]:
model: Claude Code (claude CLI, 로컬 로그인 사용)
이 claude CLI로 분석을 진행할까요? (y/n) [n]: y
새로 분석할 커밋 10개, diff 기준 예상 입력 토큰 약 3,500개 (대략치, 실제와 다를 수 있음)
진행할까요? (y/n) [n]: y
분석 톤 (soft/medium/hard) [soft]:
생성할 글의 언어 (en/ko) [ko]:
10개 커밋 분석 중 (/Users/you/my-project)...
  aeda217 GitHub 하나뿐이던 로그인 버튼에 Google·GitLab·Apple·Passkey를 붙였습니다  (claude 세션 6c33497a-...)
  ...
10개 새로 생성, 총 10개 포스트, 저장 위치: ~/.repotale/my-project/posts
127.0.0.1:4321 에서 서비스 중 - 종료하려면 Ctrl+C
```

로컬 `claude` CLI(로그인된 상태)만 있으면 끝 — API 키도, Node.js도, 로그인 화면도 필요 없다. `repotale` 바이너리 하나가 분석부터 웹 서버까지 다 함. 다시 실행해도 이미 분석한 커밋은 건너뛴다.

언어 선택은 최초 1회만 물어보고 `~/.repotale/settings`에 저장된다 (`.zshrc`처럼 사람이 직접 열어 고칠 수도 있는 평문 파일). 나중에 바꾸고 싶으면 `repotale --language`. UI 언어와 별개로, **생성되는 글의 언어**도 매 실행마다 따로 고를 수 있다 (기본값은 UI 언어를 따라감).

분석 직전엔 diff 크기 기준 대략적인 토큰 사용량을 보여주고 진행 여부를 물어봐서, 모르는 새 비용이 나가지 않게 했다.

### 실제로 이 repo를 분석하면 나오는 글 (진짜 결과, 손질 안 함)

> **GitHub 하나뿐이던 로그인 버튼에 Google·GitLab·Apple·Passkey를 붙였습니다**
>
> 오늘 로그인 화면을 뜯어고쳤습니다. 지금까지는 GitHub 계정이 없으면 아예 들어올 방법이 없었는데, 이제 Google, GitLab, Apple로도 로그인할 수 있고 비밀번호 없이 지문·얼굴로 여는 Passkey, 회사 계정으로 한 번에 들어오는 SSO까지 열어뒀습니다. 어떤 걸로 들어오든 이메일이 같으면 자동으로 한 사람의 계정으로 합쳐지도록 해서, "예전에 뭘로 가입했더라" 하고 헤매는 일이 없게 했고요.
>
> `auth` `feature` `refactor`

commit 메시지나 diff엔 저런 문장이 없다 — 그걸 읽고 재구성한 게 repotale이 하는 일이다.

## 왜

- 바이브코딩은 사람의 이해 속도를 앞지른다. 코드는 있는데 "왜 이렇게 만들었는지" 설명할 사람이 없는 상태가 쌓인다.
- commit log, changelog는 개발자용 포맷이다. Product Owner, 비개발자 창업자, 갓 합류한 개발자는 못 읽는다.
- repotale은 What(무슨 기능) / Where(어디에) / How(어떻게) / Why(왜) 네 관점으로 커밋을 재구성해서, 코드를 모르는 사람도 읽을 수 있는 서사로 바꾼다.

기획 배경 전체는 [repo-story-idea.md](./repo-story-idea.md), 인지부채 개념은 [docs/comprehension-debt-concept.md](./docs/comprehension-debt-concept.md) 참고.

## 구조

```
repotale/       Go CLI — repotale 실행 시 로컬 git log 분석 + 내장 웹서버로 결과 표시 (기본 동작)
                repotale login / connect / open / update — GitHub 원격 연동, 자가 업데이트 보조 명령
web/            Next.js 앱 — (진행 중) 로그인/원격 repo 연동용, 로컬 분석과는 무관
docs/           기획·설계 문서, 미해결 이슈
experiments/    프롬프트 실험 기록
environments/   환경별 .env 예시 (local/devel/staging/production)
docker-compose.yml  web + Postgres 실행 (원격 연동 모드용)
```

## 지금 상태

| 기능 | 상태 |
|---|---|
| 로컬 repo 분석 (`repotale` 실행 → 커밋 → 블로그) | 됨 |
| 내장 웹서버로 결과 보기 (Node/DB/로그인 불필요, `repotale` 바이너리 하나로) | 됨 |
| 이미 분석한 커밋 재분석 스킵 | 됨 |
| 분석 전 예상 토큰 사용량 확인 + 진행 여부 확인 | 됨 |
| 분석 톤 프리셋 (soft/medium/hard) + `.repotale/prompt.md` 완전 커스텀 | 됨 |
| UI 언어(en/ko) 선택 및 저장, 생성되는 글의 언어 별도 선택 | 됨 |
| GitHub/Google/GitLab/Apple 로그인 (원격 연동용) | 됨 |
| GitHub repo 연동 화면 (`/connect`) | 됨 |
| push 시 자동 재분석 (GitHub 웹훅) | 웹훅 수신은 되나 분석 파이프라인 미연결 |
| 원격 repo(PR 기반) 분석 결과를 웹에서 보기 | 미완성 — `/repo/[owner]/[repo]`는 아직 데모 데이터 |

즉 **로컬 모드는 실사용 가능**, **원격(GitHub) 모드는 로그인/연동까지만 되고 분석·발행은 아직**.

## 설치 / 실행

### 로컬 분석 (지금 되는 것)

필요한 건 [Claude Code](https://claude.com/claude-code) 설치 + 로그인뿐. API 키 불필요.

```bash
go install github.com/ask4git/repotale/repotale@v0.1.0  # 또는: make cli-build
repotale            # 분석하고 싶은 repo에서
```

결과는 `~/.repotale/<repo이름>/posts/*.json`에 저장되고, 분석이 끝나면 `repotale`이 알아서 `http://127.0.0.1:4321`에 웹 서버를 띄워 보여준다.

톤 프리셋 3개(soft/medium/hard)는 `~/.repotale/presets/*.md`에 실제 파일로 만들어져 있어서 직접 열어 고칠 수 있다 — 처음 한 번만 생성되고, 그 후엔 덮어쓰지 않는다. repo에 `.repotale/prompt.md`를 만들면 프리셋/톤 선택 자체를 건너뛰고 그 내용을 그대로 쓴다 (완전 대체, CLAUDE.md랑 같은 방식).

### 원격(GitHub) 연동 — 진행 중

OAuth 로그인 + GitHub App 웹훅까지 갖춘 클라우드용 흐름. Postgres 필요.

```bash
make up   # docker-compose로 web + db
# 또는 브루/로컬 Postgres + `make dev`
```

`environments/local.env`에 `AUTH_SECRET`, `GITHUB_CLIENT_ID/SECRET` 등 채워야 함 — 각 값 발급 방법은 [docs/github-app-setup.md](./docs/github-app-setup.md) 참고.

### CLI 원격 명령

```bash
repotale login              # GitHub PAT 로그인
repotale connect owner/repo  # repo 연결
repotale open                # 웹 대시보드 열기
```

## 다음

- 원격 repo 분석 파이프라인을 실제 웹훅/연동 화면에 연결
- 이해확인 계기판, 라이선스/과금 모델 — [docs/](./docs)의 미해결 이슈 참고

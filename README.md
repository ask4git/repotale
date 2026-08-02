<p align="center">
  <strong>repotale</strong> — 커밋을 하나씩 읽는 대신, 이야기로 읽습니다
</p>

<p align="center">
  바이브 코딩으로 코드는 빠르게 쌓이지만, 그걸 이해하는 속도는 따라가지 못합니다 — 이른바 "comprehension debt"입니다.<br>
  repotale은 git 히스토리를 <strong>비개발자도 읽을 수 있는 블로그</strong>로 바꿔서 이 빚을 갚아나갑니다.
</p>

---

## 30초 안에 살펴보기

repo 안에서 실행하면 최근 커밋을 읽고, 그 내용을 로컬 웹 화면에 블로그 형태로 보여줍니다.

```bash
go install github.com/ask4git/repotale/repotale@v0.1.0 && export PATH="$PATH:$(go env GOPATH)/bin" && repotale
```

`go install`은 컴파일만 담당하는 Go 툴체인 명령이라 설치 과정 중에는 아무것도 출력하지 않습니다(npm의 postinstall 같은 훅이 Go에는 없습니다). 그래서 설치, PATH 설정, 첫 실행을 한 줄로 이어붙여 설치 직후 바로 아래 화면을 볼 수 있게 했습니다.

<details>
<summary>설치가 "module ... does not contain package" 같은 에러로 실패한다면</summary>

릴리즈 직후라 Go 모듈 프록시나 체크섬 DB가 해당 버전을 아직 인덱싱하지 못했을 가능성이 있습니다. 보통 며칠 안에 자연스럽게 해결됩니다. 급하다면 아래처럼 우회할 수 있습니다.

```bash
GOPROXY=direct GOSUMDB=off go install github.com/ask4git/repotale/repotale@v0.1.0
```

`GOSUMDB=off`는 체크섬 검증을 끄는 옵션이므로 임시 우회용으로만 사용하고, 평소에는 기본 명령을 그대로 쓰는 것을 권장합니다.

</details>

이후 새 버전으로 업데이트할 때는 `repotale update`만 실행하면 됩니다.

PATH 설정을 매번 새 터미널에서도 유지하고 싶다면 한 번만 실행해두면 됩니다.

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

필요한 건 로그인된 로컬 `claude` CLI뿐입니다. API 키도, Node.js도, 별도 로그인 화면도 필요하지 않습니다. `repotale` 바이너리 하나가 분석부터 웹 서버 구동까지 모두 처리합니다. 다시 실행해도 이미 분석한 커밋은 건너뛰므로 같은 커밋을 두 번 분석하느라 비용이 낭비되지 않습니다.

언어 선택은 최초 실행 시 한 번만 물어보고 `~/.repotale/settings`에 저장됩니다. `.zshrc`처럼 사람이 직접 열어서 고칠 수 있는 평문 파일입니다. 나중에 바꾸고 싶다면 `repotale --language`를 실행하면 됩니다. UI 언어와는 별개로 **생성되는 글의 언어**도 실행할 때마다 따로 선택할 수 있으며, 기본값은 UI 언어를 따라갑니다.

분석을 시작하기 전에는 diff 크기를 기준으로 대략적인 토큰 사용량을 먼저 보여주고 진행 여부를 확인합니다. 모르는 사이에 비용이 나가는 일을 막기 위해서입니다.

### 실제로 이 repo를 분석하면 나오는 글

아래는 실제 결과이며, 따로 손질하지 않았습니다.

> **GitHub 하나뿐이던 로그인 버튼에 Google·GitLab·Apple·Passkey를 붙였습니다**
>
> 오늘 로그인 화면을 뜯어고쳤습니다. 지금까지는 GitHub 계정이 없으면 아예 들어올 방법이 없었는데, 이제 Google, GitLab, Apple로도 로그인할 수 있고 비밀번호 없이 지문·얼굴로 여는 Passkey, 회사 계정으로 한 번에 들어오는 SSO까지 열어뒀습니다. 어떤 걸로 들어오든 이메일이 같으면 자동으로 한 사람의 계정으로 합쳐지도록 해서, "예전에 뭘로 가입했더라" 하고 헤매는 일이 없게 했고요.
>
> `auth` `feature` `refactor`

commit 메시지나 diff 어디에도 이런 문장은 없습니다. 이걸 읽고 재구성하는 것이 repotale이 하는 일입니다.

## 왜 만들었나

- 바이브 코딩은 사람이 이해하는 속도보다 빠르게 진행됩니다. 코드는 존재하지만 "왜 이렇게 만들었는지" 설명할 수 있는 사람은 없는 상태가 쌓여갑니다.
- commit log와 changelog는 개발자를 위한 포맷입니다. Product Owner, 비개발자 창업자, 갓 합류한 개발자는 읽기 어렵습니다.
- repotale은 What(무슨 기능인지), Where(어디에 들어갔는지), How(어떻게 구현했는지), Why(왜 그렇게 했는지) 네 가지 관점으로 커밋을 재구성해서, 코드를 몰라도 읽을 수 있는 이야기로 바꿔줍니다.

기획 배경 전체는 [repo-story-idea.md](./repo-story-idea.md)에서, 인지부채 개념은 [docs/comprehension-debt-concept.md](./docs/comprehension-debt-concept.md)에서 확인할 수 있습니다.

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
| 로컬 repo 분석 (`repotale` 실행 → 커밋 → 블로그) | 완료 |
| 내장 웹서버로 결과 확인 (Node/DB/로그인 불필요, `repotale` 바이너리 하나로) | 완료 |
| 이미 분석한 커밋 재분석 건너뛰기 | 완료 |
| 분석 전 예상 토큰 사용량 확인 및 진행 여부 확인 | 완료 |
| 분석 톤 프리셋 (soft/medium/hard) 및 `.repotale/prompt.md` 완전 커스텀 | 완료 |
| UI 언어(en/ko) 선택 및 저장, 생성되는 글의 언어 별도 선택 | 완료 |
| GitHub/Google/GitLab/Apple 로그인 (원격 연동용) | 완료 |
| GitHub repo 연동 화면 (`/connect`) | 완료 |
| push 시 자동 재분석 (GitHub 웹훅) | 웹훅 수신은 되지만 분석 파이프라인과 아직 연결되지 않음 |
| 원격 repo(PR 기반) 분석 결과를 웹에서 확인 | 미완성 — `/repo/[owner]/[repo]`는 아직 데모 데이터로 표시됩니다 |

정리하면, **로컬 모드는 지금 바로 사용할 수 있고**, **원격(GitHub) 모드는 로그인과 연동까지만 되어 있으며 분석과 발행은 아직 준비 중**입니다.

## 설치와 실행

### 로컬 분석 (지금 사용할 수 있는 기능)

[Claude Code](https://claude.com/claude-code)를 설치하고 로그인만 되어 있으면 됩니다. 별도 API 키는 필요하지 않습니다.

```bash
go install github.com/ask4git/repotale/repotale@v0.1.0  # 또는: make cli-build
repotale            # 분석하고 싶은 repo 안에서 실행
```

결과는 `~/.repotale/<repo이름>/posts/*.json`에 저장되며, 분석이 끝나면 `repotale`이 자동으로 `http://127.0.0.1:4321`에 웹 서버를 띄워 보여줍니다.

톤 프리셋 3종(soft/medium/hard)은 `~/.repotale/presets/*.md`에 실제 파일로 생성되어 직접 열어 수정할 수 있습니다. 최초 한 번만 생성되며 이후에는 덮어쓰지 않습니다. repo에 `.repotale/prompt.md` 파일을 두면 프리셋 선택 과정 자체를 건너뛰고 그 내용을 그대로 사용합니다. CLAUDE.md와 같은 방식의 완전한 커스텀 오버라이드입니다.

### 원격(GitHub) 연동 (진행 중인 기능)

OAuth 로그인과 GitHub App 웹훅을 갖춘 클라우드용 흐름이며, Postgres가 필요합니다.

```bash
make up   # docker-compose로 web + db 실행
# 또는 로컬 Postgres 설치 후 `make dev`
```

`environments/local.env`에 `AUTH_SECRET`, `GITHUB_CLIENT_ID/SECRET` 등을 채워야 합니다. 각 값을 발급받는 방법은 [docs/github-app-setup.md](./docs/github-app-setup.md)에서 확인할 수 있습니다.

### CLI 원격 명령

```bash
repotale login              # GitHub PAT 로그인
repotale connect owner/repo  # repo 연결
repotale open                # 웹 대시보드 열기
```

## 다음 계획

- 원격 repo 분석 파이프라인을 실제 웹훅/연동 화면에 연결합니다.
- 이해확인 계기판, 라이선스와 과금 모델은 [docs/](./docs)의 미해결 이슈 문서를 참고해주세요.

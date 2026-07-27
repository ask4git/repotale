# GitHub App 설정 (수동)

> repo 연결 이후 push마다 웹훅으로 재분석을 트리거하기 위한 GitHub App.
> 로그인에 쓰는 GitHub OAuth App(`GITHUB_CLIENT_ID`/`GITHUB_CLIENT_SECRET`)과는
> 별개다. 아래는 github.com에서 사람이 직접 해야 하는 등록 절차라 코드로
> 자동화할 수 없다.

## 1. App 생성

GitHub → Settings → Developer settings → GitHub Apps → New GitHub App

- **Webhook URL**: `https://<배포 도메인>/api/webhooks/github`
  (로컬 개발은 ngrok 등으로 터널링한 URL)
- **Webhook secret**: 임의의 랜덤 문자열 생성 후 저장 (아래 3번 참고)

## 2. 권한 (Permissions)

Repository permissions에서 아래 두 개만 Read-only로 설정한다.

- **Contents**: Read-only
- **Metadata**: Read-only (자동 필수 항목)

## 3. Webhook 이벤트

Subscribe to events에서 다음 하나만 체크한다.

- **Push**

## 4. 값 발급 및 환경변수 반영

App 생성 후 아래 세 값을 `environments/*.env` (각 환경별 실제 값 파일,
`.env.example`이 아님)에 채운다.

| App 설정 화면 값 | 환경변수 |
| --- | --- |
| App ID | `GITHUB_APP_ID` |
| Generate a private key로 받은 `.pem` 파일 내용 | `GITHUB_APP_PRIVATE_KEY` |
| 1번에서 설정한 Webhook secret | `GITHUB_WEBHOOK_SECRET` |

`GITHUB_APP_PRIVATE_KEY`는 PEM 전체(줄바꿈 포함)를 한 줄로 넣어야 하므로
줄바꿈을 `\n`으로 이스케이프해서 저장한다.

## 5. 설치 (Install)

App을 만든 계정/조직에서 대상 repo에 App을 install해야 웹훅이 실제로
발생한다. Install 시 생성되는 `installation_id`는 이후 설치 토큰 발급
(`getInstallationAccessToken`, `web/lib/github-app.ts`)에 필요하며, 어느
repo가 어느 installation에 연결됐는지 저장하는 부분은 이 문서의 범위 밖이다.

import { auth } from "@/auth";
import { redirect } from "next/navigation";
import { headers } from "next/headers";
import { TopNav } from "@/components/top-nav";
import { RepoList } from "./repo-list";
import { LinkGithubButton } from "./link-github";

async function fetchRepos(accessToken: string) {
  const res = await fetch(
    "https://api.github.com/user/repos?sort=updated&per_page=20",
    { headers: { Authorization: `Bearer ${accessToken}` } },
  );
  if (!res.ok) throw new Error(`GitHub API error: ${res.status}`);
  return res.json();
}

export default async function ConnectPage() {
  const requestHeaders = await headers();
  const session = await auth.api.getSession({ headers: requestHeaders });
  if (!session) redirect("/login");

  // repo 목록을 읽으려면 GitHub 계정이 연결돼 있어야 함 - 다른 provider로
  // 로그인한 경우 아직 연결 전일 수 있으므로 여기서 그 상태를 감지한다.
  let accessToken: string | undefined;
  try {
    const result = await auth.api.getAccessToken({
      body: { providerId: "github", userId: session.user.id },
      headers: requestHeaders,
    });
    accessToken = result.accessToken;
  } catch {
    accessToken = undefined;
  }

  return (
    <>
      <TopNav userName={session.user?.name} userImage={session.user?.image} />
      <main className="mx-auto max-w-3xl px-6 py-12">
        <h1 className="text-xl font-semibold tracking-tight">
          연동할 저장소를 선택하세요
        </h1>
        <p className="mt-2 text-sm text-muted-foreground">
          연동한 repo는 커밋이 쌓일 때마다 새 글이 자동으로 생성됩니다.
        </p>
        <div className="mt-8">
          {accessToken ? (
            <RepoList repos={await fetchRepos(accessToken)} />
          ) : (
            <div className="flex flex-col items-start gap-3 rounded-lg border p-6">
              <p className="text-sm text-muted-foreground">
                repo 목록을 불러오려면 GitHub 계정 연결이 필요합니다.
              </p>
              <LinkGithubButton />
            </div>
          )}
        </div>
      </main>
    </>
  );
}

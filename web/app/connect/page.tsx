import { auth } from "@/auth";
import { redirect } from "next/navigation";
import { headers } from "next/headers";
import { TopNav } from "@/components/top-nav";
import { RepoList } from "./repo-list";

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

  const { accessToken } = await auth.api.getAccessToken({
    body: { providerId: "github", userId: session.user.id },
    headers: requestHeaders,
  });
  const repos = await fetchRepos(accessToken!);

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
          <RepoList repos={repos} />
        </div>
      </main>
    </>
  );
}

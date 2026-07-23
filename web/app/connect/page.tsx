import { auth, signOut } from "@/auth";
import { redirect } from "next/navigation";
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
  const session = await auth();
  if (!session) redirect("/login");

  const repos = await fetchRepos(session.accessToken!);

  return (
    <main>
      <h1>연동할 repo를 선택하세요</h1>
      <p>{session.user?.name ?? session.user?.email}로 로그인됨</p>
      <form
        action={async () => {
          "use server";
          await signOut({ redirectTo: "/login" });
        }}
      >
        <button type="submit">로그아웃</button>
      </form>
      <RepoList repos={repos} />
    </main>
  );
}

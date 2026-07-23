import { auth, signIn } from "@/auth";
import { redirect } from "next/navigation";

export default async function LoginPage() {
  const session = await auth();
  if (session) redirect("/connect");

  return (
    <main>
      <h1>repotale</h1>
      <p>GitHub repo를 연동하면, 무엇이 왜 바뀌었는지 블로그로 정리해드립니다.</p>
      <form
        action={async () => {
          "use server";
          await signIn("github", { redirectTo: "/connect" });
        }}
      >
        <button type="submit">GitHub로 로그인</button>
      </form>
    </main>
  );
}

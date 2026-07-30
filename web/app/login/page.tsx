import { auth } from "@/auth";
import { redirect } from "next/navigation";
import { headers } from "next/headers";
import { LoginButtons } from "./login-buttons";

export default async function LoginPage() {
  const session = await auth.api.getSession({ headers: await headers() });
  if (session) redirect("/connect");

  return (
    <main className="flex min-h-screen flex-col items-center justify-center px-6">
      <div className="w-full max-w-sm">
        <div className="mb-10 text-center">
          <p className="text-sm font-medium tracking-tight text-muted-foreground">
            repotale
          </p>
          <h1 className="mt-3 text-2xl font-semibold tracking-tight">
            내 프로덕트가 어떻게<br />만들어졌는지 읽어보세요
          </h1>
          <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
            GitHub repo를 연동하면 커밋 히스토리를 읽기 쉬운
            블로그 형식으로 정리해드립니다.
          </p>
        </div>

        <LoginButtons />

        <p className="mt-6 text-center text-xs text-muted-foreground">
          로그인하면 repo의 커밋 히스토리를 읽기 위한 최소 권한만 요청합니다.
          같은 이메일의 다른 계정은 자동으로 연동됩니다.
        </p>
      </div>
    </main>
  );
}

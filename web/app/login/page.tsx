import { auth, signIn } from "@/auth";
import { redirect } from "next/navigation";
import { Button } from "@/components/ui/button";

function GithubMark() {
  return (
    <svg viewBox="0 0 24 24" className="size-4" fill="currentColor" aria-hidden>
      <path d="M12 .5C5.65.5.5 5.65.5 12c0 5.08 3.29 9.39 7.86 10.91.57.1.78-.25.78-.55 0-.27-.01-1-.02-1.96-3.2.7-3.88-1.54-3.88-1.54-.52-1.33-1.28-1.68-1.28-1.68-1.04-.71.08-.7.08-.7 1.15.08 1.76 1.18 1.76 1.18 1.03 1.75 2.7 1.25 3.36.96.1-.75.4-1.25.73-1.53-2.55-.29-5.24-1.28-5.24-5.68 0-1.25.45-2.28 1.18-3.08-.12-.29-.51-1.46.11-3.04 0 0 .96-.31 3.15 1.18a10.9 10.9 0 0 1 5.74 0c2.19-1.49 3.15-1.18 3.15-1.18.62 1.58.23 2.75.11 3.04.74.8 1.18 1.83 1.18 3.08 0 4.41-2.7 5.38-5.27 5.67.42.36.78 1.07.78 2.15 0 1.56-.01 2.81-.01 3.19 0 .3.2.66.79.55A10.52 10.52 0 0 0 23.5 12C23.5 5.65 18.35.5 12 .5Z" />
    </svg>
  );
}

export default async function LoginPage() {
  const session = await auth();
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

        <form
          action={async () => {
            "use server";
            await signIn("github", { redirectTo: "/connect" });
          }}
        >
          <Button type="submit" size="lg" className="w-full gap-2 rounded-full">
            <GithubMark />
            GitHub로 로그인
          </Button>
        </form>

        <p className="mt-6 text-center text-xs text-muted-foreground">
          로그인하면 repo의 커밋 히스토리를 읽기 위한 최소 권한만 요청합니다.
        </p>
      </div>
    </main>
  );
}

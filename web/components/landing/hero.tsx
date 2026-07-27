import Link from "next/link";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

function GithubMark() {
  return (
    <svg viewBox="0 0 24 24" className="size-4" fill="currentColor" aria-hidden>
      <path d="M12 .5C5.65.5.5 5.65.5 12c0 5.08 3.29 9.39 7.86 10.91.57.1.78-.25.78-.55 0-.27-.01-1-.02-1.96-3.2.7-3.88-1.54-3.88-1.54-.52-1.33-1.28-1.68-1.28-1.68-1.04-.71.08-.7.08-.7 1.15.08 1.76 1.18 1.76 1.18 1.03 1.75 2.7 1.25 3.36.96.1-.75.4-1.25.73-1.53-2.55-.29-5.24-1.28-5.24-5.68 0-1.25.45-2.28 1.18-3.08-.12-.29-.51-1.46.11-3.04 0 0 .96-.31 3.15 1.18a10.9 10.9 0 0 1 5.74 0c2.19-1.49 3.15-1.18 3.15-1.18.62 1.58.23 2.75.11 3.04.74.8 1.18 1.83 1.18 3.08 0 4.41-2.7 5.38-5.27 5.67.42.36.78 1.07.78 2.15 0 1.56-.01 2.81-.01 3.19 0 .3.2.66.79.55A10.52 10.52 0 0 0 23.5 12C23.5 5.65 18.35.5 12 .5Z" />
    </svg>
  );
}

export function Hero() {
  return (
    <section className="px-6 py-24 text-center">
      <div className="mx-auto max-w-2xl">
        <Badge variant="secondary">repotale</Badge>
        <h1 className="mt-4 text-4xl font-semibold tracking-tight text-balance">
          코드는 AI가 앞질렀고,
          <br />
          이해는 사람이 뒤처졌습니다
        </h1>
        <p className="mt-4 text-base leading-relaxed text-muted-foreground text-balance">
          예전 부채가 '나중에 갚을 코드'였다면, 지금 부채는 '이미 있는데
          아무도 모르는 코드'입니다. repotale은 GitHub repo를 연동해 AI agent가
          만든 기능을 읽을 수 있는 blog post로 자동 발행하고, 그 인지부채를
          갚아드립니다.
        </p>
        <div className="mt-8 flex justify-center">
          <Button
            render={<Link href="/login" />}
            size="lg"
            className="gap-2 rounded-full"
          >
            <GithubMark />
            GitHub로 시작하기
          </Button>
        </div>
      </div>
    </section>
  );
}

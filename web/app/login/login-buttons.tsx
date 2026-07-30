"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { authClient } from "@/lib/auth-client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { KeyRound, Mail } from "lucide-react";

// 각 브랜드의 공식 로고마크를 근사한 SVG. 픽셀 단위로 브랜드 가이드라인을
// 검수한 리소스는 아니라서, 대외 공개 전엔 각 사 press kit의 공식 에셋으로
// 교체하는 걸 권장함.
function GithubMark() {
  return (
    <svg viewBox="0 0 24 24" className="size-[18px]" fill="currentColor" aria-hidden>
      <path d="M12 .5C5.65.5.5 5.65.5 12c0 5.08 3.29 9.39 7.86 10.91.57.1.78-.25.78-.55 0-.27-.01-1-.02-1.96-3.2.7-3.88-1.54-3.88-1.54-.52-1.33-1.28-1.68-1.28-1.68-1.04-.71.08-.7.08-.7 1.15.08 1.76 1.18 1.76 1.18 1.03 1.75 2.7 1.25 3.36.96.1-.75.4-1.25.73-1.53-2.55-.29-5.24-1.28-5.24-5.68 0-1.25.45-2.28 1.18-3.08-.12-.29-.51-1.46.11-3.04 0 0 .96-.31 3.15 1.18a10.9 10.9 0 0 1 5.74 0c2.19-1.49 3.15-1.18 3.15-1.18.62 1.58.23 2.75.11 3.04.74.8 1.18 1.83 1.18 3.08 0 4.41-2.7 5.38-5.27 5.67.42.36.78 1.07.78 2.15 0 1.56-.01 2.81-.01 3.19 0 .3.2.66.79.55A10.52 10.52 0 0 0 23.5 12C23.5 5.65 18.35.5 12 .5Z" />
    </svg>
  );
}

function GoogleMark() {
  return (
    <svg viewBox="0 0 48 48" className="size-[18px]" aria-hidden>
      <path fill="#FFC107" d="M43.6 20.5H24v8h11.3c-1.6 4.7-6.1 8-11.3 8-6.6 0-12-5.4-12-12s5.4-12 12-12c3.1 0 5.8 1.1 8 3l6-6C34.9 5.1 29.7 3 24 3 12.4 3 3 12.4 3 24s9.4 21 21 21 21-9.4 21-21c0-1.4-.1-2.7-.4-3.5z" />
      <path fill="#FF3D00" d="M6.3 14.7l6.6 4.8C14.6 15.9 18.9 13 24 13c3.1 0 5.8 1.1 8 3l6-6C34.9 5.1 29.7 3 24 3 16.3 3 9.6 7.3 6.3 14.7z" />
      <path fill="#4CAF50" d="M24 45c5.6 0 10.6-2.1 14.4-5.6l-6.6-5.6C29.6 35.9 26.9 37 24 37c-5.2 0-9.6-3.3-11.3-8l-6.6 5.1C9.6 40.7 16.3 45 24 45z" />
      <path fill="#1976D2" d="M43.6 20.5H24v8h11.3c-.8 2.3-2.2 4.2-4.1 5.6l6.6 5.6C41.9 36.5 45 30.9 45 24c0-1.4-.1-2.7-.4-3.5z" />
    </svg>
  );
}

function GitlabMark() {
  return (
    <svg viewBox="0 0 24 24" className="size-[18px]" aria-hidden>
      <path fill="#e24329" d="M12 21.94 15.68 10.6H8.32L12 21.94z" />
      <path fill="#fc6d26" d="M12 21.94 8.32 10.6H2.98L12 21.94z" />
      <path fill="#fca326" d="M2.98 10.6 1.63 14.7a.85.85 0 0 0 .31.95L12 21.94 2.98 10.6z" />
      <path fill="#e24329" d="M2.98 10.6h5.34L6.03 4.14a.42.42 0 0 0-.8 0L2.98 10.6z" />
      <path fill="#fc6d26" d="M12 21.94 15.68 10.6h5.34L12 21.94z" />
      <path fill="#fca326" d="M21.02 10.6 22.37 14.7a.85.85 0 0 1-.31.95L12 21.94l9.02-11.34z" />
      <path fill="#e24329" d="M21.02 10.6h-5.34l2.29-6.46a.42.42 0 0 1 .8 0l2.25 6.46z" />
    </svg>
  );
}

function AppleMark() {
  return (
    <svg viewBox="0 0 24 24" className="size-[18px]" fill="currentColor" aria-hidden>
      <path d="M17.05 12.63c-.03-2.5 2.04-3.7 2.13-3.76-1.16-1.7-2.97-1.93-3.62-1.96-1.54-.16-3 .9-3.78.9-.78 0-1.98-.88-3.26-.86-1.68.03-3.23.98-4.09 2.48-1.75 3.03-.45 7.5 1.25 9.96.83 1.2 1.82 2.55 3.12 2.5 1.25-.05 1.72-.8 3.23-.8 1.5 0 1.93.8 3.25.78 1.34-.02 2.19-1.22 3.01-2.42.95-1.39 1.34-2.73 1.36-2.8-.03-.01-2.6-1-2.63-3.98z" />
      <path d="M14.86 5.13c.68-.83 1.15-1.98 1.02-3.13-.99.04-2.18.66-2.89 1.48-.63.73-1.19 1.9-1.04 3.02 1.09.09 2.21-.55 2.91-1.37z" />
    </svg>
  );
}

const SOCIAL_PROVIDERS: {
  id: "github" | "google" | "gitlab" | "apple";
  label: string;
  icon: React.ReactNode;
  className: string;
}[] = [
  {
    id: "github",
    label: "GitHub로 로그인",
    icon: <GithubMark />,
    className: "border-transparent bg-[#181717] text-white hover:bg-[#2b2b2b]",
  },
  {
    id: "google",
    label: "Google로 로그인",
    icon: <GoogleMark />,
    className:
      "border border-[#dadce0] bg-white text-[#3c4043] hover:bg-[#f8f9fa] dark:border-[#5f6368] dark:bg-[#131314] dark:text-[#e8eaed] dark:hover:bg-[#1e1f20]",
  },
  {
    id: "gitlab",
    label: "GitLab로 로그인",
    icon: <GitlabMark />,
    className:
      "border border-[#dadce0] bg-white text-[#171321] hover:bg-[#f8f9fa] dark:border-[#3a3a3a] dark:bg-[#171321] dark:text-white dark:hover:bg-[#241f33]",
  },
  {
    id: "apple",
    label: "Apple로 로그인",
    icon: <AppleMark />,
    className: "border-transparent bg-black text-white hover:bg-black/90 dark:bg-white dark:text-black dark:hover:bg-white/90",
  },
];

export function LoginButtons() {
  const router = useRouter();
  const [ssoEmail, setSsoEmail] = useState("");
  const [pending, setPending] = useState<string | null>(null);

  async function signInSocial(provider: "github" | "google" | "gitlab" | "apple") {
    setPending(provider);
    await authClient.signIn.social({ provider, callbackURL: "/connect" });
  }

  async function signInPasskey() {
    setPending("passkey");
    await authClient.signIn.passkey({
      fetchOptions: { onSuccess: () => router.push("/connect") },
    });
    setPending(null);
  }

  async function signInSso() {
    if (!ssoEmail) return;
    setPending("sso");
    await authClient.signIn.sso({ email: ssoEmail, callbackURL: "/connect" });
  }

  return (
    <div className="flex flex-col gap-3">
      {SOCIAL_PROVIDERS.map((provider) => (
        <Button
          key={provider.id}
          size="lg"
          variant="outline"
          className={`w-full gap-2 rounded-full ${provider.className}`}
          disabled={pending !== null}
          onClick={() => signInSocial(provider.id)}
        >
          {provider.icon}
          {provider.label}
        </Button>
      ))}

      <Button
        size="lg"
        variant="outline"
        className="w-full gap-2 rounded-full"
        disabled={pending !== null}
        onClick={signInPasskey}
      >
        <KeyRound className="size-4" />
        Passkey로 로그인
      </Button>

      <div className="mt-2 flex items-center gap-2">
        <Input
          type="email"
          placeholder="회사 이메일로 SSO 로그인"
          value={ssoEmail}
          onChange={(e) => setSsoEmail(e.target.value)}
          disabled={pending !== null}
        />
        <Button
          variant="secondary"
          size="icon"
          disabled={pending !== null || !ssoEmail}
          onClick={signInSso}
          aria-label="SSO로 로그인"
        >
          <Mail className="size-4" />
        </Button>
      </div>
    </div>
  );
}

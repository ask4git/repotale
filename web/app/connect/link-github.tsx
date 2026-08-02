"use client";

import { authClient } from "@/lib/auth-client";
import { Button } from "@/components/ui/button";

export function LinkGithubButton() {
  return (
    <Button
      size="lg"
      className="gap-2 rounded-full"
      onClick={() =>
        authClient.linkSocial({ provider: "github", callbackURL: "/connect" })
      }
    >
      GitHub 계정 연결하기
    </Button>
  );
}

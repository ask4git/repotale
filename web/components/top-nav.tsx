"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { authClient } from "@/lib/auth-client";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";

export function TopNav({
  userName,
  userImage,
}: {
  userName?: string | null;
  userImage?: string | null;
}) {
  const router = useRouter();

  return (
    <header className="border-b">
      <div className="mx-auto flex h-14 max-w-3xl items-center justify-between px-6">
        <Link href="/connect" className="text-sm font-semibold tracking-tight">
          repotale
        </Link>
        <div className="flex items-center gap-3">
          <Avatar className="size-7">
            <AvatarImage src={userImage ?? undefined} alt={userName ?? ""} />
            <AvatarFallback className="text-xs">
              {userName?.[0]?.toUpperCase() ?? "?"}
            </AvatarFallback>
          </Avatar>
          <Button
            variant="ghost"
            size="sm"
            onClick={() =>
              authClient.signOut({
                fetchOptions: { onSuccess: () => router.push("/login") },
              })
            }
          >
            로그아웃
          </Button>
        </div>
      </div>
    </header>
  );
}

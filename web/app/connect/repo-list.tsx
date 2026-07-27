"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Lock, Search } from "lucide-react";

type Repo = {
  id: number;
  full_name: string;
  private: boolean;
  description: string | null;
};

export function RepoList({ repos }: { repos: Repo[] }) {
  const router = useRouter();
  const [query, setQuery] = useState("");
  // ponytail: connection state is client-only, no DB write yet. Persist once
  // a repos table exists.
  const [connected, setConnected] = useState<number | null>(null);

  const filtered = repos.filter((repo) =>
    repo.full_name.toLowerCase().includes(query.toLowerCase()),
  );

  return (
    <div>
      <div className="relative mb-4">
        <Search className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          placeholder="저장소 검색"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="pl-9"
        />
      </div>

      <ul className="divide-y rounded-lg border">
        {filtered.map((repo) => {
          const [owner, name] = repo.full_name.split("/");
          const isConnected = connected === repo.id;
          return (
            <li
              key={repo.id}
              className="flex items-center justify-between gap-4 px-4 py-3.5"
            >
              <div className="min-w-0">
                <div className="flex items-center gap-2">
                  <span className="truncate text-sm font-medium">
                    {repo.full_name}
                  </span>
                  {repo.private && (
                    <Badge variant="secondary" className="gap-1 text-xs">
                      <Lock className="size-3" />
                      private
                    </Badge>
                  )}
                </div>
                {repo.description && (
                  <p className="mt-0.5 truncate text-sm text-muted-foreground">
                    {repo.description}
                  </p>
                )}
              </div>
              <Button
                size="sm"
                variant={isConnected ? "secondary" : "outline"}
                className="shrink-0 rounded-full"
                onClick={() => {
                  setConnected(repo.id);
                  router.push(`/repo/${owner}/${name}`);
                }}
              >
                {isConnected ? "연결됨" : "연결"}
              </Button>
            </li>
          );
        })}
        {filtered.length === 0 && (
          <li className="px-4 py-8 text-center text-sm text-muted-foreground">
            일치하는 저장소가 없습니다.
          </li>
        )}
      </ul>
    </div>
  );
}

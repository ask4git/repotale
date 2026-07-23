"use client";

import { useState } from "react";

type Repo = {
  id: number;
  full_name: string;
  private: boolean;
  description: string | null;
};

export function RepoList({ repos }: { repos: Repo[] }) {
  // ponytail: connection state is client-only, no DB write yet. Persist once
  // a repos table exists.
  const [connected, setConnected] = useState<number | null>(null);

  return (
    <ul>
      {repos.map((repo) => (
        <li key={repo.id}>
          <span>{repo.full_name}</span>
          {repo.private && <span> (private)</span>}
          {repo.description && <p>{repo.description}</p>}
          <button
            type="button"
            disabled={connected === repo.id}
            onClick={() => setConnected(repo.id)}
          >
            {connected === repo.id ? "연결됨" : "연결"}
          </button>
        </li>
      ))}
    </ul>
  );
}

import Link from "next/link";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Sparkles } from "lucide-react";
import { getLocalPosts } from "@/lib/local-posts";

export default async function LocalRepoPage({
  params,
}: {
  params: Promise<{ repoId: string }>;
}) {
  const { repoId } = await params;
  const posts = await getLocalPosts(repoId);

  return (
    <>
      <header className="border-b">
        <div className="mx-auto flex h-14 max-w-3xl items-center px-6">
          <Link href="/" className="text-sm font-semibold tracking-tight">
            repotale <span className="text-muted-foreground">(local)</span>
          </Link>
        </div>
      </header>
      <main className="mx-auto max-w-3xl px-6 py-12">
        <p className="text-sm text-muted-foreground">{repoId}</p>
        <h1 className="mt-1 text-xl font-semibold tracking-tight">
          이 repo에서 무슨 일이 있었나
        </h1>

        {posts.length === 0 ? (
          <Card className="mt-8">
            <CardContent className="flex flex-col items-center gap-3 py-12 text-center">
              <Sparkles className="size-5 text-muted-foreground" />
              <p className="text-sm font-medium">아직 분석 결과가 없습니다</p>
              <p className="max-w-xs text-sm text-muted-foreground">
                이 repo 안에서 <code>repotale</code>을 실행하면 여기 나타납니다.
              </p>
            </CardContent>
          </Card>
        ) : (
          <ul className="mt-8 space-y-4">
            {posts.map((post) => (
              <li key={post.slug}>
                <Card>
                  <CardHeader>
                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                      <time>{post.publishedAt}</time>
                      {post.tags.map((tag) => (
                        <Badge key={tag} variant="secondary" className="text-xs">
                          {tag}
                        </Badge>
                      ))}
                    </div>
                    <CardTitle className="text-base leading-snug font-semibold">
                      {post.title}
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm leading-relaxed text-muted-foreground">
                      {post.excerpt}
                    </p>
                    <p className="mt-3 text-xs text-muted-foreground">
                      commit {post.commitSha.slice(0, 7)}
                    </p>
                  </CardContent>
                </Card>
              </li>
            ))}
          </ul>
        )}
      </main>
    </>
  );
}

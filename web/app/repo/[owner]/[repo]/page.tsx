import { auth } from "@/auth";
import { redirect } from "next/navigation";
import Link from "next/link";
import { TopNav } from "@/components/top-nav";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ArrowUpRight, Sparkles } from "lucide-react";
import { demoPosts } from "./posts";

export default async function RepoSummaryPage({
  params,
}: {
  params: Promise<{ owner: string; repo: string }>;
}) {
  const session = await auth();
  if (!session) redirect("/login");

  const { owner, repo } = await params;
  const fullName = `${owner}/${repo}`;
  const posts = demoPosts[fullName];

  return (
    <>
      <TopNav userName={session.user?.name} userImage={session.user?.image} />
      <main className="mx-auto max-w-3xl px-6 py-12">
        <p className="text-sm text-muted-foreground">{fullName}</p>
        <h1 className="mt-1 text-xl font-semibold tracking-tight">
          이 repo에서 무슨 일이 있었나
        </h1>

        {!posts ? (
          <Card className="mt-8">
            <CardContent className="flex flex-col items-center gap-3 py-12 text-center">
              <Sparkles className="size-5 text-muted-foreground" />
              <p className="text-sm font-medium">아직 분석 결과가 없습니다</p>
              <p className="max-w-xs text-sm text-muted-foreground">
                최초 분석에는 커밋 히스토리 규모에 따라 몇 분 정도 걸릴 수
                있어요. 완료되면 이 페이지에 첫 글이 나타납니다.
              </p>
            </CardContent>
          </Card>
        ) : (
          <ul className="mt-8 space-y-4">
            {posts.map((post) => (
              <li key={post.slug}>
                <Card className="transition-colors hover:border-foreground/20">
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
                    <Link
                      href={post.prUrl}
                      target="_blank"
                      rel="noreferrer"
                      className="mt-3 inline-flex items-center gap-1 text-sm text-blue-600 hover:underline dark:text-blue-400"
                    >
                      PR #{post.prNumber} 보기
                      <ArrowUpRight className="size-3.5" />
                    </Link>
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

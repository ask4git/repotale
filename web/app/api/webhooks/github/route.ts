import { NextRequest, NextResponse } from "next/server";
import { verifyWebhookSignature } from "@/lib/github-app";

export async function POST(req: NextRequest) {
  // must read the raw body before any JSON parsing, or the HMAC won't match
  const rawBody = await req.text();
  const signature = req.headers.get("x-hub-signature-256");
  const secret = process.env.GITHUB_WEBHOOK_SECRET;

  if (!secret || !verifyWebhookSignature(rawBody, signature, secret)) {
    return NextResponse.json({ error: "invalid signature" }, { status: 401 });
  }

  const event = req.headers.get("x-github-event");
  const payload = JSON.parse(rawBody);

  if (event === "push") {
    const fullName = payload.repository?.full_name;
    const commits = payload.commits ?? [];
    console.log(`[webhook] push received for ${fullName}: ${commits.length} commits`);
    // TODO(pipeline): trigger analyzeRepo() here once web/lib/pipeline/ lands
  }

  return NextResponse.json({ ok: true });
}

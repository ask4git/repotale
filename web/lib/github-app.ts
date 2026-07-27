import crypto from "crypto";
import { SignJWT, importPKCS8 } from "jose";

// Verifies a GitHub webhook payload against the `X-Hub-Signature-256` header.
// GitHub signs the raw request body with HMAC-SHA256 using the webhook secret.
export function verifyWebhookSignature(
  rawBody: string,
  signatureHeader: string | null,
  secret: string
): boolean {
  if (!signatureHeader) return false;

  const expected =
    "sha256=" + crypto.createHmac("sha256", secret).update(rawBody).digest("hex");
  const expectedBuf = Buffer.from(expected);
  const actualBuf = Buffer.from(signatureHeader);

  if (expectedBuf.length !== actualBuf.length) return false;
  return crypto.timingSafeEqual(expectedBuf, actualBuf);
}

async function signAppJwt(): Promise<string> {
  const appId = process.env.GITHUB_APP_ID;
  const privateKey = process.env.GITHUB_APP_PRIVATE_KEY;
  if (!appId || !privateKey) {
    throw new Error("GITHUB_APP_ID / GITHUB_APP_PRIVATE_KEY is not set");
  }

  // env vars store the PEM with literal "\n" escapes since PEM is multi-line.
  const key = await importPKCS8(privateKey.replace(/\\n/g, "\n"), "RS256");
  const now = Math.floor(Date.now() / 1000);

  return new SignJWT({})
    .setProtectedHeader({ alg: "RS256" })
    .setIssuedAt(now - 60) // allow for clock drift, per GitHub App auth docs
    .setExpirationTime(now + 600) // 10 min max
    .setIssuer(appId)
    .sign(key);
}

// Mints a short-lived installation access token for calling the GitHub API
// on behalf of an installation of this GitHub App.
export async function getInstallationAccessToken(
  installationId: string | number
): Promise<string> {
  const jwt = await signAppJwt();

  const res = await fetch(
    `https://api.github.com/app/installations/${installationId}/access_tokens`,
    {
      method: "POST",
      headers: {
        Authorization: `Bearer ${jwt}`,
        Accept: "application/vnd.github+json",
        "X-GitHub-Api-Version": "2022-11-28",
      },
    }
  );

  if (!res.ok) {
    throw new Error(
      `failed to mint installation token: ${res.status} ${await res.text()}`
    );
  }

  const data = (await res.json()) as { token: string };
  return data.token;
}

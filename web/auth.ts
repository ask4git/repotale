import { betterAuth } from "better-auth";
import { sso } from "@better-auth/sso";
import { passkey } from "@better-auth/passkey";
import { Pool } from "pg";

export const auth = betterAuth({
  database: new Pool({ connectionString: process.env.DATABASE_URL }),
  baseURL: process.env.NEXT_PUBLIC_APP_URL,
  secret: process.env.AUTH_SECRET,
  socialProviders: {
    github: {
      clientId: process.env.GITHUB_CLIENT_ID as string,
      clientSecret: process.env.GITHUB_CLIENT_SECRET as string,
      scope: ["read:user", "user:email", "repo"],
    },
    google: {
      clientId: process.env.GOOGLE_CLIENT_ID as string,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET as string,
    },
    gitlab: {
      clientId: process.env.GITLAB_CLIENT_ID as string,
      clientSecret: process.env.GITLAB_CLIENT_SECRET as string,
    },
    apple: {
      clientId: process.env.APPLE_CLIENT_ID as string,
      clientSecret: process.env.APPLE_CLIENT_SECRET as string,
    },
  },
  account: {
    accountLinking: {
      enabled: true,
      // github/google/gitlab/apple all return a verified email, so we trust
      // them to auto-link into the same user without a manual "connect" step.
      trustedProviders: ["github", "google", "gitlab", "apple"],
    },
  },
  plugins: [sso(), passkey()],
});

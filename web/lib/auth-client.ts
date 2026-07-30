import { createAuthClient } from "better-auth/react";
import { ssoClient } from "@better-auth/sso/client";
import { passkeyClient } from "@better-auth/passkey/client";

export const authClient = createAuthClient({
  plugins: [ssoClient(), passkeyClient()],
});

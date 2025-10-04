import type { JWT } from "next-auth/jwt";
import type { Session, User } from "next-auth";
import NextAuth from "next-auth";
import Credentials from "next-auth/providers/credentials";

// Function to refresh the access token
async function refreshAccessToken(token: JWT): Promise<JWT> {
  try {
    console.log("🔄 Attempting to refresh access token...");

    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        refresh_token: token.refresh_token,
      }),
    });

    const data = await response.json();
    console.log("🔍 Refresh response:", {
      success: data.success,
      hasAccessToken: !!data.data?.access_token,
    });

    if (!response.ok || !data.success) {
      console.error("❌ Refresh failed:", data.error || "Unknown error");
      throw new Error(data.error || "Failed to refresh token");
    }

    console.log("✅ Token refreshed successfully");

    return {
      ...token,
      access_token: data.data.access_token,
      // Update the expiry time (15 minutes from now)
      accessTokenExpires: Date.now() + 15 * 60 * 1000,
    };
  } catch (error) {
    console.error("❌ Error refreshing access token:", error);
    return {
      ...token,
      error: "RefreshAccessTokenError",
    };
  }
}

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [
    Credentials({
      credentials: {
        email: { label: "Email", type: "email" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials): Promise<User | null> {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/login`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(credentials),
        });

        if (!res.ok) {
          throw new Error("Login failed");
        }

        const data = await res.json();

        if (!data.success) {
          throw new Error(data.message || "Login failed");
        }

        // Backend returns user and company as separate objects
        return {
          id: data.data.user.id,
          email: data.data.user.email,
          first_name: data.data.user.first_name,
          last_name: data.data.user.last_name,
          role: data.data.user.role,
          company: {
            id: data.data.company.id,
            name: data.data.company.name,
            country: data.data.company.country,
            base_currency: data.data.company.base_currency,
          },
          access_token: data.data.access_token,
          refresh_token: data.data.refresh_token,
        };
      },
    }),
  ],

  session: { strategy: "jwt" },

  pages: { signIn: "/login" },

  callbacks: {
    async jwt({ token, user, account }: { token: JWT; user?: User; account?: any }) {
      // Initial sign in
      if (user) {
        token.id = user.id;
        token.email = user.email;
        token.first_name = user.first_name;
        token.last_name = user.last_name;
        token.role = user.role;
        token.company = user.company;
        token.access_token = user.access_token;
        token.refresh_token = user.refresh_token;
        // Set token expiry (15 minutes from now - adjust based on your backend config)
        token.accessTokenExpires = Date.now() + 15 * 60 * 1000;
        return token;
      }

      // If there's no access token expiry set, don't try to refresh
      if (!token.accessTokenExpires) {
        return token;
      }

      // Return previous token if the access token has not expired yet
      if (Date.now() < (token.accessTokenExpires as number)) {
        return token;
      }

      // Access token has expired, try to refresh it only if we have a refresh token
      if (!token.refresh_token) {
        console.error("No refresh token available");
        return {
          ...token,
          error: "RefreshAccessTokenError",
        };
      }

      console.log("Access token expired, refreshing...");
      return refreshAccessToken(token);
    },

    async session({ session, token }: { session: Session; token: JWT }) {
      // If there was an error refreshing the token, return null to force re-login
      if (token.error) {
        console.error("Token refresh error, session invalid");
        return null as any;
      }

      session.user = {
        id: token.id,
        email: token.email,
        first_name: token.first_name,
        last_name: token.last_name,
        role: token.role,
        company: token.company,
        access_token: token.access_token,
        refresh_token: token.refresh_token,
      };
      session.accessToken = token.access_token; // Add for easy access
      return session;
    },
  },
});

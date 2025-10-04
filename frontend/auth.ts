import type { JWT } from "next-auth/jwt";
import type { Session, User } from "next-auth";
import NextAuth from "next-auth";
import Credentials from "next-auth/providers/credentials";

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
      if (user) {
        token.id = user.id;
        token.email = user.email;
        token.first_name = user.first_name;
        token.last_name = user.last_name;
        token.role = user.role;
        token.company = user.company;
        token.access_token = user.access_token;
        token.refresh_token = user.refresh_token;
      }

      if (account) {
        token.access_token = account.access_token ?? token.access_token;
        token.refresh_token = account.refresh_token ?? token.refresh_token;
      }

      return token;
    },

    async session({ session, token }: { session: Session; token: JWT }) {
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

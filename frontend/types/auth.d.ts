import { DefaultSession, DefaultUser } from "next-auth";
import "next-auth/jwt";

declare module "next-auth" {
  interface User extends DefaultUser {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    role: string;
    company: {
      id: string;
      name: string;
      country: string;
    };
    access_token: string;
    refresh_token: string;
  }

  interface Session extends DefaultSession {
    user: User;
  }
}

declare module "next-auth/jwt" {
  interface JWT {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    role: string;
    company: {
      id: string;
      name: string;
      country: string;
    };
    access_token: string;
    refresh_token: string;
  }
}

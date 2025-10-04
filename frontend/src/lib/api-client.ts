import { auth } from "@/auth";
import { getSession } from "next-auth/react";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export class APIError extends Error {
  constructor(
    public status: number,
    message: string,
    public data?: any
  ) {
    super(message);
    this.name = "APIError";
  }
}

// Client-side function to refresh token
async function refreshToken(refreshToken: string): Promise<string | null> {
  try {
    const response = await fetch(`${API_URL}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    const data = await response.json();

    if (!response.ok || !data.success) {
      return null;
    }

    return data.data.access_token;
  } catch (error) {
    console.error("Error refreshing token:", error);
    return null;
  }
}

async function getAuthHeader() {
  const session = await auth();
  if (!session?.user?.access_token) {
    throw new Error("Not authenticated");
  }
  return {
    Authorization: `Bearer ${session.user.access_token}`,
  };
}

export async function apiClient<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const url = `${API_URL}${endpoint}`;

  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  // Add auth header if not a public endpoint
  if (!endpoint.includes("/auth/")) {
    try {
      const authHeader = await getAuthHeader();
      Object.assign(headers, authHeader);
    } catch (error) {
      // For server-side calls without session
      console.error("Auth error:", error);
    }
  }

  const response = await fetch(url, {
    ...options,
    headers,
  });

  const data = await response.json();

  if (!response.ok) {
    throw new APIError(response.status, data.error || data.message || "Request failed", data);
  }

  return data;
}

// Client-side API client (for use in client components)
// This version handles token refresh automatically on 401 errors
export function createClientAPI(accessToken: string, refreshTokenValue?: string) {
  return async function clientApiClient<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_URL}${endpoint}`;

    const headers: HeadersInit = {
      "Content-Type": "application/json",
      ...options.headers,
    };

    // Add auth header
    if (!endpoint.includes("/auth/")) {
      Object.assign(headers, { Authorization: `Bearer ${accessToken}` });
    }

    let response = await fetch(url, {
      ...options,
      headers,
    });

    // If 401 and we have a refresh token, try to refresh
    if (response.status === 401 && refreshTokenValue && !endpoint.includes("/auth/")) {
      console.log("Token expired, attempting refresh...");
      const newAccessToken = await refreshToken(refreshTokenValue);

      if (newAccessToken) {
        console.log("Token refreshed successfully, retrying request...");
        // Update the session with new token
        const session = await getSession();
        if (session) {
          session.accessToken = newAccessToken;
          session.user.access_token = newAccessToken;
        }

        // Retry the request with new token
        const newHeaders: HeadersInit = {
          "Content-Type": "application/json",
          ...options.headers,
          Authorization: `Bearer ${newAccessToken}`,
        };

        response = await fetch(url, {
          ...options,
          headers: newHeaders,
        });
      } else {
        console.error("Token refresh failed, redirecting to login...");
        // Token refresh failed, redirect to login
        if (typeof window !== "undefined") {
          window.location.href = "/login";
        }
      }
    }

    const data = await response.json();

    if (!response.ok) {
      throw new APIError(response.status, data.error || data.message || "Request failed", data);
    }

    return data;
  };
}

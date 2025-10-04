import { auth } from "@/auth";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export class APIError extends Error {
  constructor(public status: number, message: string, public data?: any) {
    super(message);
    this.name = "APIError";
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

export async function apiClient<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
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
    throw new APIError(
      response.status,
      data.error || data.message || "Request failed",
      data
    );
  }

  return data;
}

// Client-side API client (for use in client components)
export function createClientAPI(accessToken: string) {
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

    const response = await fetch(url, {
      ...options,
      headers,
    });

    const data = await response.json();

    if (!response.ok) {
      throw new APIError(
        response.status,
        data.error || data.message || "Request failed",
        data
      );
    }

    return data;
  };
}

import { API_URL, API_KEY } from "astro:env/server";
import type {
  HomeResponse,
  PostResponse,
  TagResponse,
  CategoryResponse,
  SearchResponse,
} from "./types";

export class ApiClientError extends Error {
  constructor(
    public readonly status: number,
    public readonly body: string,
    message: string,
  ) {
    super(message);
    this.name = "ApiClientError";
  }
}

async function apiFetch<T>(
  path: string,
  init: RequestInit = {},
): Promise<T> {
  const url = `${API_URL}${path}`;
  const res = await fetch(url, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      "X-API-Key": API_KEY,
      ...(init.headers as Record<string, string>),
    },
  });

  if (!res.ok) {
    const body = await res.text();
    throw new ApiClientError(res.status, body, `API ${res.status}: ${path}`);
  }

  return res.json() as Promise<T>;
}

// ---- Read endpoints ----

export function getHome(page = 1): Promise<HomeResponse> {
  return apiFetch<HomeResponse>(`/api/?page=${page}`);
}

export function getPost(slug: string): Promise<PostResponse> {
  return apiFetch<PostResponse>(`/api/post/${encodeURIComponent(slug)}`);
}

export function getTag(slug: string, page = 1): Promise<TagResponse> {
  return apiFetch<TagResponse>(
    `/api/tag/${encodeURIComponent(slug)}?page=${page}`,
  );
}

export function getCategory(slug: string, page = 1): Promise<CategoryResponse> {
  return apiFetch<CategoryResponse>(
    `/api/category/${encodeURIComponent(slug)}?page=${page}`,
  );
}

export function search(q: string): Promise<SearchResponse> {
  return apiFetch<SearchResponse>(`/api/search?q=${encodeURIComponent(q)}`);
}

// ---- Write endpoints (called from Astro API routes) ----

export function postComment(body: {
  message: string;
  post_id: number;
  countMe: number;
  honeypot: string;
}): Promise<unknown> {
  return apiFetch("/api/comment", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

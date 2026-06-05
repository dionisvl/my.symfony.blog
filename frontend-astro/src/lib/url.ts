import { PUBLIC_SITE_URL } from "astro:env/client";

/**
 * Convert a relative path to an absolute URL using PUBLIC_SITE_URL.
 * If path is already an absolute http(s) URL, returns it as-is.
 * Safely handles leading/trailing slashes.
 */
export function absoluteUrl(path: string): string {
  // If path is already an absolute URL, return it
  if (/^https?:\/\//.test(path)) {
    return path;
  }

  const baseUrl = PUBLIC_SITE_URL.replace(/\/$/, ""); // Remove trailing slash
  const cleanPath = path.startsWith("/") ? path : `/${path}`;

  return `${baseUrl}${cleanPath}`;
}

/**
 * Get the canonical URL for a given pathname.
 * Canonical URLs should always be stable and production-based, never using Astro.url.href.
 */
export function canonicalFor(pathname: string): string {
  return absoluteUrl(pathname);
}

import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";

dayjs.extend(relativeTime);

/** Long human date, e.g. "Jun 4, 2026" — matches Symfony post.formattedDate. */
export function formatDate(iso: string): string {
  return dayjs(iso).format("MMM D, YYYY");
}

/** Relative "x ago" — matches Symfony comment.diffForHumans. */
export function fromNow(iso: string): string {
  return dayjs(iso).fromNow();
}

/** Strip HTML tags and truncate to `len` chars with ellipsis. */
export function excerpt(html: string | undefined, len = 200): string {
  if (!html) return "";
  const text = html.replace(/<[^>]*>/g, "").trim();
  return text.length > len ? text.slice(0, len) + "..." : text;
}

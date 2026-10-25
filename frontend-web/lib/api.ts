/**
 * Shared API configuration and URL helpers for frontend-web.
 * Configurable via NEXT_PUBLIC_API_URL; defaults to http://localhost:8080.
 */
export const API_BASE_URL: string = (
  process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
).replace(/\/+$/, '');

/**
 * Builds an absolute API URL by trimming trailing slashes and ensuring a leading slash on the path.
 * @param path The API endpoint path, e.g. '/api/login'
 */
export function getApiUrl(path: string): string {
  const cleanPath = path.startsWith('/') ? path : `/${path}`;
  return `${API_BASE_URL}${cleanPath}`;
}

/**
 * Centralized authenticated API client for frontend-web.
 * - Injects customer JWT Bearer token into Authorization header.
 * - Handles JSON and non-JSON responses safely.
 * - Preserves HTTP status codes.
 * - Handles HTTP 401 by clearing invalid token and signaling redirect to /signin.
 * - Preserves specific errors (400, 403, 404, 409, 500).
 * - Supports request cancellation via AbortSignal.
 */

import { getApiUrl } from './api';
import { getAuthToken, clearAuthToken } from './auth';

let isRedirecting = false;

function redirectToSignin(): void {
  clearAuthToken();
  if (typeof window !== 'undefined' && !window.location.pathname.startsWith('/signin') && !isRedirecting) {
    isRedirecting = true;
    window.location.href = '/signin';
  }
}

export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

export interface RequestOptions extends Omit<RequestInit, 'headers'> {
  headers?: Record<string, string>;
  signal?: AbortSignal;
}

export async function authenticatedFetch<T>(
  path: string,
  options: RequestOptions = {}
): Promise<T> {
  if (typeof window === 'undefined') {
    throw new Error('authenticatedFetch cannot be called during server rendering');
  }

  const token = getAuthToken();
  if (!token) {
    redirectToSignin();
    throw new ApiError(401, 'Sesi Anda telah berakhir. Silakan login kembali.');
  }

  const url = getApiUrl(path);
  const headers: Record<string, string> = {
    'Authorization': `Bearer ${token}`,
    ...options.headers,
  };

  if (options.body && typeof options.body === 'string' && !headers['Content-Type']) {
    headers['Content-Type'] = 'application/json';
  }

  let response: Response;
  try {
    response = await fetch(url, {
      ...options,
      headers,
    });
  } catch (err: unknown) {
    if (err instanceof DOMException && err.name === 'AbortError') {
      throw err;
    }
    const message = err instanceof Error ? err.message : 'Koneksi ke server gagal';
    throw new ApiError(0, message);
  }

  let parsedBody: unknown = null;
  const contentType = response.headers.get('content-type') || '';
  if (contentType.includes('application/json')) {
    try {
      parsedBody = await response.json();
    } catch {
      parsedBody = null;
    }
  } else {
    try {
      parsedBody = await response.text();
    } catch {
      parsedBody = null;
    }
  }

  if (!response.ok) {
    if (response.status === 401) {
      redirectToSignin();
    }

    let errorMessage = `HTTP Error ${response.status}`;
    if (parsedBody && typeof parsedBody === 'object') {
      if ('error' in parsedBody) {
        const errVal = (parsedBody as { error: unknown }).error;
        if (typeof errVal === 'string' && errVal.trim() !== '') {
          errorMessage = errVal;
        }
      } else if ('message' in parsedBody) {
        const msgVal = (parsedBody as { message: unknown }).message;
        if (typeof msgVal === 'string' && msgVal.trim() !== '') {
          errorMessage = msgVal;
        }
      }
    } else if (typeof parsedBody === 'string' && parsedBody.trim() !== '') {
      // Bounded text snippet if plaintext error
      const trimmed = parsedBody.trim();
      errorMessage = trimmed.length > 200 ? trimmed.substring(0, 200) : trimmed;
    }

    throw new ApiError(response.status, errorMessage);
  }

  return parsedBody as T;
}

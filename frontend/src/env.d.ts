/// <reference types="astro/client" />

declare const umami:
  | {
      track: (event: string, data?: Record<string, unknown>) => void;
    }
  | undefined;

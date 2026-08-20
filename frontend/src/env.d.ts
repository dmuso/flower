declare global {
  interface ImportMetaEnv {
    readonly FRONTEND_API_URL?: string;
    readonly ENVIRONMENT?: string;
  }
}

export {};

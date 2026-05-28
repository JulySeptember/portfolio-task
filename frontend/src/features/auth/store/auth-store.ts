// src/features/auth/store/auth-store.ts

import { create } from "zustand";

type AuthState = {
  accessToken: string | null;

  idToken: string | null;

  refreshToken: string | null;

  isAuthenticated: boolean;

  setTokens: (tokens: {
    accessToken: string;
    idToken: string;
    refreshToken?: string;
  }) => void;

  clearTokens: () => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  accessToken: null,

  idToken: null,

  refreshToken: null,

  isAuthenticated: false,

  setTokens: ({ accessToken, idToken, refreshToken }) => {
    localStorage.setItem("access_token", accessToken);

    localStorage.setItem("id_token", idToken);

    if (refreshToken) {
      localStorage.setItem("refresh_token", refreshToken);
    }

    set({
      accessToken,

      idToken,

      refreshToken: refreshToken ?? null,

      isAuthenticated: true,
    });
  },

  clearTokens: () => {
    localStorage.removeItem("access_token");

    localStorage.removeItem("id_token");

    localStorage.removeItem("refresh_token");

    set({
      accessToken: null,

      idToken: null,

      refreshToken: null,

      isAuthenticated: false,
    });
  },
}));

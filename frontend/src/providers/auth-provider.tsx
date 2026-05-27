"use client";

import { createContext, useContext } from "react";

import { logout } from "@/features/auth/api/logout";

import { useCurrentUser } from "@/features/auth/hooks/use-current-user";

type AuthContextValue = {
  user: ReturnType<typeof useCurrentUser>["data"];

  isLoading: boolean;

  isAuthenticated: boolean;

  logout: () => void;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { data, isLoading } = useCurrentUser();

  return (
    <AuthContext.Provider
      value={{
        user: data,

        isLoading,

        isAuthenticated: !!data,

        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }

  return context;
}

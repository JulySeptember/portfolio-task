// src/app/(protected)/layout.tsx

import type { ReactNode } from "react";

import { AppHeader } from "@/components/layout/app-header";

import { getCurrentUser } from "@/features/auth/api/get-current-user";

type Props = {
  children: ReactNode;
};

export default async function ProtectedLayout({ children }: Props) {
  const user = await getCurrentUser();

  return (
    <>
      <AppHeader user={user} />

      <main>{children}</main>
    </>
  );
}

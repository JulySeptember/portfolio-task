import type { ReactNode } from "react";

import { cn } from "@/lib/utils/cn";

type Props = {
  children: ReactNode;

  className?: string;
};

export function PageContainer({ children, className }: Props) {
  return (
    <main className={cn("mx-auto w-full max-w-7xl p-6", className)}>
      {children}
    </main>
  );
}

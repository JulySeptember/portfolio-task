import type { ReactNode } from "react";

type Props = {
  children: ReactNode;
};

export function PageContainer({ children }: Props) {
  return <main className="mx-auto w-full max-w-7xl px-6 py-8">{children}</main>;
}

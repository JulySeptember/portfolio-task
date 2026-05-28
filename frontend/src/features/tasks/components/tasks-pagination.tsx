"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";

type Props = {
  total: number;
  limit: number;
  offset: number;
};

export function TasksPagination({ total, limit, offset }: Props) {
  const router = useRouter();
  const searchParams = useSearchParams();

  const hasPrevPage = offset > 0;
  const hasNextPage = offset + limit < total;

  function move(nextOffset: number) {
    const params = new URLSearchParams(searchParams.toString());
    params.set("offset", String(nextOffset));
    router.replace(`/tasks?${params.toString()}`);
  }

  // 0件対応
  const start = total === 0 ? 0 : offset + 1;
  const end = Math.min(offset + limit, total);

  return (
    <div className="flex items-center justify-between">
      <p className="text-muted-foreground text-sm">
        {start}-{end} / {total}
      </p>

      <div className="flex gap-2">
        <Button
          variant="outline"
          disabled={!hasPrevPage}
          onClick={() => move(Math.max(offset - limit, 0))}
        >
          Previous
        </Button>

        <Button
          variant="outline"
          disabled={!hasNextPage}
          onClick={() => move(offset + limit)}
        >
          Next
        </Button>
      </div>
    </div>
  );
}

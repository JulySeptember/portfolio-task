"use client";

import { Button } from "@/components/ui/button";

type Props = {
  offset: number;

  limit: number;

  total: number;

  currentCount: number;

  hasPrevPage: boolean;

  hasNextPage: boolean;

  onPrevious: () => void;

  onNext: () => void;
};

export function TasksPagination({
  offset,
  limit,
  total,
  currentCount,
  hasPrevPage,
  hasNextPage,
  onPrevious,
  onNext,
}: Props) {
  return (
    <div className="flex items-center justify-between">
      <Button variant="outline" disabled={!hasPrevPage} onClick={onPrevious}>
        Previous
      </Button>

      <p className="text-muted-foreground text-sm">
        {offset + 1} - {offset + currentCount} / {total}
      </p>

      <Button variant="outline" disabled={!hasNextPage} onClick={onNext}>
        Next
      </Button>
    </div>
  );
}

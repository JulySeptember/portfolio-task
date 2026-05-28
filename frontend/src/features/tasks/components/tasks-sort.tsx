"use client";

import { useEffect, useState } from "react";

import { useRouter, useSearchParams } from "next/navigation";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

type SortType = "created_at" | "due_date";

type OrderType = "ASC" | "DESC";

export function TasksSort() {
  const router = useRouter();

  const searchParams = useSearchParams();

  const [sort, setSort] = useState<SortType>(
    (searchParams.get("sort") as SortType | null) ?? "created_at",
  );

  const [order, setOrder] = useState<OrderType>(
    (searchParams.get("order") as OrderType | null) ?? "DESC",
  );

  // =========================
  // sync from url
  // =========================

  useEffect(() => {
    setSort((searchParams.get("sort") as SortType | null) ?? "created_at");

    setOrder((searchParams.get("order") as OrderType | null) ?? "DESC");
  }, [searchParams]);

  // =========================
  // update url
  // =========================

  function update(nextSort: SortType, nextOrder: OrderType) {
    const params = new URLSearchParams(searchParams.toString());

    params.set("sort", nextSort);

    params.set("order", nextOrder);

    // sort変えたら先頭ページへ
    params.set("offset", "0");

    router.replace(`/tasks?${params.toString()}`);
  }

  return (
    <div className="flex items-center gap-2">
      {/* sort */}

      <Select
        value={sort}
        onValueChange={(value: SortType) => {
          setSort(value);

          update(value, order);
        }}
      >
        <SelectTrigger className="h-11 min-w-0 flex-1 sm:w-44 sm:flex-none">
          <SelectValue />
        </SelectTrigger>

        <SelectContent>
          <SelectItem value="created_at">Created At</SelectItem>

          <SelectItem value="due_date">Due Date</SelectItem>
        </SelectContent>
      </Select>

      {/* order */}

      <Select
        value={order}
        onValueChange={(value: OrderType) => {
          setOrder(value);

          update(sort, value);
        }}
      >
        <SelectTrigger className="h-11 w-28 shrink-0 sm:w-32">
          <SelectValue />
        </SelectTrigger>

        <SelectContent>
          <SelectItem value="DESC">DESC</SelectItem>

          <SelectItem value="ASC">ASC</SelectItem>
        </SelectContent>
      </Select>
    </div>
  );
}

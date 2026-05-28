"use client";

import { useRouter, useSearchParams } from "next/navigation";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export function TasksSort() {
  const router = useRouter();

  const searchParams = useSearchParams();

  const sort = searchParams.get("sort") ?? "created_at";

  const order = searchParams.get("order") ?? "DESC";

  function update(key: string, value: string) {
    const params = new URLSearchParams(searchParams.toString());

    params.set(key, value);

    router.replace(`/tasks?${params.toString()}`);
  }

  return (
    <div className="flex items-center gap-2">
      <Select value={sort} onValueChange={(value) => update("sort", value)}>
        <SelectTrigger className="h-11 min-w-0 flex-1 sm:w-44 sm:flex-none">
          <SelectValue />
        </SelectTrigger>

        <SelectContent>
          <SelectItem value="created_at">Created At</SelectItem>

          <SelectItem value="due_date">Due Date</SelectItem>
        </SelectContent>
      </Select>

      <Select value={order} onValueChange={(value) => update("order", value)}>
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

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

    router.push(`/tasks?${params.toString()}`);
  }

  return (
    <div className="flex items-center gap-3">
      <Select value={sort} onValueChange={(value) => update("sort", value)}>
        <SelectTrigger className="w-44">
          <SelectValue />
        </SelectTrigger>

        <SelectContent>
          <SelectItem value="created_at">Created At</SelectItem>
          <SelectItem value="due_date">Due Date</SelectItem>
        </SelectContent>
      </Select>

      <Select value={order} onValueChange={(value) => update("order", value)}>
        <SelectTrigger className="w-32">
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

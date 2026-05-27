"use client";

import { useRouter, useSearchParams } from "next/navigation";

import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";

export function TasksFilter() {
  const router = useRouter();

  const searchParams = useSearchParams();

  const status = searchParams.get("status") ?? "ALL";

  function updateStatus(value: string) {
    const params = new URLSearchParams(searchParams.toString());

    if (value === "ALL") {
      params.delete("status");
    } else {
      params.set("status", value);
    }

    params.set("offset", "0");

    router.push(`/tasks?${params.toString()}`);
  }

  return (
    <ToggleGroup
      type="single"
      value={status === "ALL" ? "" : status}
      onValueChange={(value) => {
        updateStatus(value || "ALL");
      }}
      className="flex items-center gap-2"
    >
      <ToggleGroupItem value="TODO" className="h-9 px-3 text-sm">
        TODO
      </ToggleGroupItem>

      <ToggleGroupItem value="DOING" className="h-9 px-3 text-sm">
        DOING
      </ToggleGroupItem>

      <ToggleGroupItem value="DONE" className="h-9 px-3 text-sm">
        DONE
      </ToggleGroupItem>
    </ToggleGroup>
  );
}

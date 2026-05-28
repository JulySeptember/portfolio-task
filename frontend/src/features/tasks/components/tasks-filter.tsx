"use client";

import { useRouter, useSearchParams } from "next/navigation";

import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";

type Status = "ALL" | "TODO" | "DOING" | "DONE";

export function TasksFilter() {
  const router = useRouter();

  const searchParams = useSearchParams();

  const status = (searchParams.get("status") as Status | null) ?? "ALL";

  function updateStatus(value: Status) {
    const params = new URLSearchParams(searchParams.toString());

    if (value === "ALL") {
      params.delete("status");
    } else {
      params.set("status", value);
    }

    // filter変更時は1ページ目へ戻す
    params.set("offset", "0");

    router.replace(`/tasks?${params.toString()}`);
  }

  return (
    <ToggleGroup
      type="single"
      value={status}
      onValueChange={(value) => {
        if (!value) {
          updateStatus("ALL");
          return;
        }

        updateStatus(value as Status);
      }}
      className="flex flex-wrap items-center gap-2"
    >
      <ToggleGroupItem value="ALL" className="h-9 px-3 text-sm">
        ALL
      </ToggleGroupItem>

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

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
      value={status}
      onValueChange={(value) => {
        if (value) updateStatus(value);
      }}
    >
      <ToggleGroupItem value="ALL">ALL</ToggleGroupItem>

      <ToggleGroupItem value="TODO">TODO</ToggleGroupItem>

      <ToggleGroupItem value="DOING">DOING</ToggleGroupItem>

      <ToggleGroupItem value="DONE">DONE</ToggleGroupItem>
    </ToggleGroup>
  );
}

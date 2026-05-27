"use client";

import { useRouter, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";

export function TasksToolbar() {
  const router = useRouter();

  const searchParams = useSearchParams();

  function openCreate() {
    const params = new URLSearchParams(searchParams.toString());

    params.set("create", "true");

    router.push(`/tasks?${params.toString()}`);
  }

  return (
    <div className="flex justify-end">
      <Button onClick={openCreate}>Create Task</Button>
    </div>
  );
}

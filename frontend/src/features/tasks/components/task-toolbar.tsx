"use client";

import { useRouter, useSearchParams } from "next/navigation";

import { Plus } from "lucide-react";

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
      <Button onClick={openCreate} className="h-11 rounded-xl px-5">
        <Plus className="mr-2 h-4 w-4" />
        Create Task
      </Button>
    </div>
  );
}

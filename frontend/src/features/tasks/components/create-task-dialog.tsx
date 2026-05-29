// src/features/tasks/components/create-task-dialog.tsx

"use client";

import { useRouter, useSearchParams } from "next/navigation";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { TaskEditor } from "./task-editor";

export function CreateTaskDialog() {
  const router = useRouter();

  const searchParams = useSearchParams();

  const open = searchParams.get("create") === "true";

  function closeDialog() {
    const params = new URLSearchParams(searchParams.toString());

    params.delete("create");

    const query = params.toString();

    router.replace(query ? `/tasks?${query}` : "/tasks", {
      scroll: false,
    });
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(nextOpen) => {
        if (!nextOpen) {
          closeDialog();
        }
      }}
    >
      <DialogContent
        onPointerDownOutside={(e) => e.preventDefault()}
        onEscapeKeyDown={(e) => e.preventDefault()}
        className="
          w-[95vw]
          max-w-4xl
          overflow-y-auto
          rounded-2xl
          p-6
          sm:p-8
        "
      >
        <DialogHeader className="sr-only">
          <DialogTitle>Create Task</DialogTitle>

          <DialogDescription>
            Create a new task with title, description, and due date.
          </DialogDescription>
        </DialogHeader>

        <TaskEditor mode="create" onSuccess={closeDialog} />
      </DialogContent>
    </Dialog>
  );
}

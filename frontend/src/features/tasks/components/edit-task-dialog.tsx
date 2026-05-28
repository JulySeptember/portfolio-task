// src/features/tasks/components/edit-task-dialog.tsx

"use client";

import { useRouter, useSearchParams } from "next/navigation";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { useTask } from "../hooks/use-task";

import { TaskEditor } from "./task-editor";

type Props = {
  taskId: number;
};

export function EditTaskDialog({ taskId }: Props) {
  const router = useRouter();

  const searchParams = useSearchParams();

  const { data: task, isLoading } = useTask(taskId);

  function closeDialog() {
    const params = new URLSearchParams(searchParams.toString());

    params.delete("taskId");

    const query = params.toString();

    router.replace(query ? `/tasks?${query}` : "/tasks", {
      scroll: false,
    });
  }

  if (isLoading || !task) {
    return null;
  }

  return (
    <Dialog
      open
      onOpenChange={(open) => {
        if (!open) {
          closeDialog();
        }
      }}
    >
      <DialogContent
        className="
          w-[95vw]
          max-w-6xl
          overflow-y-auto
          rounded-2xl
          p-6
          sm:p-8
        "
        onPointerDownOutside={(e) => e.preventDefault()}
        onEscapeKeyDown={(e) => e.preventDefault()}
      >
        <DialogHeader className="sr-only">
          <DialogTitle>Edit Task</DialogTitle>

          <DialogDescription>
            Edit task details and update task status.
          </DialogDescription>
        </DialogHeader>

        <TaskEditor
          mode="edit"
          task={task}
          onSuccess={closeDialog}
          showOpenPageButton
        />
      </DialogContent>
    </Dialog>
  );
}

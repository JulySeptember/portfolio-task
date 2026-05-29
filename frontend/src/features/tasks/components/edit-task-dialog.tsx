"use client";

import { useRouter, useSearchParams } from "next/navigation";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import type { Task } from "../schemas/task-schema";

import { TaskEditor } from "./task-editor";

type Props = {
  task: Task;
  onClose: () => void;
};

export function EditTaskDialog({ task, onClose }: Props) {
  const router = useRouter();

  const searchParams = useSearchParams();

  const open = searchParams.has("taskId");
  function closeDialog() {
    const params = new URLSearchParams(searchParams.toString());

    params.delete("taskId");

    const query = params.toString();

    router.replace(query ? `/tasks?${query}` : "/tasks", {
      scroll: false,
    });

    onClose();
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
          <DialogTitle>Edit Task</DialogTitle>
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

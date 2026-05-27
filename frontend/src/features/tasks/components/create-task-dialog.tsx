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

import { useCreateTask } from "../hooks/use-create-task";

import { TaskForm } from "./task-form";

export function CreateTaskDialog() {
  const router = useRouter();

  const searchParams = useSearchParams();

  const open = searchParams.get("create") === "true";

  const createTask = useCreateTask();

  function closeDialog() {
    const params = new URLSearchParams(searchParams.toString());

    params.delete("create");

    router.replace(`/tasks?${params.toString()}`, {
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
      <DialogContent className="w-[95vw] max-w-6xl! overflow-y-auto p-10">
        <div className="mx-auto w-full max-w-4xl">
          <DialogHeader className="space-y-3">
            <DialogTitle className="text-4xl font-bold">
              Create Task
            </DialogTitle>

            <DialogDescription className="text-base">
              Create a new task.
            </DialogDescription>
          </DialogHeader>

          <TaskForm
            submitLabel="Create Task"
            isPending={createTask.isPending}
            defaultValues={{
              title: "",
              description: "",
              status: "TODO",
              due_date: "",
            }}
            onSubmit={(values) => {
              createTask.mutate(
                {
                  ...values,

                  due_date:
                    values.due_date && values.due_date !== ""
                      ? new Date(values.due_date).toISOString()
                      : null,
                },
                {
                  onSuccess: () => {
                    closeDialog();
                  },
                },
              );
            }}
          />
        </div>
      </DialogContent>
    </Dialog>
  );
}

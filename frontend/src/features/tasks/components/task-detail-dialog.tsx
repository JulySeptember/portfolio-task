// src/features/tasks/components/task-detail-dialog.tsx

"use client";

import { useRouter, useSearchParams } from "next/navigation";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { Button } from "@/components/ui/button";

import { useTask } from "../hooks/use-task";

import { useDeleteTask } from "../hooks/use-delete-task";

import { useUpdateTask } from "../hooks/use-update-task";

import { TaskForm } from "./task-form";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";

type Props = {
  taskId: number;
};

export function TaskDetailDialog({ taskId }: Props) {
  const router = useRouter();

  const searchParams = useSearchParams();

  const { data: task, isLoading } = useTask(taskId);

  const updateTask = useUpdateTask();

  const deleteTask = useDeleteTask();

  function closeDialog() {
    const params = new URLSearchParams(searchParams.toString());

    params.delete("taskId");

    router.replace(`/tasks?${params.toString()}`, {
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
      <DialogContent className="w-[95vw] max-w-6xl! overflow-y-auto p-10">
        <div className="mx-auto w-full max-w-4xl space-y-8">
          <DialogHeader className="space-y-3">
            <DialogTitle className="text-4xl font-bold">Edit Task</DialogTitle>

            <DialogDescription className="text-base">
              Update task information.
            </DialogDescription>
          </DialogHeader>
          <TaskForm
            submitLabel="Save Changes"
            isPending={updateTask.isPending}
            defaultValues={{
              title: task.title,
              description: task.description,
              status: task.status,
              due_date: task.dueDate ? task.dueDate.slice(0, 16) : "",
            }}
            onSubmit={(values) => {
              updateTask.mutate(
                {
                  id: task.id,

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
          <div className="flex justify-end border-t pt-6">
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button variant="destructive">Delete Task</Button>
              </AlertDialogTrigger>

              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Delete task?</AlertDialogTitle>

                  <AlertDialogDescription>
                    This action cannot be undone.
                  </AlertDialogDescription>
                </AlertDialogHeader>

                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>

                  <AlertDialogAction
                    onClick={() => {
                      deleteTask.mutate(task.id, {
                        onSuccess: () => {
                          closeDialog();
                        },
                      });
                    }}
                  >
                    {deleteTask.isPending ? "Deleting..." : "Delete"}
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>{" "}
        </div>
      </DialogContent>
    </Dialog>
  );
}

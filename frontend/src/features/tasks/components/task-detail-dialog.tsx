"use client";

import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { Button } from "@/components/ui/button";

import { Input } from "@/components/ui/input";

import { Label } from "@/components/ui/label";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { Textarea } from "@/components/ui/textarea";

import { TaskStatusBadge } from "./task-status-badge";

import {
  taskFormSchema,
  type Task,
  type TaskFormValues,
  type TaskStatus,
} from "../schemas/task-schema";

import { useUpdateTask } from "../hooks/use-update-task";

import { useDeleteTask } from "../hooks/use-delete-task";

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
  task: Task;

  open: boolean;

  onOpenChange: (open: boolean) => void;
};

export function TaskDetailDialog({ task, open, onOpenChange }: Props) {
  const updateTask = useUpdateTask();

  const deleteTask = useDeleteTask();

  const form = useForm<TaskFormValues>({
    resolver: zodResolver(taskFormSchema),

    values: {
      title: task.title,

      description: task.description,

      status: task.status,

      due_date: task.dueDate ? task.dueDate.slice(0, 16) : "",
    },
  });

  const onSubmit = (values: TaskFormValues) => {
    updateTask.mutate(
      {
        id: task.id,

        title: values.title,

        description: values.description,

        status: values.status,

        due_date:
          values.due_date && values.due_date !== ""
            ? new Date(values.due_date).toISOString()
            : null,
      },
      {
        onSuccess: () => {
          onOpenChange(false);
        },
      },
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        className="max-h-[90vh] w-[95vw] max-w-6xl! overflow-y-auto p-10"
        onPointerDownOutside={(e) => e.preventDefault()}
        onEscapeKeyDown={(e) => e.preventDefault()}
      >
        <div className="mx-auto w-full max-w-4xl">
          <DialogHeader className="space-y-4">
            <div className="flex items-center gap-4">
              <DialogTitle className="text-4xl font-bold">
                {task.title}
              </DialogTitle>

              <TaskStatusBadge status={task.status} />
            </div>

            <div className="space-y-1 text-sm text-muted-foreground">
              <p>Created: {new Date(task.createdAt).toLocaleString()}</p>

              <p>Updated: {new Date(task.updatedAt).toLocaleString()}</p>
            </div>
          </DialogHeader>

          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="mt-10 space-y-8"
          >
            <div className="space-y-3">
              <Label className="text-base font-medium">Title</Label>

              <Input
                className="h-12 text-base"
                placeholder="Enter task title"
                {...form.register("title")}
              />

              {form.formState.errors.title && (
                <p className="text-sm text-red-500">
                  {form.formState.errors.title.message}
                </p>
              )}
            </div>

            <div className="space-y-3">
              <Label className="text-base font-medium">Description</Label>

              <Textarea
                className="min-h-52 text-base"
                placeholder="Enter task description"
                {...form.register("description")}
              />

              {form.formState.errors.description && (
                <p className="text-sm text-red-500">
                  {form.formState.errors.description.message}
                </p>
              )}
            </div>

            <div className="grid gap-8 md:grid-cols-2">
              <div className="space-y-3">
                <Label className="text-base font-medium">Status</Label>

                <Select
                  value={form.watch("status")}
                  onValueChange={(value: TaskStatus) =>
                    form.setValue("status", value)
                  }
                >
                  <SelectTrigger className="h-12 w-full text-base">
                    <SelectValue />
                  </SelectTrigger>

                  <SelectContent>
                    <SelectItem value="TODO">TODO</SelectItem>

                    <SelectItem value="DOING">DOING</SelectItem>

                    <SelectItem value="DONE">DONE</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-3">
                <Label className="text-base font-medium">Due Date</Label>

                <Input
                  type="datetime-local"
                  className="h-12 text-base"
                  {...form.register("due_date")}
                />
              </div>
            </div>

            <div className="sticky bottom-0 flex items-center justify-between border-t bg-background pt-6">
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button
                    type="button"
                    variant="destructive"
                    className="h-12 px-8 text-base"
                    disabled={deleteTask.isPending}
                  >
                    {deleteTask.isPending ? "削除中..." : "Delete"}
                  </Button>
                </AlertDialogTrigger>

                <AlertDialogContent className="max-w-lg p-8">
                  <AlertDialogHeader className="space-y-4">
                    <AlertDialogTitle className="text-2xl font-bold">
                      タスクを削除しますか？
                    </AlertDialogTitle>

                    <AlertDialogDescription className="text-base leading-relaxed">
                      「{task.title}」を削除します。
                      <br />
                      <br />
                      この操作は取り消せません。
                    </AlertDialogDescription>
                  </AlertDialogHeader>

                  <AlertDialogFooter className="mt-6 gap-3">
                    <AlertDialogCancel className="h-11 text-base">
                      Cancel
                    </AlertDialogCancel>

                    <AlertDialogAction
                      className="h-11 text-base"
                      onClick={() => {
                        deleteTask.mutate(task.id, {
                          onSuccess: () => {
                            onOpenChange(false);
                          },
                        });
                      }}
                    >
                      Delete
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>

              <Button
                type="submit"
                className="h-12 px-8 text-base"
                disabled={updateTask.isPending}
              >
                {updateTask.isPending ? "Updating..." : "Save Changes"}
              </Button>
            </div>
          </form>
        </div>
      </DialogContent>
    </Dialog>
  );
}

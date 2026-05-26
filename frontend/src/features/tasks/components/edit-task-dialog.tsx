"use client";

import { useEffect, useState } from "react";

import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { Button } from "@/components/ui/button";

import { Input } from "@/components/ui/input";

import { Label } from "@/components/ui/label";

import { Textarea } from "@/components/ui/textarea";

import {
  taskRequestSchema,
  type Task,
  type TaskRequest,
} from "../schemas/task-schema";

import { useUpdateTask } from "../hooks/use-update-task";

type Props = {
  task: Task;
};

export function EditTaskDialog({ task }: Props) {
  const [open, setOpen] = useState(false);

  const updateTask = useUpdateTask();

  const form = useForm<TaskRequest>({
    resolver: zodResolver(taskRequestSchema),

    defaultValues: {
      title: task.title,

      description: task.description,

      status: task.status,

      due_date: task.dueDate,
    },
  });

  useEffect(() => {
    form.reset({
      title: task.title,

      description: task.description,

      status: task.status,

      due_date: task.dueDate,
    });
  }, [form, task]);

  const onSubmit = (values: TaskRequest) => {
    updateTask.mutate(
      {
        id: task.id,

        data: values,
      },

      {
        onSuccess: () => {
          setOpen(false);
        },
      },
    );
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm" variant="outline">
          Edit
        </Button>
      </DialogTrigger>

      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Task</DialogTitle>
        </DialogHeader>

        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">Title</Label>

            <Input id="title" {...form.register("title")} />

            {form.formState.errors.title && (
              <p className="text-sm text-red-500">
                {form.formState.errors.title.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>

            <Textarea id="description" {...form.register("description")} />

            {form.formState.errors.description && (
              <p className="text-sm text-red-500">
                {form.formState.errors.description.message}
              </p>
            )}
          </div>

          <div className="flex justify-end gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => setOpen(false)}
            >
              Cancel
            </Button>

            <Button type="submit" disabled={updateTask.isPending}>
              {updateTask.isPending ? "Saving..." : "Save"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

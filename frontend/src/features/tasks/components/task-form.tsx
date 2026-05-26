"use client";

import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";

import { taskRequestSchema, type TaskRequest } from "../schemas/task-schema";

import { useCreateTask } from "../hooks/use-create-task";

import { Button } from "@/components/ui/button";

import { Input } from "@/components/ui/input";

import { Label } from "@/components/ui/label";

import { Textarea } from "@/components/ui/textarea";

export function TaskForm() {
  const createTask = useCreateTask();

  const form = useForm<TaskRequest>({
    resolver: zodResolver(taskRequestSchema),

    defaultValues: {
      title: "",
      description: "",
      status: "TODO",
      due_date: null,
    },
  });

  const onSubmit = (values: TaskRequest) => {
    createTask.mutate(values);

    form.reset();
  };

  return (
    <form
      onSubmit={form.handleSubmit(onSubmit)}
      className="space-y-4 rounded-lg border p-4"
    >
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

      <Button type="submit" disabled={createTask.isPending}>
        {createTask.isPending ? "Creating..." : "Create Task"}
      </Button>
    </form>
  );
}

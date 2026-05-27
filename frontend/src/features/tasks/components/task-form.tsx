// src/features/tasks/components/task-form.tsx

"use client";

import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";

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

import {
  taskFormSchema,
  type TaskFormValues,
  type TaskStatus,
} from "../schemas/task-schema";

type Props = {
  defaultValues: TaskFormValues;

  submitLabel: string;

  isPending?: boolean;

  onSubmit: (values: TaskFormValues) => void;
};

export function TaskForm({
  defaultValues,
  submitLabel,
  isPending,
  onSubmit,
}: Props) {
  const form = useForm<TaskFormValues>({
    resolver: zodResolver(taskFormSchema),

    defaultValues,
  });

  return (
    <form onSubmit={form.handleSubmit(onSubmit)} className="mt-10 space-y-8">
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
          className="min-h-40 text-base"
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

      <div className="flex justify-end pt-4">
        <Button
          type="submit"
          className="h-12 px-8 text-base"
          disabled={isPending}
        >
          {isPending ? "Saving..." : submitLabel}
        </Button>
      </div>
    </form>
  );
}

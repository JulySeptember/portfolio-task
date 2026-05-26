"use client";

import { useState } from "react";

import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
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

import {
  taskFormSchema,
  type TaskFormValues,
  type TaskStatus,
} from "../schemas/task-schema";

import { useCreateTask } from "../hooks/use-create-task";

export function CreateTaskDialog() {
  const [open, setOpen] = useState(false);

  const createTask = useCreateTask();

  const form = useForm<TaskFormValues>({
    resolver: zodResolver(taskFormSchema),

    defaultValues: {
      title: "",

      description: "",

      status: "TODO",

      due_date: "",
    },
  });

  const onSubmit = (values: TaskFormValues) => {
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
          form.reset();

          setOpen(false);
        },
      },
    );
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="h-12 px-6 text-base">Create Task</Button>
      </DialogTrigger>

      <DialogContent className="w-[95vw] max-w-6xl! overflow-y-auto p-10">
        <div className="mx-auto w-full max-w-4xl">
          <DialogHeader className="space-y-3">
            <DialogTitle className="text-4xl font-bold">
              Create Task
            </DialogTitle>

            <DialogDescription className="text-base">
              Create a new task with status and due date.
            </DialogDescription>
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
                  defaultValue={form.getValues("status")}
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
                disabled={createTask.isPending}
              >
                {createTask.isPending ? "Creating..." : "Create Task"}
              </Button>
            </div>
          </form>
        </div>
      </DialogContent>
    </Dialog>
  );
}

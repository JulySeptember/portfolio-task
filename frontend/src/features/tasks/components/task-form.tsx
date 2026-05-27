// src/features/tasks/components/task-form.tsx

"use client";

import type { ReactNode } from "react";
import { useEffect, useRef, useState } from "react";
import { format } from "date-fns";
import { CalendarIcon } from "lucide-react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { cn } from "@/lib/utils";
import { TaskStatusSelect } from "./task-status-select";
import {
  taskFormSchema,
  type TaskFormInput,
  type TaskFormValues,
  type TaskStatus,
} from "../schemas/task-schema";

type Props = {
  defaultValues: TaskFormInput;
  submitLabel: string;
  isPending?: boolean;
  onSubmit: (values: TaskFormValues) => void;
  footer?: ReactNode;
  secondaryAction?: ReactNode;
  autoResizeDescription?: boolean;
};

export function TaskForm({
  defaultValues,
  submitLabel,
  isPending,
  onSubmit,
  footer,
  secondaryAction,
  autoResizeDescription = false,
}: Props) {
  const [open, setOpen] = useState(false);
  const form = useForm<TaskFormInput>({
    resolver: zodResolver(taskFormSchema),
    defaultValues,
  });

  const dueDate = form.watch("due_date");
  const description = form.watch("description");
  const textareaRef = useRef<HTMLTextAreaElement | null>(null);

  // Auto-resize description for full-page forms
  useEffect(() => {
    if (!autoResizeDescription) return;
    const textarea = textareaRef.current;
    if (!textarea) return;
    textarea.style.height = "auto";
    textarea.style.height = `${textarea.scrollHeight}px`;
  }, [description, autoResizeDescription]);

  return (
    <form
      onSubmit={form.handleSubmit((values) => {
        if (isPending) return;
        onSubmit(values);
      })}
      className="mt-4 w-full min-w-0 space-y-6"
    >
      {/* Title */}
      <div className="min-w-0 space-y-2">
        <Label className="text-sm font-medium tracking-tight">Title</Label>
        <Input
          className="h-12 w-full min-w-0 rounded-xl px-4 text-base md:text-lg"
          placeholder="Enter task title"
          {...form.register("title")}
        />
        {form.formState.errors.title && (
          <p className="text-sm text-red-500">
            {form.formState.errors.title.message}
          </p>
        )}
      </div>

      {/* Status / Due Date */}
      <div className="grid min-w-0 gap-4 md:grid-cols-2">
        <div className="min-w-0 space-y-2">
          <Label className="text-sm font-medium tracking-tight">Status</Label>
          <TaskStatusSelect
            value={form.watch("status")}
            onChange={(value: TaskStatus) => form.setValue("status", value)}
            className="h-12! w-full rounded-xl text-sm md:text-base"
          />
        </div>
        <div className="min-w-0 space-y-2">
          <Label className="text-sm font-medium tracking-tight">Due Date</Label>
          <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
              <Button
                type="button"
                variant="outline"
                className={cn(
                  "h-12 w-full justify-start overflow-hidden rounded-xl px-4 text-left text-sm font-normal md:text-base",
                  !dueDate && "text-muted-foreground",
                )}
              >
                <CalendarIcon className="mr-3 h-4 w-4 shrink-0" />
                <span className="truncate">
                  {dueDate
                    ? format(new Date(dueDate), "yyyy/MM/dd")
                    : "Select due date"}
                </span>
              </Button>
            </PopoverTrigger>
            <PopoverContent
              className="w-auto max-w-[calc(100vw-2rem)] rounded-xl p-0"
              align="start"
            >
              <Calendar
                mode="single"
                selected={dueDate ? new Date(dueDate) : undefined}
                onSelect={(date: Date | undefined) => {
                  if (!date) return;
                  form.setValue("due_date", date.toISOString());
                  setOpen(false);
                }}
              />
            </PopoverContent>
          </Popover>
        </div>
      </div>

      {/* Description */}
      <div className="min-w-0 space-y-2">
        <Textarea
          placeholder="Add details about this task..."
          {...form.register("description")}
          ref={(el) => {
            form.register("description").ref(el);
            textareaRef.current = el;
          }}
          className={cn(
            "w-full min-w-0 rounded-2xl px-5 py-4 text-sm leading-7 md:text-base",
            autoResizeDescription
              ? "min-h-80 resize-none overflow-hidden md:min-h-125"
              : "min-h-56 resize-y md:min-h-120",
          )}
        />
        {form.formState.errors.description && (
          <p className="text-sm text-red-500">
            {form.formState.errors.description.message}
          </p>
        )}
      </div>

      {/* Footer */}
      <div className="flex flex-col gap-4 border-t pt-5 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-2">{secondaryAction}</div>
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
          <div className="w-full sm:w-auto">{footer}</div>
          <Button
            type="submit"
            className="h-12 w-full sm:w-42.5 shrink-0 rounded-xl px-8 text-base font-medium transition-none focus-visible:ring-0 focus-visible:outline-none"
          >
            {isPending ? "Saving..." : submitLabel}
          </Button>
        </div>
      </div>
    </form>
  );
}

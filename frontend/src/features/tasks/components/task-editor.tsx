// src/features/tasks/components/task-editor.tsx
"use client";

import { useState } from "react";

import { format } from "date-fns";

import { CalendarIcon } from "lucide-react";

import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";

import { useRouter, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";

import { Input } from "@/components/ui/input";

import { Label } from "@/components/ui/label";

import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

import { Calendar } from "@/components/ui/calendar";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";

import { cn } from "@/lib/utils";

import { TaskStatusSelect } from "./task-status-select";

import { TaskDescriptionField } from "./task-description-field";

import { useCreateTask } from "../hooks/use-create-task";

import { useDeleteTask } from "../hooks/use-delete-task";

import { useUpdateTask } from "../hooks/use-update-task";

import {
  taskFormSchema,
  type Task,
  type TaskFormValues,
  type TaskStatus,
} from "../schemas/task-schema";

// =========================
// props
// =========================

type Props =
  | {
      mode: "create";

      onSuccess?: () => void;
    }
  | {
      mode: "edit";

      task: Task;

      onSuccess?: () => void;
    };

// =========================
// default values helper
// =========================

function getDefaultValues(props: Props): TaskFormValues {
  // create mode
  if (props.mode === "create") {
    return {
      title: "",

      description: "",

      status: "TODO",

      due_date: "",
    };
  }

  // edit mode
  return {
    title: props.task.title,

    description: props.task.description ?? "",

    status: props.task.status as TaskStatus,

    due_date: props.task.dueDate ?? "",
  };
}

// =========================
// component
// =========================

export function TaskEditor(props: Props) {
  const router = useRouter();

  const searchParams = useSearchParams();

  // mutations

  const createTask = useCreateTask();

  const updateTask = useUpdateTask();

  const deleteTask = useDeleteTask();

  // calendar popover state

  const [open, setOpen] = useState(false);

  // form default values

  const defaultValues = getDefaultValues(props);

  // react-hook-form

  const form = useForm<TaskFormValues>({
    resolver: zodResolver(taskFormSchema),

    defaultValues,
  });

  // watched fields

  const due_date = form.watch("due_date");

  // =========================
  // submit
  // =========================

  async function handleSubmit(values: TaskFormValues) {
    const payload = {
      title: values.title,

      description: values.description ?? "",

      status: values.status,

      due_date:
        values.due_date && values.due_date !== ""
          ? new Date(values.due_date).toISOString()
          : null,
    };

    // =========================
    // create
    // =========================

    if (props.mode === "create") {
      createTask.mutate(payload, {
        onSuccess: () => {
          props.onSuccess?.();
        },
      });

      return;
    }

    // =========================
    // update
    // =========================

    updateTask.mutate(
      {
        publicId: props.task.publicId,

        ...payload,
      },
      {
        onSuccess: () => {
          props.onSuccess?.();
        },
      },
    );
  }

  return (
    <form
      onSubmit={form.handleSubmit(async (values) => {
        // prevent double submit

        if (createTask.isPending || updateTask.isPending) {
          return;
        }

        await handleSubmit(values);
      })}
      className="
        mt-4
        w-full
        min-w-0
        space-y-6
      "
    >
      {/* =========================
          title
      ========================= */}

      <div className="min-w-0 space-y-2">
        <Label className="text-sm font-medium tracking-tight">Title</Label>

        <Input
          className="
            h-12
            w-full
            rounded-xl
            px-4
            text-base
            md:text-lg
          "
          placeholder="Enter task title"
          onKeyDown={(e) => {
            // prevent accidental submit

            if (e.key === "Enter") {
              e.preventDefault();
            }
          }}
          {...form.register("title")}
        />

        {/* validation error */}

        {form.formState.errors.title && (
          <p className="text-sm text-red-400">
            {form.formState.errors.title.message}
          </p>
        )}
      </div>

      {/* =========================
          status / due date
      ========================= */}

      <div
        className="
          grid
          gap-4
          md:grid-cols-2
        "
      >
        {/* status */}

        <div className="space-y-2">
          <Label className="text-sm font-medium tracking-tight">Status</Label>

          <TaskStatusSelect
            value={form.watch("status")}
            onChange={(value: TaskStatus) => {
              form.setValue("status", value);
            }}
            className="
              h-12!
              rounded-xl
              text-sm
              md:text-base
            "
          />
        </div>

        {/* due date */}

        <div className="space-y-2">
          <Label className="text-sm font-medium tracking-tight">Due Date</Label>

          <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
              <Button
                type="button"
                variant="outline"
                className={cn(
                  `
                    h-12
                    w-full
                    justify-start
                    overflow-hidden
                    rounded-xl
                    px-4
                    text-left
                    text-sm
                    font-normal
                    md:text-base
                  `,
                  !due_date && "text-muted-foreground",
                )}
              >
                <CalendarIcon
                  className="
                    mr-3
                    h-4
                    w-4
                    shrink-0
                  "
                />

                <span className="truncate">
                  {due_date
                    ? format(new Date(due_date), "yyyy/MM/dd")
                    : "Select due date"}
                </span>
              </Button>
            </PopoverTrigger>

            <PopoverContent
              className="
                w-auto
                rounded-xl
                p-0
              "
              align="start"
            >
              <Calendar
                mode="single"
                selected={due_date ? new Date(due_date) : undefined}
                onSelect={(date) => {
                  if (!date) {
                    return;
                  }

                  form.setValue("due_date", date.toISOString());

                  setOpen(false);
                }}
              />
            </PopoverContent>
          </Popover>
        </div>
      </div>

      {/* =========================
          description
      ========================= */}

      <TaskDescriptionField form={form} />

      {/* =========================
          footer
      ========================= */}

      <div
        className="
          flex
          flex-col
          gap-3
          border-t
          pt-5
          sm:flex-row
          sm:items-center
          sm:justify-between
        "
      >
        {/* left side */}

        <div />

        {/* right side */}

        <div
          className="
            flex
            flex-col
            gap-3
            sm:flex-row
            sm:items-center
          "
        >
          {/* =========================
              delete button
          ========================= */}

          {props.mode === "edit" && (
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button
                  variant="destructive"
                  className="
                    h-12
                    w-full
                    rounded-xl
                    px-6
                    text-base
                    font-medium
                    sm:w-44
                  "
                >
                  Delete Task
                </Button>
              </AlertDialogTrigger>

              <AlertDialogContent className="rounded-2xl">
                <AlertDialogHeader className="space-y-3">
                  <AlertDialogTitle className="text-xl">
                    Delete this task?
                  </AlertDialogTitle>

                  <p
                    className="
                      text-sm
                      leading-6
                      text-muted-foreground
                    "
                  >
                    This action cannot be undone.
                  </p>
                </AlertDialogHeader>

                <AlertDialogFooter
                  className="
                    flex-col-reverse
                    gap-3
                    sm:flex-row
                    sm:justify-end
                  "
                >
                  <AlertDialogCancel
                    className="
                      h-12
                      rounded-xl
                      px-6
                      text-base
                    "
                  >
                    Cancel
                  </AlertDialogCancel>

                  <AlertDialogAction
                    disabled={deleteTask.isPending}
                    className="
                      h-12
                      rounded-xl
                      px-6
                      text-base
                      font-medium
                    "
                    onClick={() => {
                      if (props.mode !== "edit") {
                        return;
                      }

                      deleteTask.mutate(props.task.publicId, {
                        onSuccess: () => {
                          // keep current filters/sort/pagination

                          const params = new URLSearchParams(
                            searchParams.toString(),
                          );

                          // close dialog

                          params.delete("taskId");

                          router.replace(
                            params.toString()
                              ? `/tasks?${params.toString()}`
                              : "/tasks",
                          );
                        },
                      });
                    }}
                  >
                    Delete Task
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          )}

          {/* =========================
              save button
          ========================= */}

          <Button
            type="submit"
            disabled={createTask.isPending || updateTask.isPending}
            className="
              h-12
              w-full
              rounded-xl
              px-8
              text-base
              font-medium
              sm:w-42.5
            "
          >
            SAVE
          </Button>
        </div>
      </div>
    </form>
  );
}

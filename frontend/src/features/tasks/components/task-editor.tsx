// src/features/tasks/components/task-editor.tsx

"use client";

import Link from "next/link";

import { Button } from "@/components/ui/button";

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

import { useCreateTask } from "../hooks/use-create-task";

import { useDeleteTask } from "../hooks/use-delete-task";

import { useUpdateTask } from "../hooks/use-update-task";

import { useRouter } from "next/navigation";

import type { Task, TaskFormValues, TaskStatus } from "../schemas/task-schema";

import { TaskForm } from "./task-form";

type Props =
  | {
      mode: "create";

      onSuccess?: () => void;

      showOpenPageButton?: boolean;

      autoResizeDescription?: boolean;
    }
  | {
      mode: "edit";

      task: Task;

      onSuccess?: () => void;

      showOpenPageButton?: boolean;

      autoResizeDescription?: boolean;
    };

export function TaskEditor(props: Props) {
  const createTask = useCreateTask();

  const updateTask = useUpdateTask();

  const deleteTask = useDeleteTask();

  const router = useRouter();

  async function handleSubmit(values: TaskFormValues) {
    const payload = {
      ...values,

      due_date:
        values.due_date && values.due_date !== ""
          ? new Date(values.due_date).toISOString()
          : null,
    };

    if (props.mode === "create") {
      createTask.mutate(payload, {
        onSuccess: () => {
          props.onSuccess?.();
        },
      });

      return;
    }

    updateTask.mutate(
      {
        id: props.task.id,
        ...payload,
      },
      {
        onSuccess: () => {
          props.onSuccess?.();
        },
      },
    );
  }

  if (props.mode === "create") {
    return (
      <TaskForm
        submitLabel="Create Task"
        isPending={createTask.isPending}
        defaultValues={{
          title: "",
          description: "",
          status: "TODO",
          due_date: "",
        }}
        onSubmit={handleSubmit}
        autoResizeDescription={props.autoResizeDescription}
        secondaryAction={
          props.showOpenPageButton === false ? null : (
            <Button asChild variant="outline" className="rounded-xl">
              <Link href="/tasks/new">Open Full Page</Link>
            </Button>
          )
        }
      />
    );
  }

  return (
    <TaskForm
      submitLabel="Save Changes"
      isPending={updateTask.isPending}
      autoResizeDescription={props.autoResizeDescription}
      defaultValues={{
        title: props.task.title,
        description: props.task.description,
        status: props.task.status as TaskStatus,
        due_date: props.task.dueDate ? props.task.dueDate.slice(0, 16) : "",
      }}
      onSubmit={handleSubmit}
      secondaryAction={
        props.showOpenPageButton ? (
          <Button asChild variant="outline" className="rounded-xl">
            <Link href={`/tasks/${props.task.id}`}>Open Full Page</Link>
          </Button>
        ) : null
      }
      footer={
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button variant="destructive" className="h-12 rounded-xl px-6">
              Delete Task
            </Button>
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
                disabled={deleteTask.isPending}
                onClick={() => {
                  deleteTask.mutate(props.task.id, {
                    onSuccess: () => {
                      props.onSuccess?.();

                      if (!props.onSuccess) {
                        router.push("/tasks");
                      }
                    },
                  });
                }}
              >
                {deleteTask.isPending ? "Deleting..." : "Delete"}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      }
    />
  );
}

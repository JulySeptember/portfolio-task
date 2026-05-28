"use client";

import { useRouter } from "next/navigation";

import Link from "next/link";

import { Expand } from "lucide-react";

import { Button } from "@/components/ui/button";

import { encodeId } from "@/lib/utils/hash-id";

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

import type { Task, TaskFormValues, TaskStatus } from "../schemas/task-schema";

import { TaskForm } from "./task-form";

type Props =
  | {
      mode: "create";

      onSuccess?: () => void;

      showOpenPageButton?: boolean;

      autoResizeDescription?: boolean;

      onOpenFullPage?: () => void;
    }
  | {
      mode: "edit";

      task: Task;

      onSuccess?: () => void;

      showOpenPageButton?: boolean;

      autoResizeDescription?: boolean;

      onOpenFullPage?: () => void;
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

  function openFullPage() {
    if (props.mode === "edit" && props.task?.id) {
      const hashedId = encodeId(props.task.id);

      router.push(`/tasks/${hashedId}`);
    }
  }

  function FullButton({
    href,
    onClick,
  }: {
    href?: string;

    onClick?: () => void;
  }) {
    const className = `
      h-10
      gap-1.5
      rounded-xl
      px-3
      text-muted-foreground
      hover:bg-muted
      hover:text-foreground
    `;

    const content = (
      <>
        <Expand className="h-4 w-4 shrink-0" />

        <span
          className="
            text-[10px]
            font-medium
            tracking-[0.2em]
            opacity-60
          "
        >
          FULL
        </span>
      </>
    );

    if (href) {
      return (
        <Button asChild type="button" variant="ghost" className={className}>
          <Link href={href}>{content}</Link>
        </Button>
      );
    }

    return (
      <Button
        type="button"
        variant="ghost"
        className={className}
        onClick={onClick}
      >
        {content}
      </Button>
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
            <FullButton href="/tasks/new" />
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
        props.showOpenPageButton ? <FullButton onClick={openFullPage} /> : null
      }
      footer={
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button
              variant="destructive"
              className="
                h-12
                w-full
                rounded-xl
                px-6
                sm:w-auto
              "
            >
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

// src/features/tasks/components/tasks-table.tsx

"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { Trash2 } from "lucide-react";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

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

import { formatDateOnly } from "@/lib/utils/date";

import { useTasks } from "../hooks/use-tasks";
import { useDeleteTask } from "../hooks/use-delete-task";
import { useUpdateTaskStatus } from "../hooks/use-update-task-status";

import type { Task } from "../schemas/task-schema";

import { TasksPagination } from "./tasks-pagination";
import { TaskStatusSelect } from "./task-status-select";

type Props = {
  limit: number;
  offset: number;
  status?: "TODO" | "DOING" | "DONE";
  sort?: "created_at" | "due_date";
  order?: "ASC" | "DESC";
  onOpenTask: (task: Task) => void;
};

export function TasksTable({
  limit,
  offset,
  status,
  sort,
  order,
  onOpenTask,
}: Props) {
  const router = useRouter();

  const searchParams = useSearchParams();

  const updateStatus = useUpdateTaskStatus();

  const deleteTask = useDeleteTask();

  const { data, isLoading, error } = useTasks({
    limit,
    offset,
    status,
    sort,
    order,
  });

  if (isLoading) {
    return null;
  }

  if (error) {
    console.error(error);

    return (
      <div className="rounded-xl border p-6 text-sm text-red-500">
        Failed to load tasks
      </div>
    );
  }

  if (!data) {
    return null;
  }

  const tasks: Task[] = data.items;

  function openTask(task: Task) {
    onOpenTask(task);
  }

  if (tasks.length === 0) {
    return (
      <div className="rounded-xl border p-8 text-center">
        <p className="text-muted-foreground text-sm">No tasks found</p>

        <div className="border-t px-6 pt-6">
          <TasksPagination
            total={data.count ?? 0}
            limit={data.limit ?? limit}
            offset={data.offset ?? offset}
          />
        </div>
      </div>
    );
  }

  return (
    <div
      className="
        animate-in
        fade-in-0
        slide-in-from-bottom-4
        duration-500
        space-y-4
      "
    >
      {/* mobile */}
      <div className="space-y-3 md:hidden">
        {tasks.map((task) => (
          <div
            key={task.publicId}
            role="button"
            tabIndex={0}
            onClick={() => openTask(task)}
            className="space-y-4 rounded-xl border p-4 transition-colors hover:bg-muted/50"
          >
            <div className="space-y-2">
              <p className="line-clamp-1 wrap-break-word font-medium">
                {task.title}
              </p>

              {task.description && (
                <p
                  className="
                    text-muted-foreground
                    line-clamp-2
                    wrap-break-word
                    text-sm
                  "
                >
                  {task.description}
                </p>
              )}
            </div>

            <div className="flex items-center justify-between gap-3">
              <div
                className="w-28 shrink-0"
                onClick={(e) => e.stopPropagation()}
              >
                <TaskStatusSelect
                  value={task.status}
                  onChange={(value) =>
                    updateStatus.mutate({
                      publicId: task.publicId,
                      input: {
                        status: value,
                      },
                    })
                  }
                  className="h-10 w-full rounded-lg text-xs"
                />
              </div>

              <div className="shrink-0 text-sm">
                {task.dueDate ? (
                  <span>{formatDateOnly(task.dueDate)}</span>
                ) : (
                  <span className="text-muted-foreground">No date</span>
                )}
              </div>

              <div onClick={(e) => e.stopPropagation()}>
                <AlertDialog>
                  <AlertDialogTrigger asChild>
                    <Button
                      size="icon"
                      variant="ghost"
                      className="h-9 w-9 text-red-500 hover:text-red-600"
                    >
                      <Trash2 className="h-4 w-4" />
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
                        onClick={() => deleteTask.mutate(task.publicId)}
                      >
                        {deleteTask.isPending ? "Deleting..." : "Delete"}
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* desktop */}
      <div className="hidden rounded-xl border md:block">
        <Table className="table-fixed">
          <TableHeader>
            <TableRow>
              <TableHead className="w-[50%] pl-4 text-left">Title</TableHead>

              <TableHead className="w-[25%] text-left">Status</TableHead>

              <TableHead className="w-[20%] text-left">Due Date</TableHead>

              <TableHead className="w-16" />
            </TableRow>
          </TableHeader>

          <TableBody>
            {tasks.map((task) => (
              <TableRow
                key={task.publicId}
                className="cursor-pointer transition-colors hover:bg-muted/50"
                onClick={() => openTask(task)}
              >
                <TableCell className="pl-4">
                  <div className="space-y-1">
                    <p className="overflow-hidden text-ellipsis wrap-break-word font-medium">
                      {task.title}
                    </p>

                    {task.description && (
                      <p
                        className="
                          text-muted-foreground
                          line-clamp-2
                          overflow-hidden
                          text-ellipsis
                          wrap-break-word
                          text-sm
                        "
                      >
                        {task.description}
                      </p>
                    )}
                  </div>
                </TableCell>

                <TableCell onClick={(e) => e.stopPropagation()}>
                  <TaskStatusSelect
                    value={task.status}
                    onChange={(value) =>
                      updateStatus.mutate({
                        publicId: task.publicId,
                        input: {
                          status: value,
                        },
                      })
                    }
                  />
                </TableCell>

                <TableCell>
                  {task.dueDate ? (
                    <span className="text-sm">
                      {formatDateOnly(task.dueDate)}
                    </span>
                  ) : (
                    <span className="text-muted-foreground text-sm">-</span>
                  )}
                </TableCell>

                <TableCell
                  className="text-right"
                  onClick={(e) => e.stopPropagation()}
                >
                  <AlertDialog>
                    <AlertDialogTrigger asChild>
                      <Button
                        size="icon"
                        variant="ghost"
                        className="h-8 w-8 text-red-500 hover:text-red-600"
                      >
                        <Trash2 className="h-4 w-4" />
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
                          onClick={() => deleteTask.mutate(task.publicId)}
                        >
                          {deleteTask.isPending ? "Deleting..." : "Delete"}
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <div className="rounded-xl border px-4 py-4 md:px-6">
        <TasksPagination
          total={data.count ?? 0}
          limit={data.limit ?? limit}
          offset={data.offset ?? offset}
        />
      </div>
    </div>
  );
}

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
import { encodeId } from "@/lib/utils/hash-id";

import { useTasks } from "../hooks/use-tasks";
import { useDeleteTask } from "../hooks/use-delete-task";
import { useUpdateTaskStatus } from "../hooks/use-update-task-status";

import type { Task, TaskListResponse } from "../schemas/task-schema";
import { TasksPagination } from "./tasks-pagination";
import { TaskStatusSelect } from "./task-status-select";

type Props = {
  initialData: TaskListResponse;
  limit: number;
  offset: number;
  status?: "TODO" | "DOING" | "DONE";
  sort?: "created_at" | "due_date";
  order?: "ASC" | "DESC";
};

export function TasksTable({
  initialData,
  limit,
  offset,
  status,
  sort,
  order,
}: Props) {
  const router = useRouter();
  const searchParams = useSearchParams();

  const updateStatus = useUpdateTaskStatus();
  const deleteTask = useDeleteTask();

  const { data } = useTasks(
    { limit, offset, status, sort, order },
    { initialData },
  );

  const tasks: Task[] = data?.items ?? [];

  function openTask(taskId: number) {
    const params = new URLSearchParams(searchParams.toString());
    params.set("taskId", encodeId(taskId));
    router.push(`/tasks?${params.toString()}`);
  }

  if (tasks.length === 0) {
    return (
      <div className="rounded-xl border p-8 text-center">
        <p className="text-muted-foreground text-sm">No tasks found</p>
        <div className="border-t px-6 pt-6">
          <TasksPagination
            total={data?.count ?? 0}
            limit={data?.limit ?? limit}
            offset={data?.offset ?? offset}
          />
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Mobile Cards */}
      <div className="space-y-3 md:hidden">
        {tasks.map((task: Task) => (
          <div
            key={task.id}
            role="button"
            tabIndex={0}
            onClick={() => openTask(task.id)}
            onKeyDown={(e) => {
              if (e.key === "Enter" || e.key === " ") {
                e.preventDefault();
                openTask(task.id);
              }
            }}
            className="
        space-y-4
        rounded-xl
        border
        p-4
        transition-colors
        hover:bg-muted/50
        focus:outline-none
        focus:ring-2
        focus:ring-ring
      "
          >
            <div className="space-y-2">
              <p className="line-clamp-2 break-all font-medium">{task.title}</p>

              {task.description && (
                <p className="text-muted-foreground line-clamp-3 break-all text-sm">
                  {task.description}
                </p>
              )}
            </div>

            <div className="flex items-center justify-between gap-3">
              <div
                onClick={(e) => e.stopPropagation()}
                className="min-w-0 flex-1"
              >
                <TaskStatusSelect
                  value={task.status}
                  onChange={(value) => {
                    updateStatus.mutate({
                      id: task.id,
                      status: value,
                    });
                  }}
                  className="h-10 w-full rounded-lg text-sm"
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
                        onClick={() => {
                          deleteTask.mutate(task.id);
                        }}
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
      {/* Desktop Table */}
      <div className="hidden rounded-xl border md:block">
        <Table className="table-fixed">
          <TableHeader>
            <TableRow>
              <TableHead className="w-[50%]">Title</TableHead>
              <TableHead className="w-[25%]">Status</TableHead>
              <TableHead className="w-[20%]">Due Date</TableHead>
              <TableHead className="w-16" />
            </TableRow>
          </TableHeader>

          <TableBody>
            {tasks.map((task: Task) => (
              <TableRow
                key={task.id}
                className="cursor-pointer hover:bg-muted/50"
                onClick={() => openTask(task.id)}
              >
                <TableCell>
                  <div className="min-w-0 space-y-1">
                    <p className="truncate font-medium">{task.title}</p>
                    {task.description && (
                      <p className="text-muted-foreground line-clamp-2 break-all text-sm">
                        {task.description}
                      </p>
                    )}
                  </div>
                </TableCell>

                <TableCell onClick={(e) => e.stopPropagation()}>
                  <TaskStatusSelect
                    value={task.status}
                    onChange={(value) => {
                      updateStatus.mutate({
                        id: task.id,
                        status: value,
                      });
                    }}
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
                          onClick={() => {
                            deleteTask.mutate(task.id);
                          }}
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
          total={data?.count ?? 0}
          limit={data?.limit ?? limit}
          offset={data?.offset ?? offset}
        />
      </div>
    </div>
  );
}

"use client";

// src/app/(protected)/tasks/page.tsx

import Link from "next/link";

import { useMemo } from "react";

import { useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";

import { useTasks } from "@/features/tasks/hooks/use-tasks";

import { TasksTable } from "@/features/tasks/components/tasks-table";

import { EditTaskDialog } from "@/features/tasks/components/edit-task-dialog";

import { CreateTaskDialog } from "@/features/tasks/components/create-task-dialog";

import { TasksFilter } from "@/features/tasks/components/tasks-filter";

import { TasksSort } from "@/features/tasks/components/tasks-sort";

import { TasksTableSkeleton } from "@/features/tasks/components/tasks-table-skeleton";

export default function TasksPage() {
  const searchParams = useSearchParams();

  const limit = Number(searchParams.get("limit") ?? 10);

  const offset = Number(searchParams.get("offset") ?? 0);

  const status =
    (searchParams.get("status") as "TODO" | "DOING" | "DONE" | null) ??
    undefined;

  const sort =
    (searchParams.get("sort") as "created_at" | "due_date" | null) ?? undefined;

  const order =
    (searchParams.get("order") as "ASC" | "DESC" | null) ?? undefined;

  // =========================
  // decode hashed task id
  // =========================

  const encodedTaskId = searchParams.get("taskId");

  const decodedTaskId = useMemo(() => {
    if (!encodedTaskId) {
      return null;
    }

    try {
      const decoded = Number(atob(decodeURIComponent(encodedTaskId)));

      if (!Number.isInteger(decoded) || decoded <= 0) {
        return null;
      }

      return decoded;
    } catch {
      return null;
    }
  }, [encodedTaskId]);
  const { data, isPending } = useTasks({
    limit,
    offset,
    status,
    sort,
    order,
  });

  return (
    <div className="mx-auto w-full max-w-7xl space-y-6 p-4 sm:p-6 lg:p-8">
      {/* top */}

      <div className="space-y-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight sm:text-4xl">
            Task Management
          </h1>

          <p className="text-muted-foreground mt-2 text-sm sm:text-base">
            Manage your tasks efficiently
          </p>
        </div>

        <div className="flex justify-start">
          <Button
            asChild
            size="lg"
            className="h-12 w-full rounded-xl px-8 text-base font-medium sm:w-auto"
          >
            <Link href="/tasks?create=true">Create Task</Link>
          </Button>
        </div>
      </div>

      {/* controls */}

      <div className="flex flex-col gap-4 rounded-2xl border p-4 sm:flex-row sm:items-center sm:justify-between sm:p-5">
        <TasksFilter />

        <TasksSort />
      </div>

      {/* table */}

      {isPending ? (
        <TasksTableSkeleton />
      ) : (
        <TasksTable
          initialData={
            data ?? {
              items: [],
              count: 0,
              limit,
              offset,
            }
          }
          limit={limit}
          offset={offset}
          status={status}
          sort={sort}
          order={order}
        />
      )}

      {/* dialogs */}

      <CreateTaskDialog />

      {decodedTaskId && <EditTaskDialog taskId={decodedTaskId} />}
    </div>
  );
}

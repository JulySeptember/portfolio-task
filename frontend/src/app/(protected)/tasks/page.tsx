// src/app/(protected)/tasks/page.tsx
"use client";

import { useMemo, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import Link from "next/link";

import { Button } from "@/components/ui/button";
import { useTasks } from "@/features/tasks/hooks/use-tasks";
import { TasksTable } from "@/features/tasks/components/tasks-table";
import { EditTaskDialog } from "@/features/tasks/components/edit-task-dialog";
import { CreateTaskDialog } from "@/features/tasks/components/create-task-dialog";
import { TasksFilter } from "@/features/tasks/components/tasks-filter";
import { TasksSort } from "@/features/tasks/components/tasks-sort";

import type { Task } from "@/features/tasks/schemas/task-schema";

export default function TasksPage() {
  const searchParams = useSearchParams();
  const router = useRouter();

  const limit = Number(searchParams.get("limit") ?? 10);
  const offset = Number(searchParams.get("offset") ?? 0);

  const status =
    (searchParams.get("status") as "TODO" | "DOING" | "DONE") ?? undefined;

  const sort =
    (searchParams.get("sort") as "created_at" | "due_date") ?? undefined;

  const order = (searchParams.get("order") as "ASC" | "DESC") ?? undefined;

  const taskIdFromUrl = searchParams.get("taskId");

  const [selectedTask, setSelectedTask] = useState<Task | null>(null);

  const { data } = useTasks({
    limit,
    offset,
    status,
    sort,
    order,
  });

  const openedTask = useMemo(() => {
    if (!taskIdFromUrl || !data) return null;

    return data.items.find((t) => t.publicId === taskIdFromUrl) ?? null;
  }, [taskIdFromUrl, data]);

  function handleOpenTask(task: Task) {
    setSelectedTask(task);

    const params = new URLSearchParams(searchParams.toString());

    params.set("taskId", task.publicId);

    router.push(`/tasks?${params.toString()}`, {
      scroll: false,
    });
  }

  function handleCloseTask() {
    setSelectedTask(null);

    const params = new URLSearchParams(searchParams.toString());

    params.delete("taskId");

    router.push(params.toString() ? `/tasks?${params.toString()}` : "/tasks", {
      scroll: false,
    });
  }

  return (
    <div className="mx-auto w-full max-w-7xl space-y-6 p-4 sm:p-6 lg:p-8">
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

      <div className="flex flex-col gap-4 rounded-2xl border p-4 sm:flex-row sm:items-center sm:justify-between sm:p-5">
        <TasksFilter />

        <TasksSort />
      </div>

      <TasksTable
        limit={limit}
        offset={offset}
        status={status}
        sort={sort}
        order={order}
        onOpenTask={handleOpenTask}
      />

      <CreateTaskDialog />

      {(openedTask || selectedTask) && (
        <EditTaskDialog
          task={openedTask ?? selectedTask!}
          onClose={handleCloseTask}
        />
      )}
    </div>
  );
}

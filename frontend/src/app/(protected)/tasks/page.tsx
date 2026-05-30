// src/app/(protected)/tasks/page.tsx
"use client";

import { useMemo, useCallback } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import Link from "next/link";

import { Button } from "@/components/ui/button";
import { useTasks } from "@/features/tasks/hooks/use-tasks";
import { TasksTable } from "@/features/tasks/components/tasks-table";
import { TaskDialog } from "@/features/tasks/components/task-dialog";
import { TasksFilter } from "@/features/tasks/components/tasks-filter";
import { TasksSort } from "@/features/tasks/components/tasks-sort";

import type { Task } from "@/features/tasks/schemas/task-schema";

export default function TasksPage() {
  const searchParams = useSearchParams();
  const router = useRouter();

  const limit = Math.min(
    Math.max(Number(searchParams.get("limit") ?? 10), 1),
    50,
  );
  const offset = Math.min(
    Math.max(Number(searchParams.get("offset") ?? 0), 0),
    1000,
  );

  const status =
    (searchParams.get("status") as "TODO" | "DOING" | "DONE") ?? undefined;
  const sort =
    (searchParams.get("sort") as "created_at" | "due_date") ?? undefined;
  const order = (searchParams.get("order") as "ASC" | "DESC") ?? undefined;

  const taskIdFromUrl = searchParams.get("taskId");
  const createOpen = searchParams.get("create") === "true";

  const { data } = useTasks({ limit, offset, status, sort, order });

  // URLを唯一の状態管理に
  const openedTask = useMemo(() => {
    if (!taskIdFromUrl || !data) return null;
    return data.items.find((task) => task.publicId === taskIdFromUrl) ?? null;
  }, [taskIdFromUrl, data]);

  // タスクを開く
  const handleOpenTask = useCallback(
    (task: Task) => {
      const params = new URLSearchParams(searchParams.toString());
      params.set("taskId", task.publicId);
      router.push(`/tasks?${params.toString()}`, { scroll: false });
    },
    [searchParams, router],
  );

  // タスクを閉じる
  const handleCloseTask = useCallback(() => {
    const params = new URLSearchParams(searchParams.toString());
    params.delete("taskId");
    router.push(params.toString() ? `/tasks?${params.toString()}` : "/tasks", {
      scroll: false,
    });
  }, [searchParams, router]);

  // 作成ダイアログ開閉
  const handleCreateOpenChange = useCallback(
    (open: boolean) => {
      const params = new URLSearchParams(searchParams.toString());
      if (open) params.set("create", "true");
      else params.delete("create");
      router.push(
        params.toString() ? `/tasks?${params.toString()}` : "/tasks",
        { scroll: false },
      );
    },
    [searchParams, router],
  );

  return (
    <div className="mx-auto w-full max-w-7xl space-y-6 p-4 sm:p-6 lg:p-8">
      {/* Header */}
      <div className="space-y-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight sm:text-4xl">
            Task Management
          </h1>
          <p className="mt-2 text-sm text-muted-foreground sm:text-base">
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

      {/* Filter / Sort */}
      <div className="flex flex-col gap-4 rounded-2xl border p-4 sm:flex-row sm:items-center sm:justify-between sm:p-5">
        <TasksFilter />
        <TasksSort />
      </div>

      {/* Table */}
      <TasksTable
        limit={limit}
        offset={offset}
        status={status}
        sort={sort}
        order={order}
        onOpenTask={handleOpenTask}
      />

      {/* Create Task Dialog */}
      <TaskDialog
        mode="create"
        open={createOpen}
        onOpenChange={handleCreateOpenChange}
      />

      {/* Edit Task Dialog */}
      {openedTask && (
        <TaskDialog
          mode="edit"
          task={openedTask}
          open
          onOpenChange={handleCloseTask}
        />
      )}
    </div>
  );
}

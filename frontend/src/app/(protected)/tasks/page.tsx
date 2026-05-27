// src/app/(protected)/tasks/page.tsx

import Link from "next/link";

import { Button } from "@/components/ui/button";

import { getTasks } from "@/features/tasks/server/get-tasks";

import { TasksTable } from "@/features/tasks/components/tasks-table";

import { TaskDetailDialog } from "@/features/tasks/components/task-detail-dialog";

import { CreateTaskDialog } from "@/features/tasks/components/create-task-dialog";

import { TasksFilter } from "@/features/tasks/components/tasks-filter";

import { TasksSort } from "@/features/tasks/components/tasks-sort";

import { getCurrentUser } from "@/features/auth/api/get-current-user";

type Props = {
  searchParams: Promise<{
    limit?: string;
    offset?: string;
    status?: "TODO" | "DOING" | "DONE";
    sort?: "created_at" | "due_date";
    order?: "ASC" | "DESC";
    taskId?: string;
    create?: string;
  }>;
};

export default async function TasksPage({ searchParams }: Props) {
  const params = await searchParams;

  const limit = Number(params.limit ?? 10);

  const offset = Number(params.offset ?? 0);

  const [response, currentUser] = await Promise.all([
    getTasks({
      limit,
      offset,
      status: params.status,
      sort: params.sort,
      order: params.order,
    }),

    getCurrentUser(),
  ]);

  return (
    <div className="mx-auto w-full max-w-7xl space-y-8 p-8">
      {/* ========================= */}
      {/* top */}
      {/* ========================= */}

      <div className="flex items-start justify-between gap-6">
        <div className="space-y-3">
          <div>
            <h1 className="text-4xl font-bold tracking-tight">
              Task Management
            </h1>

            <p className="text-muted-foreground mt-2 text-base">
              Manage your tasks efficiently
            </p>
          </div>
        </div>
      </div>

      {/* ========================= */}
      {/* controls */}
      {/* ========================= */}

      <div className="flex flex-wrap items-center justify-between gap-4 rounded-2xl border p-5">
        <div className="flex flex-wrap items-center gap-3">
          <TasksFilter />

          <TasksSort />
        </div>

        <Button asChild size="lg" className="h-12 px-8 text-base font-medium">
          <Link href="/tasks?create=true">Create Task</Link>
        </Button>
      </div>

      {/* ========================= */}
      {/* table */}
      {/* ========================= */}

      <TasksTable
        initialData={response}
        limit={limit}
        offset={offset}
        status={params.status}
        sort={params.sort}
        order={params.order}
      />

      {/* ========================= */}
      {/* dialogs */}
      {/* ========================= */}

      <CreateTaskDialog />

      {params.taskId && <TaskDetailDialog taskId={Number(params.taskId)} />}
    </div>
  );
}

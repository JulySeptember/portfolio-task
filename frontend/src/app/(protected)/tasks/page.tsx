// src/app/(protected)/tasks/page.tsx

import Link from "next/link";

import { Button } from "@/components/ui/button";

import { getTasks } from "@/features/tasks/server/get-tasks";

import { TasksTable } from "@/features/tasks/components/tasks-table";

import { EditTaskDialog } from "@/features/tasks/components/edit-task-dialog";

import { CreateTaskDialog } from "@/features/tasks/components/create-task-dialog";

import { TasksFilter } from "@/features/tasks/components/tasks-filter";

import { TasksSort } from "@/features/tasks/components/tasks-sort";

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

  const response = await getTasks({
    limit,
    offset,
    status: params.status,
    sort: params.sort,
    order: params.order,
  });

  return (
    <div className="mx-auto w-full max-w-7xl space-y-6 p-4 sm:p-6 lg:p-8">
      {/* ========================= */}
      {/* top */}
      {/* ========================= */}

      <div className="space-y-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight sm:text-4xl">
            Task Management
          </h1>

          <p className="text-muted-foreground mt-2 text-sm sm:text-base">
            Manage your tasks efficiently
          </p>
        </div>
        {/* create button outside */}
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

      {/* ========================= */}
      {/* controls */}
      {/* ========================= */}

      <div className="flex flex-col gap-4 rounded-2xl border p-4 sm:flex-row sm:items-center sm:justify-between sm:p-5">
        <TasksFilter />

        <TasksSort />
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

      {params.taskId && <EditTaskDialog taskId={Number(params.taskId)} />}
    </div>
  );
}

"use client";

import { useState } from "react";

import { Button } from "@/components/ui/button";

import { TaskForm } from "@/features/tasks/components/task-form";

import { TasksTable } from "@/features/tasks/components/tasks-table";

import { useTasks } from "@/features/tasks/hooks/use-tasks";

export default function TasksPage() {
  const [status, setStatus] = useState<"TODO" | "DOING" | "DONE" | undefined>();

  const [sort, setSort] = useState<"created_at" | "due_date">("created_at");

  const [order, setOrder] = useState<"ASC" | "DESC">("DESC");

  const [offset, setOffset] = useState(0);

  const limit = 10;

  const { data, isPending, isError, error } = useTasks({
    status,
    sort,
    order,
    limit,
    offset,
  });

  if (isPending) {
    return <div>Loading...</div>;
  }

  if (isError) {
    return <div>{error.message}</div>;
  }

  const hasPrevPage = offset > 0;

  const hasNextPage = data.offset + data.items.length < data.count;

  return (
    <main className="space-y-6 p-6">
      <div>
        <h1 className="text-2xl font-bold">Tasks</h1>

        <p className="text-muted-foreground text-sm">Total: {data.count}</p>
      </div>

      <TaskForm />

      <div className="flex flex-wrap items-center gap-2">
        <Button
          variant={!status ? "default" : "outline"}
          onClick={() => {
            setStatus(undefined);

            setOffset(0);
          }}
        >
          ALL
        </Button>

        <Button
          variant={status === "TODO" ? "default" : "outline"}
          onClick={() => {
            setStatus("TODO");

            setOffset(0);
          }}
        >
          TODO
        </Button>

        <Button
          variant={status === "DOING" ? "default" : "outline"}
          onClick={() => {
            setStatus("DOING");

            setOffset(0);
          }}
        >
          DOING
        </Button>

        <Button
          variant={status === "DONE" ? "default" : "outline"}
          onClick={() => {
            setStatus("DONE");

            setOffset(0);
          }}
        >
          DONE
        </Button>
      </div>

      <div className="flex flex-wrap items-center gap-2">
        <Button
          variant={sort === "created_at" ? "default" : "outline"}
          onClick={() => {
            setSort("created_at");

            setOffset(0);
          }}
        >
          Created
        </Button>

        <Button
          variant={sort === "due_date" ? "default" : "outline"}
          onClick={() => {
            setSort("due_date");

            setOffset(0);
          }}
        >
          Due Date
        </Button>

        <Button
          variant="outline"
          onClick={() => {
            setOrder((prev) => (prev === "ASC" ? "DESC" : "ASC"));

            setOffset(0);
          }}
        >
          {order}
        </Button>
      </div>

      <TasksTable tasks={data.items} />

      <div className="flex items-center justify-between">
        <Button
          variant="outline"
          disabled={!hasPrevPage}
          onClick={() => setOffset((prev) => Math.max(prev - limit, 0))}
        >
          Previous
        </Button>

        <p className="text-muted-foreground text-sm">
          {offset + 1} - {offset + data.items.length} / {data.count}
        </p>

        <Button
          variant="outline"
          disabled={!hasNextPage}
          onClick={() => setOffset((prev) => prev + limit)}
        >
          Next
        </Button>
      </div>
    </main>
  );
}

"use client";

import { TaskForm } from "@/features/tasks/components/task-form";

import { TasksTable } from "@/features/tasks/components/tasks-table";

import { useTasks } from "@/features/tasks/hooks/use-tasks";

export default function TasksPage() {
  const { data, isPending, isError, error } = useTasks();

  if (isPending) {
    return <div>Loading...</div>;
  }

  if (isError) {
    return <div>{error.message}</div>;
  }

  return (
    <main className="space-y-6 p-6">
      <div>
        <h1 className="text-2xl font-bold">Tasks</h1>

        <p className="text-muted-foreground text-sm">Total: {data.count}</p>
      </div>

      <TaskForm />

      <TasksTable tasks={data.items} />
    </main>
  );
}

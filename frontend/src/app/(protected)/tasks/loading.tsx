import { Skeleton } from "@/components/ui/skeleton";

import { TasksFiltersSkeleton } from "@/features/tasks/components/tasks-filters-skeleton";

import { TasksTableSkeleton } from "@/features/tasks/components/tasks-table-skeleton";

export default function Loading() {
  return (
    <main className="space-y-6 p-6">
      <div className="space-y-2">
        <Skeleton className="h-8 w-48" />

        <Skeleton className="h-5 w-72" />
      </div>

      <TasksFiltersSkeleton />

      <TasksTableSkeleton />

      <div className="flex justify-center gap-2">
        <Skeleton className="h-10 w-10" />

        <Skeleton className="h-10 w-10" />

        <Skeleton className="h-10 w-10" />
      </div>
    </main>
  );
}

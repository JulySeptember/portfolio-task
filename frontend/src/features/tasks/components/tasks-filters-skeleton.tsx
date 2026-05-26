import { Skeleton } from "@/components/ui/skeleton";

export function TasksFiltersSkeleton() {
  return (
    <div className="flex flex-wrap items-center gap-4">
      <Skeleton className="h-10 w-40" />

      <Skeleton className="h-10 w-40" />

      <Skeleton className="h-10 w-32" />

      <Skeleton className="ml-auto h-10 w-36" />
    </div>
  );
}

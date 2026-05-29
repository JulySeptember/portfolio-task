import { Skeleton } from "@/components/ui/skeleton";

export function TasksTableSkeleton() {
  return (
    <div className="rounded-xl border">
      {/* header */}
      <div className="border-b p-4">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3">
          <Skeleton className="h-5 w-32" />

          <Skeleton className="h-5 w-24" />

          <Skeleton className="h-5 w-28" />
        </div>
      </div>

      {/* rows */}
      <div className="divide-y">
        {Array.from({ length: 6 }).map((_, index) => (
          <div
            key={index}
            className="
              grid
              grid-cols-1
              gap-4
              p-4
              sm:grid-cols-2
              md:grid-cols-3
            "
          >
            <div className="space-y-2">
              <Skeleton className="h-5 w-52 max-w-full" />

              <Skeleton className="h-4 w-72 max-w-full" />
            </div>

            <Skeleton className="h-10 w-32 rounded-md max-w-full" />

            <Skeleton className="h-5 w-40 max-w-full" />
          </div>
        ))}
      </div>
    </div>
  );
}

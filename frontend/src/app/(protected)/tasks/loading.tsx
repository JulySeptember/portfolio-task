import { Skeleton } from "@/components/ui/skeleton";

export default function Loading() {
  return (
    <main className="space-y-6 p-6">
      <div className="space-y-2">
        <Skeleton className="h-8 w-40" />

        <Skeleton className="h-4 w-24" />
      </div>

      <div className="rounded-lg border p-4">
        <div className="space-y-4">
          <Skeleton className="h-10 w-full" />

          <Skeleton className="h-24 w-full" />

          <Skeleton className="h-10 w-32" />
        </div>
      </div>

      <div className="rounded-lg border">
        <div className="space-y-3 p-4">
          {Array.from({ length: 5 }).map((_, index) => (
            <div
              key={index}
              className="flex items-center justify-between gap-4 border-b pb-3 last:border-0"
            >
              <div className="space-y-2">
                <Skeleton className="h-4 w-48" />

                <Skeleton className="h-3 w-72" />
              </div>

              <div className="flex items-center gap-2">
                <Skeleton className="h-8 w-24" />

                <Skeleton className="h-8 w-20" />

                <Skeleton className="h-8 w-20" />
              </div>
            </div>
          ))}
        </div>
      </div>
    </main>
  );
}

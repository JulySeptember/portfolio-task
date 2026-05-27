"use client";

import { useEffect, useState } from "react";

import { useRouter, useSearchParams } from "next/navigation";

import { AppHeader } from "@/components/layout/app-header";

import { PageContainer } from "@/components/layout/page-container";

import { CreateTaskDialog } from "@/features/tasks/components/create-task-dialog";

import { TasksFilter } from "@/features/tasks/components/tasks-filter";

import { TasksPagination } from "@/features/tasks/components/tasks-pagination";

import { TasksSort } from "@/features/tasks/components/tasks-sort";

import { TasksTable } from "@/features/tasks/components/tasks-table";

import { useTasks } from "@/features/tasks/hooks/use-tasks";

export default function TasksPage() {
  const searchParams = useSearchParams();

  const router = useRouter();

  const limit = 10;

  const [status, setStatus] = useState<"TODO" | "DOING" | "DONE" | undefined>(
    (searchParams.get("status") as "TODO" | "DOING" | "DONE" | null) ??
      undefined,
  );

  const [sort, setSort] = useState<"created_at" | "due_date">(
    (searchParams.get("sort") as "created_at" | "due_date" | null) ??
      "created_at",
  );

  const [order, setOrder] = useState<"ASC" | "DESC">(
    (searchParams.get("order") as "ASC" | "DESC" | null) ?? "DESC",
  );

  const [offset, setOffset] = useState(() => {
    const page = Number(searchParams.get("page") ?? "1");

    return (page - 1) * limit;
  });

  useEffect(() => {
    const params = new URLSearchParams();

    if (status) {
      params.set("status", status);
    }

    params.set("sort", sort);

    params.set("order", order);

    params.set("page", String(offset / limit + 1));

    router.replace(`/tasks?${params.toString()}`);
  }, [status, sort, order, offset, limit, router]);

  const { data, isPending, isError, error } = useTasks({
    status,
    sort,
    order,
    limit,
    offset,
  });

  const hasPrevPage = offset > 0;

  const hasNextPage = !!data && data.offset + data.items.length < data.count;

  return (
    <>
      <AppHeader />

      <PageContainer>
        <div className="space-y-6">
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div>
              <h1 className="text-3xl font-bold">Tasks</h1>

              {data && (
                <p className="text-muted-foreground text-sm">
                  Total: {data.count}
                </p>
              )}
            </div>

            <CreateTaskDialog />
          </div>

          <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
            <TasksFilter
              status={status}
              onChange={(value) => {
                setStatus(value);

                setOffset(0);
              }}
            />

            <TasksSort
              sort={sort}
              order={order}
              onSortChange={(value) => {
                setSort(value);

                setOffset(0);
              }}
              onOrderChange={(value) => {
                setOrder(value);

                setOffset(0);
              }}
            />
          </div>

          {isPending && <div>Loading...</div>}

          {isError && (
            <div className="rounded-lg border border-red-500/20 bg-red-500/10 p-4 text-sm text-red-500">
              {error.message}
            </div>
          )}

          {data && (
            <>
              <TasksTable tasks={data.items} />

              <TasksPagination
                offset={offset}
                limit={limit}
                total={data.count}
                currentCount={data.items.length}
                hasPrevPage={hasPrevPage}
                hasNextPage={hasNextPage}
                onPrevious={() =>
                  setOffset((prev) => Math.max(prev - limit, 0))
                }
                onNext={() => setOffset((prev) => prev + limit)}
              />
            </>
          )}
        </div>
      </PageContainer>
    </>
  );
}

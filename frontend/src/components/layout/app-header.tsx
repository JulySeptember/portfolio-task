// src/components/layout/app-header.tsx

"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";

import { Button } from "@/components/ui/button";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

import { ChevronDown, LogOut, Trash2, User } from "lucide-react";

import { useCurrentUser } from "@/features/auth/hooks/use-current-user";
import { clearTokens } from "@/features/auth/lib/token-storage";
import { buildLogoutURL } from "@/features/auth/lib/hosted-ui";

import { apiClient } from "@/lib/api/client";

export function AppHeader() {
  const queryClient = useQueryClient();
  const { data: user } = useCurrentUser();

  const [confirmOpen, setConfirmOpen] = useState(false);
  const [finalConfirmOpen, setFinalConfirmOpen] = useState(false);

  function logout() {
    const accessToken = localStorage.getItem("access_token");
    const isMockAuth = accessToken === "local-dev-token";

    clearTokens();
    queryClient.clear();

    if (isMockAuth) {
      window.location.href = "/";
      return;
    }

    window.location.href = buildLogoutURL();
  }

  const deleteMutation = useMutation({
    mutationFn: async () => {
      await apiClient("/api/v1/users/me", {
        method: "DELETE",
      });
    },

    onSuccess: () => {
      logout();
    },
  });

  return (
    <>
      <header className="border-b bg-background">
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-8">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" className="max-w-sm gap-2">
                <User className="h-4 w-4" />

                <span className="max-w-60 truncate">
                  {user?.email ?? "User"}
                </span>

                <ChevronDown className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>

            <DropdownMenuContent align="start" className="w-56">
              <DropdownMenuItem onClick={logout} className="cursor-pointer">
                <LogOut className="mr-2 h-4 w-4" />
                Logout
              </DropdownMenuItem>

              <DropdownMenuItem
                onClick={() => setConfirmOpen(true)}
                className="cursor-pointer text-destructive focus:text-destructive"
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Delete Account
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </header>

      {/* 1st Confirmation */}

      <AlertDialog open={confirmOpen} onOpenChange={setConfirmOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Account?</AlertDialogTitle>

            <AlertDialogDescription>
              Your account and all tasks will be deleted.
            </AlertDialogDescription>
          </AlertDialogHeader>

          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>

            <AlertDialogAction
              onClick={() => {
                setConfirmOpen(false);
                setFinalConfirmOpen(true);
              }}
            >
              Continue
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* 2nd Confirmation */}

      <AlertDialog open={finalConfirmOpen} onOpenChange={setFinalConfirmOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>

            <AlertDialogDescription>
              This action cannot be undone. All tasks belonging to this account
              will be permanently deleted.
            </AlertDialogDescription>
          </AlertDialogHeader>

          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>

            <AlertDialogAction
              onClick={() => deleteMutation.mutate()}
              disabled={deleteMutation.isPending}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {deleteMutation.isPending ? "Deleting..." : "Delete Forever"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

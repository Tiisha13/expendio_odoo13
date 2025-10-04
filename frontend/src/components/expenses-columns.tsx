"use client";

import { ColumnDef } from "@tanstack/react-table";
import { ArrowUpDown, Trash2 } from "lucide-react";
import { format } from "date-fns";
import { Expense } from "@/types/api";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

const getStatusBadge = (status: string) => {
  switch (status) {
    case "pending":
      return <Badge variant="secondary">Pending</Badge>;
    case "approved":
      return <Badge className="bg-green-500">Approved</Badge>;
    case "rejected":
      return <Badge variant="destructive">Rejected</Badge>;
    default:
      return <Badge>{status}</Badge>;
  }
};

export const createExpenseColumns = (onDelete: (id: string) => void): ColumnDef<Expense>[] => [
  {
    accessorKey: "expense_date",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Date
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => {
      const date = row.getValue("expense_date") as string;
      return format(new Date(date), "MMM d, yyyy");
    },
  },
  {
    accessorKey: "category",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Category
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => {
      return <span className="capitalize">{row.getValue("category")}</span>;
    },
  },
  {
    accessorKey: "description",
    header: "Description",
    cell: ({ row }) => {
      return <span className="block max-w-xs truncate">{row.getValue("description")}</span>;
    },
  },
  {
    accessorKey: "amount",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Amount
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => {
      const amount = parseFloat(row.getValue("amount"));
      const currency = row.original.currency;
      return `${amount.toFixed(2)} ${currency}`;
    },
  },
  {
    accessorKey: "status",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Status
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => {
      return getStatusBadge(row.getValue("status"));
    },
  },
  {
    id: "actions",
    header: () => <div className="text-right">Actions</div>,
    cell: ({ row }) => {
      const expense = row.original;

      if (expense.status !== "pending") {
        return null;
      }

      return (
        <div className="text-right">
          <Button variant="ghost" size="sm" onClick={() => onDelete(expense.id)}>
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      );
    },
  },
];

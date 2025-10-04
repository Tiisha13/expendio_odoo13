"use client";

import { ColumnDef } from "@tanstack/react-table";
import { ArrowUpDown, CheckCircle, XCircle } from "lucide-react";
import { format } from "date-fns";
import { Approval } from "@/types/api";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

export const createApprovalColumns = (
  onApprove: (approval: Approval) => void,
  onReject: (approval: Approval) => void
): ColumnDef<Approval>[] => [
  {
    id: "employee",
    accessorFn: (row) =>
      `${row.expense?.user?.first_name || ""} ${row.expense?.user?.last_name || ""}`.trim(),
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Employee
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => {
      const approval = row.original;
      return (
        <span className="font-medium">
          {approval.expense?.user?.first_name} {approval.expense?.user?.last_name}
        </span>
      );
    },
  },
  {
    id: "expense_date",
    accessorFn: (row) => row.expense?.expense_date || "",
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
      const approval = row.original;
      const date = approval.expense?.expense_date;
      return date ? format(new Date(date), "MMM d, yyyy") : "";
    },
  },
  {
    id: "category",
    accessorFn: (row) => row.expense?.category || "",
    header: "Category",
    cell: ({ row }) => {
      const approval = row.original;
      return <span className="capitalize">{approval.expense?.category}</span>;
    },
  },
  {
    id: "description",
    accessorFn: (row) => row.expense?.description || "",
    header: "Description",
    cell: ({ row }) => {
      const approval = row.original;
      return <span className="block max-w-xs truncate">{approval.expense?.description}</span>;
    },
  },
  {
    id: "amount",
    accessorFn: (row) => row.expense?.amount || 0,
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
      const approval = row.original;
      const amount = approval.expense?.amount;
      const currency = approval.expense?.currency;
      return amount ? `${amount.toFixed(2)} ${currency}` : "";
    },
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      const approval = row.original;
      return <Badge variant="secondary">{approval.status}</Badge>;
    },
  },
  {
    id: "actions",
    header: () => <div className="text-right">Actions</div>,
    cell: ({ row }) => {
      const approval = row.original;
      return (
        <div className="flex justify-end gap-2">
          <Button variant="outline" size="sm" onClick={() => onApprove(approval)}>
            <CheckCircle className="mr-2 h-4 w-4" />
            Approve
          </Button>
          <Button variant="outline" size="sm" onClick={() => onReject(approval)}>
            <XCircle className="mr-2 h-4 w-4" />
            Reject
          </Button>
        </div>
      );
    },
  },
];

"use client";

import { ColumnDef } from "@tanstack/react-table";
import { ArrowUpDown, Trash2 } from "lucide-react";
import { User } from "@/types/api";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export const createUserColumns = (
  onRoleChange: (userId: string, role: string) => void,
  onManagerAssign: (userId: string, managerId: string) => void,
  onDelete: (userId: string) => void,
  managers: User[]
): ColumnDef<User>[] => [
  {
    accessorKey: "name",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Name
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => {
      const user = row.original;
      return (
        <span className="font-medium">
          {user.first_name} {user.last_name}
        </span>
      );
    },
  },
  {
    accessorKey: "email",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
        >
          Email
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
  },
  {
    accessorKey: "role",
    header: "Role",
    cell: ({ row }) => {
      const user = row.original;
      return (
        <Select value={user.role} onValueChange={(value) => onRoleChange(user.id, value)}>
          <SelectTrigger className="w-[130px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="employee">Employee</SelectItem>
            <SelectItem value="manager">Manager</SelectItem>
            <SelectItem value="admin">Admin</SelectItem>
          </SelectContent>
        </Select>
      );
    },
  },
  {
    accessorKey: "manager_id",
    header: "Manager",
    cell: ({ row }) => {
      const user = row.original;

      if (user.role !== "employee") {
        return null;
      }

      return (
        <Select
          value={user.manager_id || ""}
          onValueChange={(value) => onManagerAssign(user.id, value)}
        >
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Assign manager" />
          </SelectTrigger>
          <SelectContent>
            {managers.map((m) => (
              <SelectItem key={m.id} value={m.id}>
                {m.first_name} {m.last_name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      );
    },
  },
  {
    id: "actions",
    header: () => <div className="text-right">Actions</div>,
    cell: ({ row }) => {
      const user = row.original;
      return (
        <div className="text-right">
          <Button variant="ghost" size="sm" onClick={() => onDelete(user.id)}>
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      );
    },
  },
];
